---
title: "11. Cilium Service Mesh"
weight: 11
sectionnumber: 11
OnlyWhenNot: techlab
---
With release 1.12 Cilium enabled direct ingress support and service mesh features like layer 7 loadbalancing


## Task {{% param sectionnumber %}}.1: Installation


```bash
helm upgrade -i cilium cilium/cilium --version 1.12.0 \
  --namespace kube-system \
  --reuse-values \
  --set ingressController.enabled=true \
  --wait
```

For Kubernetes Ingress to work kubeProxyReplacement needs to be set to `strict` or `partial`. This is why we stay on the `kubeless` cluster.

Wait until cilium is ready (check with `cilium status`). For Ingress to work it is necessary to restart the agent and the operator.

```
kubectl -n kube-system rollout restart deployment/cilium-operator
kubectl -n kube-system rollout restart ds/cilium
```


## Task {{% param sectionnumber %}}.2: Create Ingress

Cilium Service Mesh can handle ingress traffic with its Envoy proxy.

We will use this feature to allow traffic to our simple app from outside the cluster. Create a file named `ingress.yaml` with the text below inside:

{{< readfile file="/content/en/docs/11/ingress.yaml" code="true" lang="yaml" >}}

Apply it with:

```bash
kubectl apply -f ingress.yaml
```

Check the ingress and the service:

```bash
kubectl describe ingress backend
kubectl get svc cilium-ingress-backend
```
We see that Cilium created a Service with type Loadbalancer for our Ingress. Unfortunately, Minikube has no loadbalancer deployed, in our setup the external IP will stay pending.

As a workaround, we can test the service from inside Kubernetes.

```bash
SERVICE_IP=$(kubectl get svc cilium-ingress-backend -ojsonpath={.spec.clusterIP})
kubectl run --rm=true -it --restart=Never --image=curlimages/curl -- curl http://${SERVICE_IP}/public
```

You should get the following output:

```
[
  {
    "id": 1,
    "body": "public information"
  }
]pod "curl" deleted
```


## Task {{% param sectionnumber %}}.3: Layer 7 Loadbalancing

Ingress alone is not really a Service Mesh feature. Let us test a traffic control example by loadbalancing a service inside the proxy.

Start by creating the second service. Create a file named `backend2.yaml` and put in the text below:

{{< readfile file="/content/en/docs/11/backend2.yaml" code="true" lang="yaml" >}}

Apply it:
```bash
kubectl apply -f backend2.yaml
```

Call it:
```bash
kubectl run --rm=true -it --restart=Never --image=curlimages/curl -- curl --connect-timeout 3 http://backend-2:8080/public
```

We see output very similiar to our simple application backend, but with a changed text.

As layer 7 loadbalancing requires traffic to be routed through the proxy, we will enable this for our backend Pods using a `CiliumNetworkPolicy` with HTTP rules. We will block access to `/public` and allow requests to `/private`:

Create a file `cnp-l7-sm.yaml` with the following content:

{{< readfile file="/content/en/docs/11/cnp-l7-sm.yaml" code="true" lang="yaml" >}}

And apply the `CiliumNetworkPolicy` with:

```bash
kubectl apply -f cnp-l7-sm.yaml
```

Until now only the backend service is replying to Ingress traffic. Now we configure Envoy to loadbalance the traffic 50/50 between backend and backend-2 with retries.
We are using a CustomResource called `CiliumEnvoyConfig` for this. Create a file `envoyconfig.yaml` with the following content:

{{< readfile file="/content/en/docs/11/envoyconfig.yaml" code="true" lang="yaml" >}}

{{% alert title="Note" color="primary" %}}
If you want to read more about Envoy configuration [Envoy Architectural Overview](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/http/http) is a good place to start.
{{% /alert %}}

Apply the `CiliumEnvoyConfig` with:

```bash
kubectl apply -f envoyconfig.yaml
```

Test it by running `curl` a few times -- different backends should respond:

```bash
for i in {1..10}; do
  kubectl run --rm=true -it --image=curlimages/curl --restart=Never curl -- curl  http://backend:8080/private
done
```

We see both backends replying. If you call it many times the distribution would be equal.

```bash
[                                                                                                                               [10/1834]
  {                                                                                                                                      
    "id": 1,                                                                                                                             
    "body": "another secret information from a different backend"                                                                                                 
  }                                                                                                                                      
]pod "curl" deleted                                                                                                                      
[                                                                                                                                        
  {                                                                                                                                      
    "id": 1,                                                                                                                             
    "body": "secret information"                                                                                                         
  }                                                                                                                                      
]pod "curl" deleted
```
This basic traffic control example shows only one function of Cilium Service Mesh, other features include i.e. TLS termination, support for tracing and canary-rollouts.
