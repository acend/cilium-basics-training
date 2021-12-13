---
title: "7. Transparent Encryption"
weight: 7
sectionnumber: 7
---
## Host traffic/endpoint traffic encryption

Cilium supports the transparent encryption of Cilium-managed host traffic and traffic between Cilium-managed endpoints either using IPsec or [WireGuard®](https://www.wireguard.com/).


## Task {{% param sectionnumber %}}.1: Increase cluster size

By default minikube create single node clusters. A a second node to the cluster:

```bash
minikube node add
```


## Task {{% param sectionnumber %}}.2: Deploy a minimal app for test traffic between nodes

Now we deploy a minimal client/server app running on different nodes:

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: traffic
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: client
  namespace: traffic
  labels:
    app: client
spec:
  replicas: 1
  selector:
    matchLabels:
      app: client
  template:
    metadata:
      labels:
        app: client
    spec:
      nodeSelector:
        kubernetes.io/hostname: minikube
      containers:
      - name: frontend-container
        image: docker.io/byrnedo/alpine-curl:0.1.8
        imagePullPolicy: IfNotPresent
        command: [ "/bin/ash", "-c", "sleep 1000000000" ]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: server
  namespace: traffic
  labels:
    app: server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: server
  template:
    metadata:
      labels:
        app: server
    spec:
      nodeSelector:
        kubernetes.io/hostname: minikube-m02
      containers:
      - name: server-container
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
  name: server
  namespace: traffic
  labels:
    app: server
spec:
  type: ClusterIP
  selector:
    app: server
  ports:
  - name: http
    port: 8080
```

```
kubectl apply -f traffic-app.yaml 
```

In a new terminal, we start a loop to generate minimal traffic between the nodes.
```bash
CLIENT=$(kubectl get pods -n traffic -l app=client -o jsonpath='{.items[0].metadata.name}')
while kubectl -n traffic exec -ti ${CLIENT} -- curl -so /dev/null --write-out '%{http_code}\n' backend:8080; do  
  sleep 10
done
```


## Task {{% param sectionnumber %}}.3:  Enable traffic encryption with WireGuard

On new clusters encryption can be enabled while installing with the cli: ```bash cilium install --encryption wireguard```. Since we already have a running cluster we install it by configuring it in the cilium configurationn map.

```bash
kubectl patch -n kube-system cm cilium-config --patch '{"data":{"enable-wireguard": "true", "enable-l7-proxy": "false", "enable-wireguard-userspace-fallback": "true"}}'
kubectl -n kube-system rollout restart daemonset cilium
kubectl -n kube-system rollout status daemonset cilium
```


### Verify encryption is working


Verify the number of peers in encryption is correct (should be the sum of nodes - 1)
```bash
kubectl -n kube-system exec -ti ds/cilium -- cilium status | grep Encryption
```

You should see something similiar to this (in this example we have a two node cluster):

```bash
Encryption:             Wireguard       [cilium_wg0 (Pubkey: XbTJd5Gnp7F8cG2Ymj6q11dBx8OtP1J5ZOAhswPiYAc=, Port: 51871, Peers: 1)]
```

We can check if the traffic is really sent to the WireGuard tunnel device cilium_wg0 (hit Ctrl+C to stop sniffing):

```bash
CILIUM_AGENT=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].metadata.name}")
kubectl debug -n kube-system -i ${CILIUM_AGENT} --image=nicolaka/netshoot -- tcpdump -ni cilium_wg0
```


## Task {{% param sectionnumber %}}.4: Cleanup

Stop the traffic generation loop in the second terminal, delete the pods and scale down the cluster.
```
kubectl delete ns traffic
minikube node delete minikube-m02
```

