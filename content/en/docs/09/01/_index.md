---
title: "Enable Cluster Mesh"
weight: 91
OnlyWhenNot: techlab
---


## {{% task %}} Create a second Kubernetes Cluster

To create a Cluster Mesh, we need a second Kubernetes cluster. For the Cluster Mesh to work, the PodCIDR ranges in all clusters and nodes must be non-conflicting and have unique IP addresses. The nodes in all clusters must have IP connectivity between each other and the network between the clusters must allow inter-cluster communication.

{{% alert title="Note" color="primary" %}}
The exact ports are documented in the [Firewall Rules](https://docs.cilium.io/en/v1.12/operations/system_requirements/#firewall-requirements) section.
{{% /alert %}}

To start a second cluster run the following command:

```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version={{% param "kubernetesVersion" %}} -p cluster2
```

As Minikube with the Docker driver uses separated Docker networks, we need to make sure that your system forwards traffic between the two networks. To enable forwarding by default execute:

```bash
sudo iptables -I DOCKER-USER -j ACCEPT
```

Then install Cilium using Helm. Remember, we need a different PodCIDR for the second cluster, therefore while installing Cilium, we have to change this config:

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDRList={10.2.0.0/16} \
  --set cluster.name=cluster2 \
  --set cluster.id=2 \
  --set operator.replicas=1 \
  --set kubeProxyReplacement=disabled \
  --wait
```

Then wait until the cluster and Cilium is ready.

```bash
cilium status --wait
```

```
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
Image versions    cilium             quay.io/cilium/cilium:v{{% param "ciliumVersion.postUpgrade" %}}: 1
                  cilium-operator    quay.io/cilium/operator-generic:v{{% param "ciliumVersion.postUpgrade" %}}: 1
```

You can verify the correct PodCIDR using:

```bash
kubectl get pod -A -o wide
```

Have a look at the `coredns-` Pod and verify that it's IP is from your defined `10.2.0.0/16` range.

```
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE   IP             NODE       NOMINATED NODE   READINESS GATES
kube-system   cilium-operator-776958f5bb-m5hww   1/1     Running   0          29s   192.168.58.2   cluster2   <none>           <none>
kube-system   cilium-qg9xj                       1/1     Running   0          29s   192.168.58.2   cluster2   <none>           <none>
kube-system   coredns-558bd4d5db-z6cxh           1/1     Running   0          38s   10.2.0.240     cluster2   <none>           <none>
kube-system   etcd-cluster2                      1/1     Running   0          44s   192.168.58.2   cluster2   <none>           <none>
kube-system   kube-apiserver-cluster2            1/1     Running   0          44s   192.168.58.2   cluster2   <none>           <none>
kube-system   kube-controller-manager-cluster2   1/1     Running   0          44s   192.168.58.2   cluster2   <none>           <none>
kube-system   kube-proxy-bqk4r                   1/1     Running   0          38s   192.168.58.2   cluster2   <none>           <none>
kube-system   kube-scheduler-cluster2            1/1     Running   0          44s   192.168.58.2   cluster2   <none>           <none>
kube-system   storage-provisioner                1/1     Running   1          49s   192.168.58.2   cluster2   <none>           <none>
```

The second cluster and Cilium is ready to use.


## {{% task %}} Enable Cluster Mesh on both Cluster

Now let us enable the Cluster Mesh using the `cilium` CLI on both clusters:


{{% alert title="Note" color="primary" %}}
Although so far we used Helm to install and update Cilium, enabling Cilium Service Mesh using Helm is currently [undocumented](https://github.com/cilium/cilium/issues/19057). We make an exception from the rule to never mix Helm and CLI installations and do it with the CLI.
{{% /alert %}}

```bash
cilium clustermesh enable --context cluster1 --service-type NodePort
cilium clustermesh enable --context cluster2 --service-type NodePort
```

You can now verify the Cluster Mesh status using:

```bash
cilium clustermesh status --context cluster1 --wait
```

```
⚠️  Service type NodePort detected! Service may fail when nodes are removed from the cluster!
✅ Cluster access information is available:
  - 192.168.49.2:31839
✅ Service "clustermesh-apiserver" of type "NodePort" found
⌛ [cluster1] Waiting for deployment clustermesh-apiserver to become ready...
🔌 Cluster Connections:
🔀 Global services: [ min:0 / avg:0.0 / max:0 ]
```

To connect the two clusters, the following step needs to be done in one direction only. The connection will automatically be established in both directions:

```bash
cilium clustermesh connect --context cluster1 --destination-context cluster2
```

The output should look something like this:

```
✨ Extracting access information of cluster cluster2...
🔑 Extracting secrets from cluster cluster2...
⚠️  Service type NodePort detected! Service may fail when nodes are removed from the cluster!
ℹ️  Found ClusterMesh service IPs: [192.168.58.2]
✨ Extracting access information of cluster cluster1...
🔑 Extracting secrets from cluster cluster1...
⚠️  Service type NodePort detected! Service may fail when nodes are removed from the cluster!
ℹ️  Found ClusterMesh service IPs: [192.168.49.2]
✨ Connecting cluster cluster1 -> cluster2...
🔑 Secret cilium-clustermesh does not exist yet, creating it...
🔑 Patching existing secret cilium-clustermesh...
✨ Patching DaemonSet with IP aliases cilium-clustermesh...
✨ Connecting cluster cluster2 -> cluster1...
🔑 Secret cilium-clustermesh does not exist yet, creating it...
🔑 Patching existing secret ciliugm-clustermesh...
✨ Patching DaemonSet with IP aliases cilium-clustermesh...
✅ Connected cluster cluster1 and cluster2!
```

It may take a bit for the clusters to be connected. You can execute the following command

```bash
cilium clustermesh status --context cluster1 --wait
```
to wait for the connection to be successful. The output should be:

```
⚠️  Service type NodePort detected! Service may fail when nodes are removed from the cluster!
✅ Cluster access information is available:
  - 192.168.58.2:32117
✅ Service "clustermesh-apiserver" of type "NodePort" found
⌛ [cluster2] Waiting for deployment clustermesh-apiserver to become ready...
✅ All 1 nodes are connected to all clusters [min:1 / avg:1.0 / max:1]
🔌 Cluster Connections:
- cluster1: 1/1 configured, 1/1 connected
🔀 Global services: [ min:3 / avg:3.0 / max:3 ]
```

The two clusters are now connected.


## {{% task %}} Cluster Mesh Troubleshooting

Use the following list of steps to troubleshoot issues with Cluster Mesh:

```bash
cilium status --context cluster1
```

or

```bash
cilium status --context cluster2
```

which gives you an output similar to this:

```
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         OK
 \__/¯¯\__/    ClusterMesh:    OK
    \__/

DaemonSet         cilium                   Desired: 1, Ready: 1/1, Available: 1/1
Deployment        cilium-operator          Desired: 1, Ready: 1/1, Available: 1/1
Deployment        hubble-relay             Desired: 1, Ready: 1/1, Available: 1/1
Deployment        clustermesh-apiserver    Desired: 1, Ready: 1/1, Available: 1/1
Containers:       cilium                   Running: 1
                  cilium-operator          Running: 1
                  hubble-relay             Running: 1
                  clustermesh-apiserver    Running: 1
Cluster Pods:     6/6 managed by Cilium
Image versions    cilium                   quay.io/cilium/cilium:v{{% param "ciliumVersion.postUpgrade" %}}: 1
                  cilium-operator          quay.io/cilium/operator-generic:v{{% param "ciliumVersion.postUpgrade" %}}: 1
                  hubble-relay             quay.io/cilium/hubble-relay:v{{% param "ciliumVersion.postUpgrade" %}}: 1
                  clustermesh-apiserver    quay.io/coreos/etcd:v3.4.13: 1
                  clustermesh-apiserver    quay.io/cilium/clustermesh-apiserver:v{{% param "ciliumVersion.postUpgrade" %}}: 1

```


If you cannot resolve the issue with the above commands, follow the steps in [Cilium's Cluster Mesh Troubleshooting Guide](https://docs.cilium.io/en/v1.12/operations/troubleshooting/#troubleshooting-clustermesh).
