---
title: "8.2 Network Policies"
weight: 82
sectionnumber: 8.2
---

## Task {{% param sectionnumber %}}.1: Allowing Specific Communication Between Clusters


The following policy illustrates how to allow particular pods to communicate between two clusters.

{{< highlight yaml >}}{{< readfile file="content/en/docs/08/02/cnp.yaml" >}}{{< /highlight >}}

{{% alert title="Note" color="primary" %}}
For the Pods to resolve the `rebel-base` service name they still need connectivity to Kubernetes DNS Service. Therefore access to that is also allowed.
{{% /alert %}}

Kubernetes security policies are not automatically distributed across clusters, it is your responsibility to apply CiliumNetworkPolicy or NetworkPolicy in all clusters.

So let us apply the above `CiliumNetworkPolicy` to both clusters:

```bash
kubectl --context cluster1 apply -f cnp.yaml
kubectl --context cluster2 apply -f cnp.yaml
```

Let us run our `curl` `for` loop again:

```bash
XWINGPOD=<x-wing-podname>
for i in {1..10}; do                                       
  kubectl --context cluster1 exec -it $XWINGPOD -- curl -m 1rebel-base
done
```

and as an result you see:

```c
url: (28) Connection timed out after 1001 milliseconds
command terminated with exit code 28
curl: (28) Connection timed out after 1000 milliseconds
command terminated with exit code 28
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
curl: (28) Connection timed out after 1000 milliseconds
command terminated with exit code 28
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
curl: (28) Connection timed out after 1000 milliseconds
command terminated with exit code 28
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
curl: (28) Connection timed out after 1000 milliseconds
command terminated with exit code 28
```

all connections to `cluster2` are dropped while the ones to `cluster1` are still working.