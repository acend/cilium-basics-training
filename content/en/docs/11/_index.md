---
title: "11. Cilium Service Mesh"
weight: 11
sectionnumber: 11
---
NOTE: Cilium Service Mesh is still in Beta you will find the relevante information [here](https://github.com/cilium/cilium-service-mesh-beta).


## Task {{% param sectionnumber %}}.1: Installation

As the Cilium Service Mesh is still in Beta and uses specific images we will use a dedicated Cluster for it.

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

We deploy [the sample app from chapter 3](https://cilium-basics.training.acend.ch/docs/03/#task-31-deploy-simple-application).


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

As we see the ingress is of Type LoadBalancer and did not get an External-IP. This is because we do not have an LoadBalancer deployed with minikube.
As a workaround we can test the service from inside kubernetes.


```bash
SERVICE_IP=$(kubectl get svc cilium-ingress-backend -ojsonpath={.spec.clusterIP})
kubectl run --rm=true -it --restart=Never --image=docker.io/byrnedo/alpine-curl:0.1.8 -- curl http://${SERVICE_IP}/public
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


## Task {{% param sectionnumber %}}.3: Layer 7 Loadbalancing

For loadbalancing we need another service:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/backend2.yaml" >}}{{< /highlight >}}

And apply this again with:


```bash
kubectl apply -f backend2.yaml
```

Layer 7 loadbalancing will need to be routed through the proxy, we will enable this for our backend pods using a Cilium Network Policy:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/cnp-l7.yaml" >}}{{< /highlight >}}

Apply the CiliumNetwork Policy with:

```bash
kubectl apply -f cnp-l7.yaml
```

Now we configure envoy to LB 50/50 backend services with retries. We are using a CustomResource called `CiliumEnvoyConfig` for this:

{{< highlight yaml >}}{{< readfile file="content/en/docs/11/envoyconfig.yaml" >}}{{< /highlight >}}

Apply the CiliumEnvoyConfig with:

```bash
kubectl apply -f envoyconfig.yaml
```

Test it by running `curl` a few times...different backends should respond.

```bash
for i in {1..5}; do
  kubectl run --rm=true -it --image=curlimages/curl --restart=Never curl -- curl  http://backend-2:8080/private
done
```


// TODO: show expected Output
