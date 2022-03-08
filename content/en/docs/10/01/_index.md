---
title: "10.1 Host Firewall"
weight: 101
sectionnumber: 10.1
---


{{% alert title="Note" color="primary" %}}
This lab should be done on your `cluster1`, make sure to switch to `cluster1` with `minikube profile cluster1`
{{% /alert %}}

Cilium is also capable to act as a host firewall to enforce security policies for Kubernetes nodes. In this lab, we are going to show you briefly how this works.


## Task {{% param sectionnumber %}}.1: Enable the Host Firewall in Cilium

We need to enable the host firewall in the Cilium config. This can be done using Helm:


```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace kube-system \
  --reuse-values \
  --set hostFirewall.enabled=true \
  --set devices='{eth0}' \
  --wait
```

The devices flag refers to the network devices Cilium is configured on such as `eth0`. Omitting this option leads Cilium to auto-detect what interfaces the host firewall applies to.

Make sure to restart the `cilium` Pods with:

{{% alert title="Note" color="primary" %}}
You will see some deprecation warnings in this command. You can ignore them.
{{% /alert %}}

```bash
kubectl -n kube-system rollout restart ds/cilium
```

At this point, the Cilium-managed nodes are ready to enforce Network Policies.


## Task {{% param sectionnumber %}}.2: Attach a Label to the Node

In this lab, we will apply host policies only to nodes with the label `node-access=ssh`. We thus first need to attach that label to a node in the cluster.

```bash
kubectl label node cluster1 node-access=ssh
```


## Task {{% param sectionnumber %}}.2: Enable Policy Audit Mode for the Host Endpoint

