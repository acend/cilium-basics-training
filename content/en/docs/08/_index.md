---
title: "8. Transparent Encryption"
weight: 8
sectionnumber: 8
---
## Host traffic/endpoint traffic encryption

To secure communication inside a Kubernetes cluster Cilium supports transparent encryption of traffic between Cilium-managed endpoints either using IPsec or [WireGuard®](https://www.wireguard.com/).


## Task {{% param sectionnumber %}}.1: Increase cluster size

By default Minikube creates single-node clusters. Add a second node to the cluster:

```bash
minikube -p cluster1 node add
```


## Task {{% param sectionnumber %}}.2: Move frontend app to a different node

To see traffic between nodes, we move the frontend pod from Chapter 3 to the newly created node:

{{< highlight yaml >}}{{< readfile file="content/en/docs/08/patch.yaml" >}}{{< /highlight >}}

create a file `patch.yaml` with the above content. You can patch the frontend deployment now:

```bash
kubectl patch deployments.apps frontend --type merge --patch-file patch.yaml
```
We should see the frontend now running on the new node `cluster1-m02`:

```bash
kubectl get pods -o wide
```

```
NAME                           READY   STATUS        RESTARTS      AGE   IP           NODE           NOMINATED NODE   READINESS GATES
backend-65f7c794cc-hh6pw       1/1     Running       0             22m   10.1.0.39    cluster1       <none>           <none>
deathstar-6c94dcc57b-6chpk     1/1     Running       1 (10m ago)   11m   10.1.0.207   cluster1       <none>           <none>
deathstar-6c94dcc57b-vtt8b     1/1     Running       0             11m   10.1.0.220   cluster1       <none>           <none>
frontend-6db4b77ff6-kznfl      1/1     Running       0             35s   10.1.1.7     cluster1-m02   <none>           <none>
not-frontend-8f467ccbd-4jl6z   1/1     Running       0             22m   10.1.0.115   cluster1       <none>           <none>
tiefighter                     1/1     Running       0             11m   10.1.0.185   cluster1       <none>           <none>
xwing                          1/1     Running       0             11m   10.1.0.205   cluster1       <none>           <none>

```


## Task {{% param sectionnumber %}}.3:  Enable node traffic encryption with WireGuard

Enabling WireGuard based encryption with Helm is simple:

```bash
helm upgrade -i cilium cilium/cilium \
  --namespace kube-system \
  --reuse-values \
  --set l7Proxy=false \
  --set encryption.enabled=true \
  --set encryption.type=wireguard \
  --set enryption.wireguard.userspaceFallback=true \
  --wait
```

Afterwards restart the Cilium DaemonSet:

```bash
kubectl -n kube-system rollout restart ds cilium
```

Currently, L7 policy enforcement and visibility is [not supported](https://github.com/cilium/cilium/issues/15462) with WireGuard, this is why we have to disable it.


### Verify encryption is working


Verify the number of peers in encryption is 1 (this can take a while, the number is sum of nodes - 1)
```bash
kubectl -n kube-system exec -ti ds/cilium -- cilium status | grep Encryption
```

You should see something similar to this (in this example we have a two-node cluster):

```bash
Encryption:             Wireguard       [cilium_wg0 (Pubkey: XbTJd5Gnp7F8cG2Ymj6q11dBx8OtP1J5ZOAhswPiYAc=, Port: 51871, Peers: 1)]
```

We now check if the traffic is really sent to the WireGuard tunnel device cilium_wg0.

```bash
CILIUM_AGENT=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].metadata.name}")
kubectl debug -n kube-system -i ${CILIUM_AGENT} --image=nicolaka/netshoot -- tcpdump -ni cilium_wg0 -X port 8080
```

If you don't see any traffic, you can generate it yourself in a second terminal. For those using the Webshell a second Terminal can be opened using the menu `Terminal` then `Split Terminal`. Now in this second terminal run:

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
kubectl exec -ti ${FRONTEND} -- curl -Is backend:8080
```
You should now see traffic flowing through the WireGuard tunnel interface cilium_wg0.

You can close the second terminal with `exit`. Then hit `Ctrl+C` to stop sniffing.

{{% alert title="Note" color="primary" %}}
As we are sniffing in the WireGuard interface `cilium_wg0` you see the unencrypted traffic.
{{% /alert %}}


## Task {{% param sectionnumber %}}.4:  CleanUp

To not mess up the next ClusterMesh Lab we are going to disable WireGuard encryption again:

```bash
helm upgrade -i cilium cilium/cilium \
  --namespace kube-system \
  --reuse-values \
  --set l7Proxy=true \
  --set encryption.enabled=false \
  --wait
```

and then restart the Cilium Daemonset:

```bash
kubectl -n kube-system rollout restart ds cilium
```

Verify that it is disabled again:

```bash
kubectl -n kube-system exec -ti ds/cilium -- cilium status | grep Encryption
```

```
Encryption:                       Disabled
```


## Legal

“WireGuard” is a registered trademark of Jason A. Donenfeld
