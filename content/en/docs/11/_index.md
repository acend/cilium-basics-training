---
title: "11. Cilium Service Mesh"
weight: 11
sectionnumber: 11
---
Cilium Service Mesh enables functions like ingress or layer 7 loadbalancing.


## Task {{% param sectionnumber %}}.1: Installation

Cilium Service Mesh is still in beta, if you want more information about the current status you can find it [here](https://github.com/cilium/cilium-service-mesh-beta). The beta version uses specific images, because of that we will use a dedicated cluster and install Cilium with the CLI.

{{% alert title="Note" color="primary" %}}
In case cluster1 is still running, stop it
```bash
minikube stop -p cluster1
```
to free up resources and speed up things.
{{% /alert %}}


```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version={{% param "kubernetesVersion" %}} -p servicemesh
cilium install --version -service-mesh:v1.11.0-beta.1 --config enable-envoy-config=true --kube-proxy-replacement=probe
```

Wait until cilium is ready (check with `cilium status`) and also enable the Hubble UI:

```bash
cilium hubble enable --ui 
```


## Task {{% param sectionnumber %}}.2: Create Ingress

Cilium Service Mesh can handle ingress traffic with its Envoy proxy.


We deploy [the sample app from chapter 3](https://cilium-basics-pr-109.training.acend.ch/docs/03/01/#task-311-install-the-hubble-cli).


{{< highlight yaml >}}{{< readfile file="content/en/docs/03/01/simple-app.yaml" >}}{{< /highlight >}}

Apply it with:

```bash
kubectl apply -f simple-app.yaml
```

Now we add an ingress resource. Create a file named `ingress.yaml` with the text below inside:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/ingress.yaml" >}}{{< /highlight >}}

Apply it with:

```bash
kubectl apply -f ingress.yaml
```

Check the ingress and the service:

```bash
kubectl describe ingress
kubectl get svc cilium-ingress-backend
```
We have successfully created an ingress service. Unfortunately, Minikube has no loadbalancer deployed and we will not get a public IP for our service (status stays pending)

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
{{% alert title="Note" color="primary" %}}
We can also use `minikube tunnel -p servicemesh` and then curl the Cluster-IP directly from our browser.
{{% /alert %}}


## Task {{% param sectionnumber %}}.3: Layer 7 Loadbalancing

Ingress alone is not really a Service Mesh feature. Let us test a traffic control example by loadbalancing a service inside the proxy.

Start by creating the second service. Create a file named `backend2.yaml` and put in the text below:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/backend2.yaml" >}}{{< /highlight >}}

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

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/cnp-l7-sm.yaml" >}}{{< /highlight >}}

And apply the `CiliumNetworkPolicy` with:

```bash
kubectl apply -f cnp-l7-sm.yaml
```

Until now only the backend service is replying to Ingress traffic. Now we configure Envoy to loadbalance the traffic 50/50 between backend and backend-2 with retries.
We are using a CustomResource called `CiliumEnvoyConfig` for this. Create a file `envoyconfig.yaml` with the following content:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/envoyconfig.yaml" >}}{{< /highlight >}}

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
    "body": "another secret information"                                                                                                 
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


## Task {{% param sectionnumber %}}.4: Cleanup

You can delete the Service Mesh cluster now.

```bash
minikube delete -p servicemesh
```

