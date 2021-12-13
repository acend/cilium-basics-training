---
title: "7. Transparent Encryption"
weight: 7
sectionnumber: 7
---
## Host traffic/endpoint traffic encryption

Cilium supports the transparent encryption of Cilium-managed host traffic and traffic between Cilium-managed endpoints either using IPsec or [WireGuard®](https://www.wireguard.com/).


### WireGuard Encryption


### Enable traffic encryption with WireGuard

```bash
kubectl patch -n kube-system cm cilium-config --patch '{"data":{"enable-wireguard": "true", "enable-l7-proxy": "false", "enable-wireguard-userspace-fallback": "true"}}'
kubectl -n kube-system rollout restart daemonset cilium
kubectl -n kube-system rollout status daemonset cilium
```


### Verify encryption is working


Add a second node in case you only have a 1 node cluster

```bash
minikube node add
```

Now verify the number of peers in encryption is correct (should be the sum of nodes - 1)
```bash
kubectl -n kube-system exec -ti ds/cilium -- cilium status | grep Encryption
```

You should see something similiar to this:

```bash
Encryption:             Wireguard       [cilium_wg0 (Pubkey: XbTJd5Gnp7F8cG2Ymj6q11dBx8OtP1J5ZOAhswPiYAc=, Port: 51871, Peers: 1)]
```
