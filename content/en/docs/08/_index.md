---
title: "8. Cluster Mesh"
weight: 8
sectionnumber: 8
---


## Task {{% param sectionnumber %}}.1: Create a second Kubernetes Cluster

In order to create a Cluster Mesh we need a second Kubernetes Cluster. For the Cluster Mesh to work, the PodCIDR ranges in all clusters and all nodes must be non-conflicting and unique IP addresses. The Nodes in all clusters must have IP connectivity between each other and the network between the clusters must allow the inter-cluster communication.

{{% alert title="Note" color="primary" %}}
The exact ports are documented in the [Firewall Rules](https://docs.cilium.io/en/v1.11/operations/system_requirements/#firewall-requirements) section.
{{% /alert %}}

To start a second cluster run the following command:

```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version=1.21.6 -p cluster2
```

and the install Cilium using the `cilium` CLI. Remember, we need a different PodCIDR for the second cluster, therefore while installing Cilium, we have to change this config:

```bash
cilium install --config cluster-pool-ipv4-cidr=172.16.0.0/20
```

Then wait until the Cluster and Cilium is ready.

```
cilium status                                               
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         disabled
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

DaemonSet         cilium             Desired: 1, Ready: 1/1, Available: 1/1
Deployment        cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
Containers:       cilium-operator    Running: 1
                  cilium             Running: 1
Cluster Pods:     1/1 managed by Cilium
Image versions    cilium             quay.io/cilium/cilium:v1.11.0: 1
                  cilium-operator    quay.io/cilium/operator-generic:v1.11.0: 1
```

You can verify the correct podCidr using:

```bash
kubectl get pod -A -o wide                 

```

Have a look at the `codedns-` Pod and verify that it's IP is from your defined `172.16.0.0/24` range.

```
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE     IP             NODE       NOMINATED NODE   READINESS GATES
kube-system   cilium-htngj                       1/1     Running   0          78s     192.168.58.2   cluster2   <none>           <none>
kube-system   cilium-operator-776958f5bb-558hr   1/1     Running   0          78s     192.168.58.2   cluster2   <none>           <none>
kube-system   coredns-558bd4d5db-m9pdb           1/1     Running   0          3m59s   172.16.0.131   cluster2   <none>           <none>
kube-system   etcd-cluster2                      1/1     Running   0          4m7s    192.168.58.2   cluster2   <none>           <none>
kube-system   kube-apiserver-cluster2            1/1     Running   0          4m15s   192.168.58.2   cluster2   <none>           <none>
kube-system   kube-controller-manager-cluster2   1/1     Running   0          4m7s    192.168.58.2   cluster2   <none>           <none>
kube-system   kube-proxy-cgbbt                   1/1     Running   0          3m59s   192.168.58.2   cluster2   <none>           <none>
kube-system   kube-scheduler-cluster2            1/1     Running   0          4m7s    192.168.58.2   cluster2   <none>           <none>
kube-system   storage-provisioner                1/1     Running   1          4m12s   192.168.58.2   cluster2   <none>           <none>
```

Great the second cluster and Cilium is ready to use.


## Task {{% param sectionnumber %}}.2: Enable Cluster Mesh on both Cluster
