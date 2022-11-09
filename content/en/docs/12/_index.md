---
title: "Cilium Enterprise"
weight: 12

OnlyWhenNot: techlab
---

So far, we used the Cilium CNI in the Open Source Software (OSS) version. Cilium OSS has [joined the CNCF as an incubating project](https://www.cncf.io/blog/2021/10/13/cilium-joins-cncf-as-an-incubating-project/) and only recently during KubeCon 2022 NA [applied to become a CNCF graduated project](https://cilium.io/blog/2022/10/27/cilium-applies-for-graduation/). [Isovalent](https://isovalent.com/), the company behind Cilium also offers enterprise support for the Cilium CNI. In this lab, we are going to look at some of the enterprise features.


## {{% task %}}  Create a Kubernetes Cluster and install Cilium Enterprise

We are going to spin up a new Kubernetes cluster with the following command:


```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version={{% param "kubernetesVersion" %}} -p cilium-enterprise 
```

Now check that everything is up and running:

```bash
kubectl get node
```

This should produce a similar output:

```
NAME                STATUS   ROLES                  AGE   VERSION
cilium-enterprise   Ready    control-plane,master   86s   v{{% param "kubernetesVersion" %}}
```

Alright, everything is up and running and we can continue with the Cilium Enterprise Installation. First we need to add the Helm chart repository:

```
helm repo add isovalent https://....
```

{{% alert title="Note" color="primary" %}}
Your trainer will provide you with the Helm chart url.
{{% /alert %}}


Next, create a `cilium-enterprise-values.yaml` file with the following content:

```yaml
cilium:
  hubble:
    enabled: false
    relay:
      enabled: false
  nodeinit:
    enabled: true
  ipam:
    mode: cluster-pool
hubble-enterprise:
  enabled: false
  enterprise:
    enabled: false
hubble-ui:
  enabled: false
```

And then install Cilium enterprise with Helm:

```bash
helm install cilium-enterprise isovalent/cilium-enterprise --version {{% param "ciliumVersion.enterprise" %}}  \
  --namespace kube-system -f cilium-enterprise-values.yaml
```


To confirm that the cilium daemonset is running Cilium Enterprise, execute the following command and verify that the container registry for `cilium-agent` is set to `quay.io/isovalent/cilium`:

```bash
kubectl get ds -n kube-system cilium -o jsonpath='{.spec.template.spec.containers[0].image}' | cut -d: -f1
```

Run the following command and validate that cilium daemonset is up and running:

```bash
kubectl get ds -n kube-system cilium
```

This should give you an output similar to this:

```
NAME                DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
cilium              1         1         1       1            1           <none>          91s
```

