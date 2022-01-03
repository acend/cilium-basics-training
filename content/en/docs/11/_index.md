---
title: "11. Cilium Service Mesh"
weight: 11
sectionnumber: 11
---
NOTE: Cilium Service Mesh is still in Beta you will find the relevante information [here](https://github.com/cilium/cilium-service-mesh-beta).


## Task {{% param sectionnumber %}}.1: Installation
As the Cilium Service Mesh is still in Beta and uses specific images we will use a dedicated Cluster for it.

```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version=1.23.0 -p serviceMesh
cilium install --version -service-mesh:v1.11.0-beta.1 --config enable-envoy-config=true --kube-proxy-replacement=probe
cilium hubble enable --ui 
```


## Task {{% param sectionnumber %}}.2: Create Ingress
We deploy [the sample app from chapter 3](https://cilium-basics.training.acend.ch/docs/03/#task-31-deploy-simple-application). 
´´´bash
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  labels:
    app: frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend-container
        image: docker.io/byrnedo/alpine-curl:0.1.8
        imagePullPolicy: IfNotPresent
        command: [ "/bin/ash", "-c", "sleep 1000000000" ]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: not-frontend
  labels:
    app: not-frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: not-frontend
  template:
    metadata:
      labels:
        app: not-frontend
    spec:
      containers:
      - name: not-frontend-container
        image: docker.io/byrnedo/alpine-curl:0.1.8
        imagePullPolicy: IfNotPresent
        command: [ "/bin/ash", "-c", "sleep 1000000000" ]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  labels:
    app: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend-container
        env:
        - name: PORT
          value: "8080"
        ports:
        - containerPort: 8080
        image: docker.io/cilium/json-mock:1.2
        imagePullPolicy: IfNotPresent
---
apiVersion: v1
kind: Service
metadata:
  name: backend
  labels:
    app: backend
spec:
  type: ClusterIP
  selector:
    app: backend
  ports:
  - name: http
    port: 8080
```

Now we add an Ingress resource
```bash
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: backend
spec:
  ingressClassName: cilium
  rules:
  - http:
      paths:
      - backend:
          service:
            name: backend
            port:
              number: 8080
        path: /
        pathType: Prefix
```
```bash
kubectl apply -f ingress.yaml
```
Check the Ingress service
```bash
kubectl describe ingress
kubect get svc cilium-ingress-backend
```

As we see the ingress is of Type LoadBalancer and did not get an External-IP. This is because we do not have an LoadBalancer deployed with minikube.
As a workaround we can test the service from inside kubernetes.


```bash
serviceIP=$(kubectl get svc cilium-ingress-backend -ojsonpath={.spec.clusterIP})
kubectl run --rm=true -it --restart=Never --image=docker.io/byrnedo/alpine-curl:0.1.8 -- curl http://${serviceIP}/public
```

## Task {{% param sectionnumber %}}.3: Layer 7 Loadbalancing

For loadbalancing we need another service:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend-2
  labels:
    app: backend-2
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend-2
  template:
    metadata:
      labels:
        app: backend-2
    spec:
      containers:
      - name: backend-container
        env:
        - name: PORT
          value: "8080"
        ports:
        - containerPort: 8080
        image: docker.io/jmalloc/echo-server
        imagePullPolicy: IfNotPresent
---
apiVersion: v1
kind: Service
metadata:
  name: backend-2
  labels:
    app: backend-2
spec:
  type: ClusterIP
  selector:
    app: backend-2
  ports:
  - name: http
    port: 8080
```

```bash
kubectl apply -f backend-2.yaml
```

Layer 7 loadbalancing will need to be routed through the proxy, we will enable this for our backend pods.

```bash
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "rule1"
spec:
  description: "enable L7 without blocking"
  endpointSelector:
    matchLabels:
      app: backend
  ingress:
  - fromEntities:
    - "all"
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
      rules:
        http:
        - method: "GET"
          path: "/"
---
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "rule2"
spec:
  description: "enable L7 without blocking"
  endpointSelector:
    matchLabels:
      app: backend-2
  ingress:
  - fromEntities:
    - "all"
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
      rules:
        http:
        - method: "GET"
          path: "/"
```

```bash
kubectl apply -f l7.yaml
```


Now we configure envoy to LB 50/50 backend services with retries.

```yaml
apiVersion: cilium.io/v2alpha1
kind: CiliumEnvoyConfig
metadata:
  name: envoy-lb-listener
spec:
  services:
    - name: backend
      namespace: default
    - name: backend-2
      namespace: default
  resources:
    - "@type": type.googleapis.com/envoy.config.listener.v3.Listener
      name: envoy-lb-listener
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                stat_prefix: envoy-lb-listener
                rds:
                  route_config_name: lb_route
                http_filters:
                  - name: envoy.filters.http.router
    - "@type": type.googleapis.com/envoy.config.route.v3.RouteConfiguration
      name: lb_route
      virtual_hosts:
        - name: "lb_route"
          domains: ["*"]
          routes:
            - match:
                prefix: "/private"
              route:
                weighted_clusters:
                  clusters:
                    - name: "default/backend"
                      weight: 50
                    - name: "default/backend-2"
                      weight: 50
                retry_policy:
                  retry_on: 5xx
                  num_retries: 3
                  per_try_timeout: 1s
    - "@type": type.googleapis.com/envoy.config.cluster.v3.Cluster
      name: "default/backend"
      connect_timeout: 5s
      lb_policy: ROUND_ROBIN
      type: EDS
      outlier_detection:
        split_external_local_origin_errors: true
        consecutive_local_origin_failure: 2
    - "@type": type.googleapis.com/envoy.config.cluster.v3.Cluster
      name: "default/backend-2"
      connect_timeout: 3s
      lb_policy: ROUND_ROBIN
      type: EDS
      outlier_detection:
        split_external_local_origin_errors: true
        consecutive_local_origin_failure: 2
```

Test it, run curl a few times...different backends should respond

```bash
kubectl run --rm=true -it --image=curlimages/curl --restart=Never curl -- sh
curl  http://backend-2:8080/private
exit
kubectl delete pods curl
```
