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

Create a file `patch.yaml` with the follwing content_

{{< highlight yaml >}}{{< readfile file="content/en/docs/08/patch.yaml" >}}{{< /highlight >}}

You can patch the frontend deployment now:

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


## Task {{% param sectionnumber %}}.3:  Sniff traffic between nodes

To check if we see unencrypted traffic between nodes we will use tcpdump.
Let us filter on the host interfce for all packets containing the string `password`:

```bash
CILIUM_AGENT=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].metadata.name}")
kubectl debug -n kube-system -i ${CILIUM_AGENT} --image=nicolaka/netshoot -- tcpdump -ni eth0 -vv | grep password
```

In a second terminal we will call our backend service with a password. For those using the Webshell a second Terminal can be opened using the menu `Terminal` then `Split Terminal`, also don't forget to ssh into the VM again. Now in this second terminal run:

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
for i in {1..10}; do
 kubectl exec -ti ${FRONTEND} -- curl -Is backend:8080?password=secret
done
```

You should now see our string `password` sniffed in the network traffic. Hit `Ctrl+C` to stop sniffing but keep the second terminal open.


## Task {{% param sectionnumber %}}.4:  Enable node traffic encryption with WireGuard

Enabling WireGuard based encryption with Helm is simple:

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace kube-system \
  --reuse-values \
  --set l7Proxy=false \
  --set encryption.enabled=true \
  --set encryption.type=wireguard \
  --set enryption.wireguard.userspaceFallback=true \
  --wait
```

Afterwards restart the Cilium DaemonSet (ignore the deprecation warnings):

```bash
kubectl -n kube-system rollout restart ds cilium
```

{{% alert title="Note" color="primary" %}}
You will see some warnings in this command. You can ignore them.
{{% /alert %}}


Currently, L7 policy enforcement and visibility is [not supported](https://github.com/cilium/cilium/issues/15462) with WireGuard, this is why we have to disable it.


## Task {{% param sectionnumber %}}.5:  Verify encryption is working


Verify the number of peers in encryption is 1 (this can take a while, the number is sum of nodes - 1)
```bash
kubectl -n kube-system exec -ti ds/cilium -- cilium status | grep Encryption
```

You should see something similar to this (in this example we have a two-node cluster):

```
Encryption:             Wireguard       [cilium_wg0 (Pubkey: XbTJd5Gnp7F8cG2Ymj6q11dBx8OtP1J5ZOAhswPiYAc=, Port: 51871, Peers: 1)]
```

We now check if the traffic is really encrypted, we start sniffing again:

```bash
CILIUM_AGENT=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].metadata.name}")
kubectl debug -n kube-system -i ${CILIUM_AGENT} --image=nicolaka/netshoot -- tcpdump -ni eth0 -vv | grep password
```

Now in the other terminal generate traffic:

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
for i in {1..10}; do
 kubectl exec -ti ${FRONTEND} -- curl -Is backend:8080?password=secret
done
```
As you should see the traffic is encrypted now and we can't find our string anymore in plaintext on eth0. To sniff the traffic before it is encrypted replace the interface `eth0` with the WireGuard interface `cilium_wg0`.

Hit `Ctrl+C` to stop sniffing. You can close the second terminal with `exit`.


## Task {{% param sectionnumber %}}.6: CleanUp

To not mess up the next ClusterMesh Lab we are going to disable WireGuard encryption again:

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}}\
  --namespace kube-system \
  --reuse-values \
  --set l7Proxy=true \
  --set encryption.enabled=false \
  --wait
```

and then restart the Cilium Daemonset (ignore the deprecation warnings):

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
