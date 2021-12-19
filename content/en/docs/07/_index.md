---
title: "7. Transparent Encryption"
weight: 7
sectionnumber: 7
---
## Host traffic/endpoint traffic encryption

To secure communication inside a kubernetes cluster Cilium supports transparent encryption of traffic between Cilium-managed endpoints either using IPsec or [WireGuard®](https://www.wireguard.com/).


## Task {{% param sectionnumber %}}.1: Increase cluster size

By default minikube create single node clusters. Add a second node to the cluster:

```bash
minikube -p cluster1 node add
```


## Task {{% param sectionnumber %}}.2: Move frontend app to different node

To see traffic between nodes we move the frontend pod from Chapter 3 to the newly created node:

```yaml
spec:
  template:
    spec:
      nodeSelector:
        kubernetes.io/hostname: cluster1-m02 
```

```bash
kubectl patch deployments.apps frontend --type merge --patch "$(cat patch.yaml)"
```
We should see the frontend now running on the new node

```bash
kubectl get pods -o wide
```


## Task {{% param sectionnumber %}}.3:  Enable node traffic encryption with WireGuard

Enabling WireGuard based encryption with helm is simple:

```bash
helm upgrade -i cilium cilium/cilium \
  --namespace kube-system \
  --reuse-values \
  --set l7Proxy=false \
  --set encryption.enabled=true \
  --set encryption.type=wireguard \
  --set enryption.wireguard.userspaceFallback=true \
  --wait
kubectl -n kube-system rollout restart ds cilium
```
Currently L7 policy enforcement and visibility is [not supported](https://github.com/cilium/cilium/issues/15462) with WireGuard, this is why we have to disable it.


### Verify encryption is working


Verify the number of peers in encryption is correct (should be the sum of nodes - 1)
```bash
kubectl -n kube-system exec -ti ds/cilium -- cilium status | grep Encryption
```

You should see something similiar to this (in this example we have a two node cluster):

```bash
Encryption:             Wireguard       [cilium_wg0 (Pubkey: XbTJd5Gnp7F8cG2Ymj6q11dBx8OtP1J5ZOAhswPiYAc=, Port: 51871, Peers: 1)]
```

We can check if the traffic is really sent to the WireGuard tunnel device cilium_wg0 (hit Ctrl+C to stop sniffing).

```bash
CILIUM_AGENT=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].metadata.name}")
kubectl debug -n kube-system -i ${CILIUM_AGENT} --image=nicolaka/netshoot -- tcpdump -ni cilium_wg0
```
If you don't see any traffic, generate it yourself. Open a new terminal and call the backend service from our frontend pod.

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
kubectl exec -ti ${FRONTEND} -- curl -Is backend:8080
```
You should now see traffic flowing through the WireGuard tunnel interface cilium_wg0.

## Legal
“WireGuard” is a registered trademark of Jason A. Donenfeld.

