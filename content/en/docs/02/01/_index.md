---
title: "2.1 Install Cilium"
weight: 21
sectionnumber: 2.1
---

Cilium can be installed using multiple ways:

* Cilium CLI
* Using Helm

In this lab, we are going to use [Helm](https://helm.sh) which is recommended for production use.
The [Cilium command-line](https://github.com/cilium/cilium-cli/) tool is used (Cilium CLI) for verification and troubleshooting.


## Task {{% param sectionnumber %}}.1: Install a Kubernetes Cluster

We are going to spin up a new Kubernetes cluster with the following command:

{{% alert title="Note" color="primary" %}}
To start from a clean Kubernetes cluster, make sure `cluster1` is not yet available. You can verify this with `minikube profile list`. If you already have a `cluster1` you can delete the cluster with `minikube delete -p cluster1`.
{{% /alert %}}

```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version={{% param "kubernetesVersion" %}} -p cluster1 
```

{{% alert title="Note" color="primary" %}}
During this training, you will create multiple clusters. For this, we use a feature in Minikube called profile which you see with the `-p cluster1` option. You can list all your profiles with `minikube profile list` and you can change to another cluster with `minikube profile <profilename>`, this will also set your current context for `kubectl` to the specified profile/cluster.
{{% /alert %}}

Minikube installed a new Kubernetes cluster without any Container Network Interface (CNI). CNI installation happens in the next task.

Minikube added a new context to your Kubernetes config file and set this as default. Verify this with the following command:

```bash
kubectl config current-context
```

This should show `cluster1`. Now check that everything is up and running:

```bash
kubectl get node
```

This should produce a similar output:

```
NAME       STATUS   ROLES                  AGE   VERSION
cluster1   Ready    control-plane,master   86s   v{{% param "kubernetesVersion" %}}
```
Depending on your Minikube version and environment your node might stay NotReady because no CNI is installed. After we installed Cilium it will become ready.

Check if all pods are running with:

```bash
kubectl get pod -A
```

which produces the following output

```
NAMESPACE     NAME                               READY   STATUS              RESTARTS   AGE
kube-system   coredns-558bd4d5db-wdlrh           0/1     ContainerCreating   0          3m1s
kube-system   etcd-minikube                      1/1     Running             0          3m7s
kube-system   kube-apiserver-cluster1            1/1     Running             0          3m16s
kube-system   kube-controller-manager-cluster1   1/1     Running             0          3m7s
kube-system   kube-proxy-9bjbq                   1/1     Running             0          3m1s
kube-system   kube-scheduler-cluster1            1/1     Running             0          3m7s
kube-system   storage-provisioner                1/1     Running             0          3m11s
```


{{% alert title="Note" color="primary" %}}
Depending on your Minikube version, coredns might start or not which is ok.
But you should not see any CNI related pods!
{{% /alert %}}


## Task {{% param sectionnumber %}}.2: Install Cilium CLI

The `cilium` CLI tool is a single binary file that can be downloaded from the project's release page. Follow the instructions depending on your operating system or environment.


### Linux/Webshell Setup

{{% alert title="Note" color="primary" %}}
If you are working in our webshell based lab setup please always follow the Linux setup.
{{% /alert %}}


Execute the following command to download the `cilium` CLI:

```bash
curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/download/v{{% param "ciliumVersion.cli" %}}/cilium-linux-amd64.tar.gz{,.sha256sum}
sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
rm cilium-linux-amd64.tar.gz{,.sha256sum}
```


### macOS Setup

Execute the following command to download the `cilium` CLI:

```bash
curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/download/v{{% param "ciliumVersion.cli" %}}/cilium-darwin-amd64.tar.gz{,.sha256sum}
shasum -a 256 -c cilium-darwin-amd64.tar.gz.sha256sum
sudo tar xzvfC cilium-darwin-amd64.tar.gz /usr/local/bin
rm cilium-darwin-amd64.tar.gz{,.sha256sum}
```


## Cilium CLI

Now that we have the `cilium` CLI let's have a look at some commands:

```bash
cilium version
```

which should give you an output similar to this:

```
cilium-cli: v0.12.0 compiled with go1.18.4 on linux/amd64
cilium image (default): v1.12.0
cilium image (stable): v1.12.0
cilium image (running): unknown. Unable to obtain cilium version, no cilium pods found in namespace "kube-system"
```

{{% alert title="Note" color="primary" %}}
It's ok if your installation does not show the same version.
{{% /alert %}}

Then let us look at

```bash
cilium status
```

```
cilium status 
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         1 errors
 \__/¯¯\__/    Operator:       disabled
 /¯¯\__/¯¯\    Hubble:         disabled
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

Containers:      cilium             
                 cilium-operator    
Cluster Pods:    0/0 managed by Cilium
Errors:          cilium    cilium    daemonsets.apps "cilium" not found

```

We don't have yet installed Cilium, therefore the error is perfectly fine.


## Task {{% param sectionnumber %}}.3: Install Cilium

Let's install Cilium with Helm. First we need to add the Cilium Helm repository:

```bash
helm repo add cilium https://helm.cilium.io/ --force-update
```

and then we can install Cilium:

{{% onlyWhenNot techlab %}}

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.preUpgrade" %}} \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDR=10.1.0.0/16 \
  --set cluster.name=cluster1 \
  --set cluster.id=1 \
  --set operator.replicas=1 \
  --set kubeProxyReplacement=disabled \
  --wait
```
{{% /onlyWhenNot %}}

{{% onlyWhen techlab %}}

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDRList={10.1.0.0/16} \
  --set cluster.name=cluster1 \
  --set cluster.id=1 \
  --set operator.replicas=1 \
  --wait
```
{{% /onlyWhen %}}

For all values possible in the Cilium Helm chart, have a look at the [Repository](https://github.com/cilium/cilium/tree/master/install/kubernetes/cilium) or the [Helm Reference](https://docs.cilium.io/en/stable/helm-reference/) in Cilium's documentation. {{% onlyWhenNot techlab %}} We disable the kubeProxyReplacement because it would cause problems with multiple clusters running on the same kernel in the later chapters.{{% /onlyWhenNot %}}

and now run again the

```bash
cilium status --wait
```

command:

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
Containers:       cilium             Running: 1
                  cilium-operator    Running: 1
Cluster Pods:     1/1 managed by Cilium
Image versions    cilium             quay.io/cilium/cilium:v{{% param "ciliumVersion.preUpgrade" %}}: 1
                  cilium-operator    quay.io/cilium/operator-generic:v{{% param "ciliumVersion.preUpgrade" %}}: 1

```

{{% alert title="Note" color="primary" %}}
If the output is not the same, make sure all Cilium container are up and in a ready state.
{{% /alert %}}


Take a look at the pods again to see what happened under the hood:

```bash
kubectl get pods -A
```

and you should see now the Pods related to Cilium:

```
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE
kube-system   cilium-hsk8g                       1/1     Running   0          89s
kube-system   cilium-operator-8dd4dc946-n9ght    1/1     Running   0          89s
kube-system   coredns-558bd4d5db-xzvc9           1/1     Running   0          111s
kube-system   etcd-minikube                      1/1     Running   0          118s
kube-system   kube-apiserver-minikube            1/1     Running   0          118s
kube-system   kube-controller-manager-minikube   1/1     Running   0          118s
kube-system   kube-proxy-bqs4d                   1/1     Running   0          111s
kube-system   kube-scheduler-minikube            1/1     Running   0          118s
kube-system   storage-provisioner                1/1     Running   1          2m3s

```

{{% alert title="Note" color="primary" %}}
It might take some time until all Pods ar in state `Runnning` and `READY`. Wait before continue.
{{% /alert %}}

Alright, Cilium is up and running, let us make some tests. The `cilium` CLI allows you to run a connectivity test:

```bash
cilium connectivity test
```

This will run for some minutes, let's wait.

{{% alert title="Note" color="primary" %}}
As we installed an older version of cilium but are using the latest `cilium` CLI, it's ok if some tests are failing.
{{% /alert %}}

```
ℹ️  Single-node environment detected, enabling single-node connectivity test
ℹ️  Monitor aggregation detected, will skip some flow validation steps
✨ [minikube] Creating namespace for connectivity check...
✨ [minikube] Deploying echo-same-node service...
✨ [minikube] Deploying same-node deployment...
✨ [minikube] Deploying client deployment...
✨ [minikube] Deploying client2 deployment...
⌛ [minikube] Waiting for deployments [client client2 echo-same-node] to become ready...
⌛ [minikube] Waiting for deployments [] to become ready...
⌛ [minikube] Waiting for CiliumEndpoint for pod cilium-test/client-6488dcf5d4-fkv57 to appear...
⌛ [minikube] Waiting for CiliumEndpoint for pod cilium-test/client2-5998d566b4-l66kc to appear...
⌛ [minikube] Waiting for CiliumEndpoint for pod cilium-test/echo-same-node-745bd5c77-dqxr9 to appear...
⌛ [minikube] Waiting for Service cilium-test/echo-same-node to become ready...
⌛ [minikube] Waiting for NodePort 192.168.49.2:31041 (cilium-test/echo-same-node) to become ready...
ℹ️  Skipping IPCache check
⌛ [minikube] Waiting for pod cilium-test/client-6488dcf5d4-fkv57 to reach default/kubernetes service...
⌛ [minikube] Waiting for pod cilium-test/client2-5998d566b4-l66kc to reach default/kubernetes service...
🔭 Enabling Hubble telescope...
⚠️  Unable to contact Hubble Relay, disabling Hubble telescope and flow validation: rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 127.0.0.1:4245: connect: connection refused"
ℹ️  Expose Relay locally with:
   cilium hubble enable
   cilium status --wait
   cilium hubble port-forward&
⚠️  Unable to detect Cilium version, assuming v1.12.0 for connectivity tests
ℹ️  Cilium version: 1.12.0

🏃 Running tests...

[=] Test [no-policies]
....................
[=] Test [allow-all]
................
[=] Test [client-ingress]
..
[=] Test [echo-ingress]
..
[=] Test [client-egress]
..
[=] Test [to-entities-world]
......
[=] Test [to-cidr-1111]
....
[=] Test [echo-ingress-l7]
..
[=] Test [client-egress-l7]
........
[=] Test [dns-only]
........
[=] Test [to-fqdns]
......
✅ All 11 tests (76 actions) successful, 0 tests skipped, 0 scenarios skipped.
```

Once done, clean up the connectivity test Namespace:

```bash
kubectl delete ns cilium-test --wait=false
```


## Task {{% param sectionnumber %}}.4: Explore your installation

We have learned about the Cilium components. Let us check out the installed CRDs now:

```bash
kubectl api-resources | grep cilium
``````

Which shows CRDs installed by Cilium:

```bash
ciliumclusterwidenetworkpolicies   ccnp           cilium.io/v2                           false        CiliumClusterwideNetworkPolicy
ciliumendpoints                    cep,ciliumep   cilium.io/v2                           true         CiliumEndpoint
ciliumexternalworkloads            cew            cilium.io/v2                           false        CiliumExternalWorkload
ciliumidentities                   ciliumid       cilium.io/v2                           false        CiliumIdentity
ciliumnetworkpolicies              cnp,ciliumnp   cilium.io/v2                           true         CiliumNetworkPolicy
ciliumnodes                        cn,ciliumn     cilium.io/v2                           false        CiliumNode
``````

And now we check all installed Cilium CRDs
```bash
kubectl get ccnp,cep,cew,ciliumid,cnp,cn -A
```

We see 1 node, 1 identity and 1 endpoint:
```bash
NAMESPACE     NAME                                               ENDPOINT ID   IDENTITY ID   INGRESS ENFORCEMENT   EGRESS ENFORCEMENT   VISIBILITY POLICY   ENDPOINT STATE   IPV4         IPV6
kube-system   ciliumendpoint.cilium.io/coredns-64897985d-7485t   465           67688                                                                        ready            10.1.0.215   

NAMESPACE   NAME                             NAMESPACE     AGE
            ciliumidentity.cilium.io/67688   kube-system   18m

NAMESPACE   NAME                            AGE
            ciliumnode.cilium.io/cluster1   18m
```

{{% alert title="Note" color="primary" %}}
It might be possible that you still see identites created by `cilium connectivity test`. They will be deleted by `cilium-operator` after max. 15 minutes.
{{% /alert %}}

{{% details title="Can you guess why only the coredns Pod is listed as an Endpoint and Identity?" %}}
This Pod is the only one which is NOT on the Host Network.
{{% /details %}}

{{% details title="Is it possible to have more CiliumNodes than nodes in a Kubernetes cluster?" %}}
A CiliumNode is a host with cilium-agent installed. So this could also be VM outside Kubernetes.
{{% /details %}}

We have discussed CNI plugin installations, let us check out the Cilium installation on the node.

We can either start a debug container on the node and chroot its /

```bash
kubectl debug node/cluster1 -it --image=busybox
chroot /host
```

or we use docker to access the node:
```bash
docker exec -it cluster1 bash
```


Now we have a shell with access to the node. We will check out the Cilium installation:

```bash
ls -l /etc/cni/net.d/
cat /etc/cni/net.d/05-cilium.conf
/opt/cni/bin/cilium-cni --help
ip a
exit #exit twice if you used kubectl debug
```
We make a few oberservations:

* Kubernetes uses the configuration file with the lowest number so it takes Cilium with the prefix 05.
* The configuration file is written as a  [CNI spec](https://github.com/containernetworking/cni/blob/master/SPEC.md#configuration-format).
* The `cilium` binary was installed to /opt/cni/bin.
* Cilium created a virtual network interfaces pair `cilium_net`,`cilium_host` and the vxlan overlay interface `cilium_vxlan`.
* We see the virtual network interface (`lxc` device) of the coredns pod (the Endpoint in Cilium terms).


## Install Cilium with the `cilium` cli

This is what the installation with the `cilium` cli would have looked like:

```
# cilium install --config cluster-pool-ipv4-cidr=10.1.0.0/16 --cluster-name cluster1 --cluster-id 1 --version {{% param "ciliumVersion.preUpgrade" %}}
```
Be careful to never use CLI and Helm together to install, this can break an already running Cilium installation.

After this initial installation, we can proceed by upgrading to a newer version in the next lab.