[Host Policies](https://docs.cilium.io/en/latest/policy/language/#hostpolicies) enforce access control over connectivity to and from nodes. Particular care must be taken to ensure that when host policies are imported, Cilium does not block access to the nodes or break the cluster’s normal behavior (for example by blocking communication with kube-apiserver).

To avoid such issues, we can switch the host firewall in audit mode, to validate the impact of host policies before enforcing them. When Policy Audit Mode is enabled, no network policy is enforced so this setting is not recommended for production deployment.

```bash
CILIUM_POD_NAME=$(kubectl -n kube-system get pods -l "k8s-app=cilium" -o jsonpath="{.items[?(@.spec.nodeName=='cluster1')].metadata.name}")
HOST_EP_ID=$(kubectl -n kube-system exec $CILIUM_POD_NAME -- cilium endpoint list -o jsonpath='{[?(@.status.identity.id==1)].id}')
kubectl -n kube-system exec $CILIUM_POD_NAME -- cilium endpoint config $HOST_EP_ID PolicyAuditMode=Enabled
```

Verification:

```bash
kubectl -n kube-system exec $CILIUM_POD_NAME -- cilium endpoint config $HOST_EP_ID | grep PolicyAuditMode
```

The output should show you:

```
PolicyAuditMode          Enabled
```


## Task {{% param sectionnumber %}}.3: Apply a Host Network Policy

Host Policies match on node labels using a Node Selector to identify the nodes to which the policy applies. The following policy applies to all nodes. It allows communications from outside the cluster only on port TCP/22. All communications from the cluster to the hosts are allowed.

Host policies don’t apply to communications between pods or between pods and the outside of the cluster, except if those pods are host-networking pods.

Create a file `ccwnp.yaml` with the following content:

{{< highlight yaml >}}{{< readfile file="content/en/docs/10/01/ccwnp.yaml" >}}{{< /highlight >}}

And then apply this `CiliumClusterwideNetworkPolicy` with:

```bash
kubectl apply -f ccwnp.yaml
```

The host is represented as a special endpoint, with label `reserved:host`, in the output of the command `cilium endpoint list`. You can therefore inspect the status of the policy using that command:

```bash
kubectl -n kube-system exec $(kubectl -n kube-system get pods -l k8s-app=cilium -o jsonpath='{.items[0].metadata.name}') -- cilium endpoint list
```
You will see that the ingress policy enforcement for the `reserved:host` endpoint is `Disabled` but with `Audit` enabled:

```
Defaulted container "cilium-agent" out of: cilium-agent, mount-cgroup (init), clean-cilium-state (init)
ENDPOINT   POLICY (ingress)   POLICY (egress)   IDENTITY   LABELS (source:key[=value])                                                  IPv6   IPv4         STATUS   
           ENFORCEMENT        ENFORCEMENT                                                                                                                   
671        Disabled (Audit)   Disabled          1          k8s:minikube.k8s.io/commit=3e64b11ed75e56e4898ea85f96b2e4af0301f43d                              ready   
                                                           k8s:minikube.k8s.io/name=cluster1                                                                        
                                                           k8s:minikube.k8s.io/updated_at=2022_02_14T13_45_35_0700                                                  
                                                           k8s:minikube.k8s.io/version=v1.25.1                                                                      
                                                           k8s:node-access=ssh                                                                                      
                                                           k8s:node-role.kubernetes.io/control-plane                                                                
                                                           k8s:node-role.kubernetes.io/master                                                                       
                                                           k8s:node.kubernetes.io/exclude-from-external-load-balancers                                              
                                                           reserved:host                                                                                            
810        Disabled           Disabled          129160     k8s:io.cilium.k8s.namespace.labels.kubernetes.io/metadata.name=kube-system          10.1.0.249   ready   
                                                           k8s:io.cilium.k8s.policy.cluster=cluster1                                                                
                                                           k8s:io.cilium.k8s.policy.serviceaccount=coredns                                                          
                                                           k8s:io.kubernetes.pod.namespace=kube-system                                                              
                                                           k8s:k8s-app=kube-dns                                                                                     
4081       Disabled           Disabled          4          reserved:health                  
```


As long as the host endpoint is running in audit mode, communications disallowed by the policy won’t be dropped. They will however be reported by `cilium monitor` as `action audit`. The audit mode thus allows you to adjust the host policy to your environment, to avoid unexpected connection breakages.

You can montitor the policy verdicts with:

```bash
kubectl -n kube-system exec $(kubectl -n kube-system get pods -l k8s-app=cilium -o jsonpath='{.items[0].metadata.name}') -- cilium monitor -t policy-verdict --related-to $HOST_EP_ID
```

Open a second terminal to produce some traffic:

{{% alert title="Note" color="primary" %}}
If you are working in our Webshell environment, make sure to first login again to your VM after opening the second terminal.
{{% /alert %}}

```bash
curl -k https://192.168.49.2:8443
```

Also try to start an SSH session (you can cancel the command when the password promt is shown):

```bash
ssh 192.168.49.2
```

In the verdict log you should see an output similar to the following one. For the `curl` request you see that the action is set to `audit`:

```
Policy verdict log: flow 0xfd71ed86 local EP ID 671, remote ID world, proto 6, ingress, action audit, match none, 192.168.49.1:50760 -> 192.168.49.2:8443 tcp SYN
Policy verdict log: flow 0xfd71ed86 local EP ID 671, remote ID world, proto 6, ingress, action audit, match none, 192.168.49.1:50760 -> 192.168.49.2:8443 tcp SYN
```

The request to the SSH port has action `allow`:

```
Policy verdict log: flow 0x6b5b1b60 local EP ID 671, remote ID world, proto 6, ingress, action allow, match L4-Only, 192.168.49.1:48254 -> 192.168.49.2:22 tcp SYN
Policy verdict log: flow 0x6b5b1b60 local EP ID 671, remote ID world, proto 6, ingress, action allow, match L4-Only, 192.168.49.1:48254 -> 192.168.49.2:22 tcp SYN
```


## Task {{% param sectionnumber %}}.4: Clean Up

Once you are confident all required communication to the host from outside the cluster is allowed, you can disable policy audit mode to enforce the host policy.

{{% alert title="Note" color="primary" %}}
When enforcing the host policy, make sure that none of the communications required to access the cluster or for the cluster to work properly are denied. They should appear as `action allow`.
{{% /alert %}}

We are not going to do this extended task (as it would require some more rules for the cluster to continue working). But the command to disable the audit mode looks like this:

```
# kubectl -n kube-system  exec $CILIUM_POD_NAME -- cilium endpoint config $HOST_EP_ID PolicyAuditMode=Disabled
```

Simply cleanup and continue:

```bash
kubectl delete ccnp demo-host-policy
kubectl label node cluster1 node-access-
```
