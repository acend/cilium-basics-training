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


## Install helm

For a complete overview refer to the helm installation [website](https://helm.sh/docs/intro/install/). If you have helm 3 already installed you can skip this step.


### Linux and MacOs Setup

Use your package manager (`apt`, `yum`, `brew` etc), download the [latest Release](https://github.com/helm/helm/releases) or use the following command to install [helm](https://helm.sh/docs/intro/install/) helm:

```bash
curl -s https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```


### Windows Setup

Get the Windows binary files from the [latest Release](https://github.com/helm/helm/releases)


## Install Cilium CLI

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
cilium-cli: v0.9.3 compiled with go1.17.3 on linux/amd64
cilium image (default): v1.10.5
cilium image (stable): v1.10.5
cilium image (running): unknown. Unable to obtain cilium version, no cilium pods found in namespace "kube-system"
```

Them lets look at

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


## Install Cilium

Let's install cilium with helm:

```bash
helm repo add cilium https://helm.cilium.io/
helm upgrade -i cilium cilium/cilium --version 1.11.0 \
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


## Install Cilium with the `cilium` cli

This is how the installation with the `cilium` cli would have looked like:

```bash
cilium install --config cluster-pool-ipv4-cidr=10.1.0.0/16 --cluster-name cluster1 --cluster-id 1 
```
