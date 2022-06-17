---
title: "5. Troubleshooting"
weight: 5
sectionnumber: 5
OnlyWhenNot: techlab
---


For more details on Troubleshooting, have a look into [Cilium's Troubleshooting Documentation](https://docs.cilium.io/en/stable/operations/troubleshooting/).


## Component & Cluster Health

An initial overview of Cilium can be retrieved by listing all pods to verify whether all pods have the status `Running`:

```bash
kubectl -n kube-system get pods -l k8s-app=cilium
```

```
NAME           READY     STATUS    RESTARTS   AGE
cilium-2hq5z   1/1       Running   0          4d
cilium-6kbtz   1/1       Running   0          4d
cilium-klj4b   1/1       Running   0          4d
cilium-zmjj9   1/1       Running   0          4d
```

If Cilium encounters a problem that it cannot recover from, it will automatically report the failure state via `cilium status` which is regularly queried by the Kubernetes liveness probe to automatically restart Cilium pods. If a Cilium Pod is in state `CrashLoopBackoff` then this indicates a permanent failure scenario.

If a particular Cilium Pod is not in a running state, the status and health of the agent on that node can be retrieved by running `cilium status` in the context of that pod:

```bash
kubectl -n kube-system exec ds/cilium -- cilium status
```

The output looks similar to this:

```
Defaulted container "cilium-agent" out of: cilium-agent, mount-cgroup (init), clean-cilium-state (init)
KVStore:                Ok   Disabled
Kubernetes:             Ok   1.23 (v1.23.0) [linux/amd64]
Kubernetes APIs:        ["cilium/v2::CiliumClusterwideNetworkPolicy", "cilium/v2::CiliumEndpoint", "cilium/v2::CiliumNetworkPolicy", "cilium/v2::CiliumNode", "core/v1::Namespace", "core/v1::Node", "core/v1::Pods", "core/v1::Service", "discovery/v1::EndpointSlice", "networking.k8s.io/v1::NetworkPolicy"]
KubeProxyReplacement:   Disabled   
Host firewall:          Disabled
Cilium:                 Ok   1.11.0 (v1.11.0-27e0848)
NodeMonitor:            Listening for events on 32 CPUs with 64x4096 of shared memory
Cilium health daemon:   Ok   
IPAM:                   IPv4: 11/254 allocated from 10.1.0.0/24, 
ClusterMesh:            0/0 clusters ready, 0 global-services
BandwidthManager:       Disabled
Host Routing:           Legacy
Masquerading:           IPTables [IPv4: Enabled, IPv6: Disabled]
Controller Status:      56/56 healthy
Proxy Status:           OK, ip 10.1.0.145, 2 redirects active on ports 10000-20000
Hubble:                 Ok   Current/Max Flows: 4095/4095 (100.00%), Flows/s: 7.17   Metrics: Disabled
Encryption:             Disabled
Cluster health:         1/1 reachable   (2022-01-10T12:29:11Z)

```

More detailed information about the status of Cilium can be inspected with:


```bash
kubectl -n kube-system exec ds/cilium -- cilium status --verbose
```

Verbose output includes detailed IPAM state (allocated addresses), Cilium controller status, and details of the Proxy status.


## Logs

To retrieve log files of a cilium pod, run:

```bash
kubectl -n kube-system logs --timestamps <pod-name>
```

The `<pod-name>` can be determined with the following command and by selecting the name of one of the pods:

```bash
kubectl -n kube-system get pods -l k8s-app=cilium
```

If the Cilium Pod was already restarted due to the liveness problem after encountering an issue, it can be useful to retrieve the logs of the Pod previous to the last restart:

```bash
kubectl -n kube-system logs --timestamps -p <pod-name>
```


## Policy Troubleshooting - Ensure Pod is managed by Cilium

A potential cause for policy enforcement not functioning as expected is that the networking of the Pod selected by the policy is not being managed by Cilium. The following situations result in unmanaged pods:

* The Pod is running in host networking and will use the host’s IP address directly. Such pods have full network connectivity but Cilium will not provide security policy enforcement for such pods.
* The Pod was started before Cilium was deployed. Cilium only manages pods that have been deployed after Cilium itself was started. Cilium will not provide security policy enforcement for such pods.

If Pod networking is not managed by Cilium, ingress and egress policy rules selecting the respective pods will not be applied. See the section Network Policy for more details.

For a quick assessment of whether any pods are not managed by Cilium, the Cilium CLI will print the number of managed pods. If this prints that all of the pods are managed by Cilium, then there is no problem:

```bash
cilium status
```

```
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         OK
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

Deployment        cilium-operator    Desired: 2, Ready: 2/2, Available: 2/2
Deployment        hubble-relay       Desired: 1, Ready: 1/1, Available: 1/1
Deployment        hubble-ui          Desired: 1, Ready: 1/1, Available: 1/1
DaemonSet         cilium             Desired: 2, Ready: 2/2, Available: 2/2
Containers:       cilium-operator    Running: 2
                  hubble-relay       Running: 1
                  hubble-ui          Running: 1
                  cilium             Running: 2
Cluster Pods:     5/5 managed by Cilium
```

You can run the following script to list the pods which are not managed by Cilium:

```bash
curl -sLO https://raw.githubusercontent.com/cilium/cilium/master/contrib/k8s/k8s-unmanaged.sh
chmod +x k8s-unmanaged.sh
./k8s-unmanaged.sh
```

```
kube-system/cilium-hqpk7
kube-system/kube-addon-manager-minikube
kube-system/kube-dns-54cccfbdf8-zmv2c
kube-system/kubernetes-dashboard-77d8b98585-g52k5
kube-system/storage-provisioner
```

{{% alert title="Note" color="primary" %}}
It's ok if you don't see any Pods listed with the above command. We don't have any unmanaged Pods in our setup.
{{% /alert %}}


## Reporting a problem - Automatic log & state collection

Before you report a problem, make sure to retrieve the necessary information from your cluster before the failure state is lost.

Execute the `cilium sysdump` command to collect troubleshooting information from your Kubernetes cluster:

```bash
cilium sysdump
```

Note that by default `cilium sysdump` will attempt to collect as many logs as possible for all the nodes in the cluster. If your cluster size is above 20 nodes, consider setting the following options to limit the size of the sysdump. This is not required, but is useful for those who have a constraint on bandwidth or upload size.

* set the `--node-list` option to pick only a few nodes in case the cluster has many of them.
* set the `--logs-since-time` option to go back in time to when the issues started.
* set the `--logs-limit-bytes` option to limit the size of the log files (note: passed onto kubectl logs; does not apply to entire collection archive).
Ideally, a sysdump that has a full history of select nodes, rather than a brief history of all the nodes, would be preferred (by using `--node-list`). The second recommended way would be to use `--logs-since-time` if you are able to narrow down when the issues started. Lastly, if the Cilium agent and Operator logs are too large, consider `--logs-limit-bytes`.

Use `--help` to see more options:

```bash
cilium sysdump --help
```
