---
title: "2. Install Cilium"
weight: 2
sectionnumber: 2
---


Cilium can be installed using multiple ways:

* Cilium CLI
* Using Helm

In this lab we are going to use [helm](https://helm.sh) since it has more options.
The [Cilium command line](https://github.com/cilium/cilium-cli/) tool is used (Cilium CLI) for verification and troubleshooting.


## Task {{% param sectionnumber %}}.1: Install helm

For a complete overview refer to the helm installation [website](https://helm.sh/docs/intro/install/). If you have helm 3 already installed you can skip this step.


### Linux and MacOs Setup

Use your package manager (`apt`, `yum`, `brew` etc), download the [latest Release](https://github.com/helm/helm/releases) or use the following command to install [helm](https://helm.sh/docs/intro/install/) helm:

```bash
curl -s https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```


### Windows Setup

Get the Windows binary files from the [latest Release](https://github.com/helm/helm/releases)


## Task {{% param sectionnumber %}}.2: Install Cilium CLI

The `cilium` CLI tool is a single binary file that can be downloaded from the project's release page. Follow the instructions depending on your operating system


### Linux Setup

Execute the following command to download the `cilium` CLI:

```bash
curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz{,.sha256sum}
sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
rm cilium-linux-amd64.tar.gz{,.sha256sum}
```


### MacOS Setup

Execute the following command to download the `cilium` CLI:

```bash
curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-darwin-amd64.tar.gz{,.sha256sum}
shasum -a 256 -c cilium-darwin-amd64.tar.gz.sha256sum
sudo tar xzvfC cilium-darwin-amd64.tar.gz /usr/local/bin
rm cilium-darwin-amd64.tar.gz{,.sha256sum}
```


### Windows Setup

Get the Windows binary files from the [latest Release](https://github.com/cilium/cilium-cli/releases/latest/)


## Cilium CLI

Now that we have the `cilium` CLI let's have a look at some commands:

```bash
cilium version
```

which should give you an output similar to this:

```
cilium-cli: v0.10.1 compiled with go1.17.6 on linux/amd64
cilium image (default): v1.11.1
cilium image (stable): v1.11.1
cilium image (running): unknown. Unable to obtain cilium version, no cilium pods found in namespace "kube-system"
```

{{% alert title="Note" color="primary" %}}
It's ok if your installation does not show the same version.
{{% /alert %}}

Then lets look at

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

We don't have yet installed cilium, therefore the error is perfectly fine.


## Task {{% param sectionnumber %}}.3: Install Cilium

Let's install cilium with helm:

```bash
helm repo add cilium https://helm.cilium.io/
helm upgrade -i cilium cilium/cilium --version 1.10.5 \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDR=10.1.0.0/16 \
  --set cluster.name=cluster1 \
  --set cluster.id=1 \
  --set operator.replicas=1 \
  --wait
```

and now run again the `cilium status` command:

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
Image versions    cilium             quay.io/cilium/cilium:v1.10.5: 1
                  cilium-operator    quay.io/cilium/operator-generic:v1.10.5: 1

```

Take a look at the pods again to see what happened under the hood:

```bash
kubectl get pods -A
```

and you should see now the Pods related to Cilium:

```
kubectl get pod -A
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

Alright, Cilium is up and running, lets make some tests. The `cilium` CLI allows you to run a connectivity test:

```bash
cilium connectivity test
```

This will run fore some minutes, let's wait.

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

Once done, clean up the connectivity test namespace:

```bash
kubectl delete ns cilium-test --wait=false
```


## Task {{% param sectionnumber %}}.4: Explore your installation

We have learned about the cilium components. Let us check out the installed CRDs now:

```bash
kubectl api-resources | grep cilium
``````

Which should output the installed CRDs:

```bash
ciliumclusterwidenetworkpolicies   ccnp           cilium.io/v2                           false        CiliumClusterwideNetworkPolicy
ciliumegressnatpolicies                           cilium.io/v2alpha1                     false        CiliumEgressNATPolicy
ciliumendpoints                    cep,ciliumep   cilium.io/v2                           true         CiliumEndpoint
ciliumexternalworkloads            cew            cilium.io/v2                           false        CiliumExternalWorkload
ciliumidentities                   ciliumid       cilium.io/v2                           false        CiliumIdentity
ciliumlocalredirectpolicies        clrp           cilium.io/v2                           true         CiliumLocalRedirectPolicy
ciliumnetworkpolicies              cnp,ciliumnp   cilium.io/v2                           true         CiliumNetworkPolicy
ciliumnodes                        cn,ciliumn     cilium.io/v2                           false        CiliumNode
``````

And now we check all installed cilium CRDs
```bash
kubectl get ccnp,cep,cew,ciliumid,clrp,cnp,cn -A
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

Can you guess why only the coredns pod is listed as an Endpoint and Identity?
<details>
  <summary>Answer</summary>
This pod is the only one which is NOT on the Host Network.
</details>

Is it possible to have more CiliumNodes than nodes in a Kubernetes Cluster?
<details>
  <summary>Answer</summary>
A CiliumNode is a host with cilium-agent installed. So this could also be VM outside Kubernetes.
</details>


We have discussed CNI plugin installations, let us check out the cilium installation on the Node.

We can either start debug debug container on the Node and chroot its /.

```bash
kubectl debug node/cluster1 -it --image=busybox
chroot /host
```

Or we use docker to access the node:
```bash
docker exec -it cluster1 bash
```


Now we have a shell with access to the node. We will check out the cilium installation:

```bash
ls -l /etc/cni/net.d/
cat /etc/cni/net.d/05-cilium.conf
/opt/cni/bin/cilium-cni --help
ip a
ls /sys/fs/bpf/tc/globals/
exit #exit twice if you used kubectl debug
```
We make a few oberservations:

* Kubernetes uses the configuration file with the lowest number so it takes cilium with the prefix 05.
* The configuration file is written as a  [CNI spec](https://github.com/containernetworking/cni/blob/master/SPEC.md#configuration-format).
* The cilium binary was installed to /opt/cni/bin.
* Cilium created two virtual network interfaces `cilium_net`,`cilium_host` (host traffic) and the vxlan overlay interface `cilium_vxlan`
* We see that cilium created eBPF Maps on the Node in /sys/fs/bpf/tc/globals/.


## Install Cilium with the `cilium` cli

This is how the installation with the `cilium` cli would have looked like:

```bash
cilium install --config cluster-pool-ipv4-cidr=10.1.0.0/16 --cluster-name cluster1 --cluster-id 1 --version v1.10.5
```
