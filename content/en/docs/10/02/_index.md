---
title: "10.2 Host Firewall"
weight: 102
sectionnumber: 10.2
---

Cilium is also capable to act as an host firwall to enforce seucrity policies for Kubernetes nodes. In this lab we are going to show you briefly how this works.


## Task {{% param sectionnumber %}}.1: Enable the Host Firewall in Cilium

We need to enable the host firewall in the cilium config. This can be done using Helm:


```bash
helm upgrade -i cilium cilium/cilium \
  --namespace kube-system \
  --reuse-values \
  --set hostFirewall.enabled=true          \
  --set devices='{eth0}'
  --wait
```

The devices flag refers to the network devices Cilium is configured on such as eth0. Omitting this option leads Cilium to auto-detect what interfaces the host firewall applies to.

At this point, the Cilium-managed nodes are ready to enforce network policies.


## Task {{% param sectionnumber %}}.2: Attach a Label to the Node

In this lab, we will apply host policies only to nodes with the label `node-access=ssh`. We thus first need to attach that label to a node in the cluster.

```bash
kubectl label node cluster1 node-access=ssh
```


## Task {{% param sectionnumber %}}.2: Enable Policy Audit Mode for the Host Endpoint

[Host Policies](https://docs.cilium.io/en/latest/policy/language/#hostpolicies) enforce access control over connectivity to and from nodes. Particular care must be taken to ensure that when host policies are imported, Cilium does not block access to the nodes or break the cluster’s normal behavior (for example by blocking communication with kube-apiserver).

To avoid such issues, we can switch the host firewall in audit mode, to validate the impact of host policies before enforcing them. When Policy Audit Mode is enabled, no network policy is enforced so this setting is not recommended for production deployment.

```bash
CILIUM_NAMESPACE=kube-system
CILIUM_POD_NAME=$(kubectl -n $CILIUM_NAMESPACE get pods -l "k8s-app=cilium" -o jsonpath="{.items[?(@.spec.nodeName=='$NODE_NAME')].metadata.name}")
HOST_EP_ID=$(kubectl -n $CILIUM_NAMESPACE exec $CILIUM_POD_NAME -- cilium endpoint list -o jsonpath='{[?(@.status.identity.id==1)].id}')
kubectl -n $CILIUM_NAMESPACE exec $CILIUM_POD_NAME -- cilium endpoint config $HOST_EP_ID PolicyAuditMode=Enabled
```

and then verify with:

```bash
kubectl -n $CILIUM_NAMESPACE exec $CILIUM_POD_NAME -- cilium endpoint config $HOST_EP_ID | grep PolicyAuditMode
```

which should give you:

```
PolicyAuditMode          Enabled
```


## Task {{% param sectionnumber %}}.3: Apply a Host Network Policy

Host Policies match on node labels using a Node Selector to identify the nodes to which the policy applies. The following policy applies to all nodes. It allows communications from outside the cluster only on port TCP/22. All communications from the cluster to the hosts are allowed.

Host policies don’t apply to communications between pods or between pods and the outside of the cluster, except if those pods are host-networking pods.

{{< highlight yaml >}}{{< readfile file="content/en/docs/10/02/ccwnp.yaml" >}}{{< /highlight >}}

You can apply this `CiliumClusterwideNetworkPolicy` with:

```bash
kubectl apply -f ccwnp.yaml
```

The host is represented as a special endpoint, with label `reserved:host`, in the output of command `cilium endpoint list`. You can therefore inspect the status of the policy using that command.

```bash
kubectl -n kube-system exec $(kubectl -n kube-system get pods -l k8s-app=cilium -o jsonpath='{.items[0].metadata.name}') -- cilium endpoint list
```

As long as the host endpoint is running in audit mode, communications disallowed by the policy won’t be dropped. They will however be reported by `cilium monitor` as `action audit`. The audit mode thus allows you to adjust the host policy to your environment, to avoid unexpected connection breakages.

```bash
kubectl -n kube-system exec $(kubectl -n kube-system get pods -l k8s-app=cilium -o jsonpath='{.items[0].metadata.name}') -- cilium monitor -t policy-verdict --related-to $HOST_EP_ID
```


## Task {{% param sectionnumber %}}.4: Clean Up

Once you are confident all required communication to the host from outside the cluster are allowed, you can disable policy audit mode to enforce the host policy.

{{% alert title="Note" color="primary" %}}
When enforce the host policy, make sure that none of the communications required to access the cluster or for the cluster to work properly are denied. They should appear as `action allow`.
{{% /alert %}}

We are not going to do this extended task (as it would require some more rules for the cluster to continue working). But the command to disable the audit mode looks like this:

```bash
kubectl -n $CILIUM_NAMESPACE exec $CILIUM_POD_NAME -- cilium endpoint config $HOST_EP_ID PolicyAuditMode=Enabled
```

Simply cleanup and continue:

```bash
kubectl delete ccnp demo-host-policy
kubectl label node $NODE_NAME node-access-
```
