---
title: "9.1 Load-balancing with Global Services"
weight: 91
sectionnumber: 9.1
---

This lab will guide you to perform load-balancing and service discovery across multiple Kubernetes clusters.


## Task {{% param sectionnumber %}}.1: Load-balancing with Global Services

Establishing load-balancing between clusters is achieved by defining a Kubernetes service with an identical name and namespace in each cluster and adding the `annotation io.cilium/global-service: "true"` to declare it global. Cilium will automatically perform load-balancing to pods in both clusters.

We are going to deploy a global service and a sample application on both of our connected clusters.

First the Kubernetes service:

{{< highlight yaml >}}{{< readfile file="content/en/docs/09/01/svc.yaml" >}}{{< /highlight >}}

Apply this with:

```bash
kubectl --context cluster1 apply -f svc.yaml
kubectl --context cluster2 apply -f svc.yaml
```

Then deploy our sample application on both clusters.

`cluster1.yaml`:

{{< highlight yaml >}}{{< readfile file="content/en/docs/09/01/cluster1.yaml" >}}{{< /highlight >}}

```bash
kubectl --context cluster1 apply -f cluster1.yaml
```

`cluster2.yaml`:

{{< highlight yaml >}}{{< readfile file="content/en/docs/09/01/cluster2.yaml" >}}{{< /highlight >}}

```bash
kubectl --context cluster2 apply -f cluster2.yaml
```

Now you can execute from either cluster the following command (there are two x-wing pods, simply select one):

```bash
XWINGPOD=$(kubectl --context cluster1 get pod -l name=x-wing -o jsonpath="{.items[0].metadata.name}")
for i in {1..10}; do                                       
  kubectl --context cluster1  exec -it $XWINGPOD -- curl -m 1 rebel-base
done
```

as a Result you get the following output:

```
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-1"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
```

and as you see, you get results from both clusters. Even ig you scale down your `rebel-base` Deployment on `cluster1` with:

```bash
kubectl --context cluster1 scale deployment rebel-base --replicas=0
```

and then execute the `curl` `for` loop again, you still get answers, this time only from `cluster2`:

```
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}
{"Galaxy": "Alderaan", "Cluster": "Cluster-2"}

```


## Task {{% param sectionnumber %}}.2: Load-balancing only to a Remote Cluster

By default, a Global Service will load-balance across backends in multiple clusters. This implicitly configures `io.cilium/shared-service: "true"`. To prevent service backends from being shared to other clusters, and to ensure that the service will only load-balance to backends in remote clusters, this option should be disabled.

So let us change our `rebel-base` service on `cluster1` with the following updated service definition:

{{< highlight yaml >}}{{< readfile file="content/en/docs/09/01/svc2.yaml" >}}{{< /highlight >}}

Apply this with:

```bash
kubectl --context cluster1 apply -f svc2.yaml
```

// TODO? What should be the result? I still see answers from both cluster...
