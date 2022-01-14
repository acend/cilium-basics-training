---
title: "11. Cilium Service Mesh"
weight: 11
sectionnumber: 11
---
Cilium Service Mesh enables functions like Ingress or Layer 7 Loadbalancing, Let us test it out.


## Task {{% param sectionnumber %}}.1: Installation

Cilium Service Mesh is still in Beta, if you want more information about the current statut you can find it [here](https://github.com/cilium/cilium-service-mesh-beta). The Beta version uses specific images, because of that we will use a dedicated cluster and install Cilium with the CLI.

{{% alert title="Note" color="primary" %}}
You can stop cluster1 with `minkube stop -p cluster1` to free up resources and speed up things.
{{% /alert %}}


```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version=1.23.0 -p servicemesh
cilium install --version -service-mesh:v1.11.0-beta.1 --config enable-envoy-config=true --kube-proxy-replacement=probe
```

Wait until cilium is ready (check with `cilium status`)

and then also enable the Hubble UI:

```bash
cilium hubble enable --ui 
```


## Task {{% param sectionnumber %}}.2: Create Ingress

Cilium Service Mesh can handle Ingress traffic with its envoy proxy.


We deploy [the sample app from chapter 3](https://cilium-basics.training.acend.ch/docs/03/#task-32-deploy-simple-application).


{{< highlight yaml >}}{{< readfile file="content/en/docs/03/simple-app.yaml" >}}{{< /highlight >}}

Apply this with:

```bash
kubectl create -f simple-app.yaml
```

Now we add an Ingress resource:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/ingress.yaml" >}}{{< /highlight >}}

Apply again with:

```bash
kubectl apply -f ingress.yaml
```

Check the ingress and the service:

```bash
kubectl describe ingress
kubectl get svc cilium-ingress-backend
```
We have successfully created an Ingress service. Unfortunately minikube has no loadbalancer deployer and we will not get a public IP for our service (status is pending)

As a workaround we can test the service from inside kubernetes.

```bash
SERVICE_IP=$(kubectl get svc cilium-ingress-backend -ojsonpath={.spec.clusterIP})
kubectl run --rm=true -it --restart=Never --image=curlimages/curl -- curl http://${SERVICE_IP}/public
```

And you get the following output:

```
[
  {
    "id": 1,
    "body": "public information"
  }
]pod "curl" deleted
```
{{% alert title="Note" color="primary" %}}
We can also use `minkube tunnel -p servicemesh` and then curl the Cluster-IP directly from our browser.
{{% /alert %}}


## Task {{% param sectionnumber %}}.3: Layer 7 Loadbalancing

Ingress is not really a Service Mesh feature. Let us test out a traffic control example by loadbalancing a service inside the proxy.


Start by creating the second service:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/backend2.yaml" >}}{{< /highlight >}}

And apply it:
```bash
kubectl apply -f backend2.yaml
```

Call it
```bash
kubectl run --rm=true -it --restart=Never --image=curlimages/curl -- curl --connect-timeout 3 http://backend-2:8080/public
```
We see an output very similiar to our simple application backend, but with a changed text.

Layer 7 loadbalancing will need to be routed through the proxy, we will enable this for our backend pods using a Cilium Network Policy with HTTP rules. We will block access to /public and allow requests to /private:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/cnp-l7.yaml" >}}{{< /highlight >}}

Apply the CiliumNetwork Policy with:

```bash
kubectl apply -f cnp-l7.yaml
```

Until now only backend service is repling to Ingress traffic. Now we configure envoy to loadbalance the traffic 50/50 between backend and backend-2 with retries.
We are using a CustomResource called `CiliumEnvoyConfig` for this:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/envoyconfig.yaml" >}}{{< /highlight >}}

{{% alert title="Note" color="primary" %}}
If you want to read more about envoy configuration [envoy Architectural Overview](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/http/http) is a good place.
{{% /alert %}}

Apply the CiliumEnvoyConfig with:

```bash
kubectl apply -f envoyconfig.yaml
```

Test it by running `curl` a few times...different backends should respond.

```bash
for i in {1..5}; do
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

You can delete the service mesh cluster now.

```bash
minkube delete -p servicemesh
```

