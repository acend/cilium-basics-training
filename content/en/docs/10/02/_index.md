---
title: "10.2 Host-Reachable Services"
weight: 102
sectionnumber: 10.2
---

When running Cilium without kube-proxy, by default, services in a Kubernetes Cluster cannot be reached from the host namespace, only from the pod namespaces. This means, if you have a pod running with `hostNetwork: true` they won't be able to reach any Kubernetes Service.

{{% alert title="Note" color="primary" %}}
Host-reachable services for TCP and UDP requires a v4.19.57, v5.1.16, v5.2.0 or more recent Linux kernel. Note that v5.0.y kernels do not have the fix required to run host-reachable services with UDP since at this point in time the v5.0.y stable kernel is end-of-life (EOL) and not maintained anymore. For only enabling TCP-based host-reachable services a v4.17.0 or newer kernel is required. The most optimal kernel with the full feature set is v5.8.
{{% /alert %}}


## Task {{% param sectionnumber %}}.1: Try access Services from HostNetwork

Let us create a simple NGINX Deployemt using `hostNetwork: true` (& `dnsPolicy: ClusterFirstWithHostNet` for the Service DNS resolution to work)

Here's the `nginx.yaml` to deploy NGINX with a service:

{{< highlight yaml >}}{{< readfile file="content/en/docs/10/02/nginx.yaml" >}}{{< /highlight >}}

Apply this with:

```bash
kubectl apply -f nginx.yaml
```

and then let us try to reach our `backend` service (from Lab 3: Network Policies) with:

```bash
kubectl exec -it <nginx-podname> -- curl -v -m 2 http://backend:8080
```

and now indeed you see, the `nginx` Pod is not able to reach our `backend` Service

```
*   Trying 10.99.155.192:8080...
* Connection timed out after 2000 milliseconds
* Closing connection 0
curl: (28) Connection timed out after 2000 milliseconds
command terminated with exit code 28
```

{{% alert title="Note" color="primary" %}}
Remember in Lab 3 we created a NetworkPolicy and only Pods with a label `app: frontend` are allowed to access the `backend` pods. Therefore our `nginx` Pod also has this label.
{{% /alert %}}


## Task {{% param sectionnumber %}}.2: Enable Host-Reachable Services


```bash
kubectl patch -n kube-system cm cilium-config --patch '{"data":{"enable-host-reachable-services": "true"}}'
```

and then restart your Cilium pods:

```bash
kubectl -n kube-system rollout restart ds cilium
```


## Task {{% param sectionnumber %}}.3: Enable Host-Reachable Services

