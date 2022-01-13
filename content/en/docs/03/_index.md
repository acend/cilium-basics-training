---
title: "5. Hubble"
weight: 5
sectionnumber: 5
---

Before we start with the CNI functionality of Cilium and its security components. We want to enable the optional Hubble component (which is disabled by default). So we can take full advantage of its eBFP observatility capabilities.


## Task {{% param sectionnumber %}}.1: Install the Hubble CLI

Similiar to the `cilium` CLI, the `hubble` CLI interfaces with Hubble and allows observing network traffic within Kubernetes.

So let us install the `hubble` CLI.


### Linux Setup

Execute the following command to download the `hubble` CLI:

```bash
export HUBBLE_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/hubble/master/stable.txt)
curl -L --remote-name-all https://github.com/cilium/hubble/releases/download/$HUBBLE_VERSION/hubble-linux-amd64.tar.gz{,.sha256sum}
sha256sum --check hubble-linux-amd64.tar.gz.sha256sum
sudo tar xzvfC hubble-linux-amd64.tar.gz /usr/local/bin
rm hubble-linux-amd64.tar.gz{,.sha256sum}
```


### MacOS Setup

Execute the following command to download the `hubble` CLI:

```bash
export HUBBLE_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/hubble/master/stable.txt)
curl -L --remote-name-all https://github.com/cilium/hubble/releases/download/$HUBBLE_VERSION/hubble-darwin-amd64.tar.gz{,.sha256sum}
shasum -a 256 -c hubble-darwin-amd64.tar.gz.sha256sum
sudo tar xzvfC hubble-darwin-amd64.tar.gz /usr/local/bin
rm hubble-darwin-amd64.tar.gz{,.sha256sum}
```


### Windows Setup

Get the Windows binary files from the [latest Release](https://github.com/cilium/hubble/releases/latest/)


## Hubble CLI

Now that we have the `hubble` CLI let's have a look at some commands:

```bash
hubble version
```

```
hubble v0.8.2 compiled with go1.16.8 on linux/amd64
```

or

```bash
hubble help
```

```
Hubble is a utility to observe and inspect recent Cilium routed traffic in a cluster.

Usage:
  hubble [command]

Available Commands:
  completion  Output shell completion code
  config      Modify or view hubble config
  help        Help about any command
  list        List Hubble objects
  observe     Observe flows of a Hubble server
  status      Display status of Hubble server
  version     Display detailed version information

Global Flags:
      --config string   Optional config file (default "/home/sebastian/.config/hubble/config.yaml")
  -D, --debug           Enable debug messages

Get help:
  -h, --help    Help for any command or subcommand

Use "hubble [command] --help" for more information about a command.

```


## Task {{% param sectionnumber %}}.2: Deploy simple application

Before we enable Hubble in Cilium we want to make sure we have at least one application to observe.

Let's have a look at the following resource definitions:

{{< highlight yaml >}}{{< readfile file="content/en/docs/03/simple-app.yaml" >}}{{< /highlight >}}

The application consists of two client deployments (`frontend` and `not-frontend`) and one backend deployment (`backend`). We are going to send requests from the frontend and not-frontend pods to the backend pod.

Deploy the app:

```bash
kubectl create -f simple-app.yaml
```

Verify with the following command that everything is up and running:

```bash
kubectl get all
```

```
NAME                                READY   STATUS              RESTARTS   AGE
pod/backend-56787b4bd7-dmzdh        0/1     ContainerCreating   0          6s
pod/frontend-7cbdcb86fd-gdb4q       0/1     ContainerCreating   0          6s
pod/not-frontend-5cf6d96558-gj4np   0/1     ContainerCreating   0          6s

NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/backend      ClusterIP   10.106.244.233   <none>        8080/TCP   6s
service/kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP    36m

NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/backend        0/1     1            0           6s
deployment.apps/frontend       0/1     1            0           6s
deployment.apps/not-frontend   0/1     1            0           6s

NAME                                      DESIRED   CURRENT   READY   AGE
replicaset.apps/backend-56787b4bd7        1         1         0       6s
replicaset.apps/frontend-7cbdcb86fd       1         1         0       6s
replicaset.apps/not-frontend-5cf6d96558   1         1         0       6s

```

Let us make life a bit easier by storing the pods name into an environment variable so we can reuse it later again:

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
echo ${FRONTEND}
NOT_FRONTEND=$(kubectl get pods -l app=not-frontend -o jsonpath='{.items[0].metadata.name}')
echo ${NOT_FRONTEND}
```


## Task {{% param sectionnumber %}}.3: Enable Hubble in Cilium

When you install Cilium using Helm, then Hubble is already enabled. The value for this is `hubble.enabled` which is set to `true` in the `values.yaml` of the Cilium Helm Chart. But we also want to enable the Hubble Relay. With the following Helm command you can enable Hubble with the Hubble Relay:

```bash
helm repo add cilium https://helm.cilium.io/
helm upgrade -i cilium cilium/cilium --version 1.11.0 \
  --namespace kube-system \
  --reuse-values \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --wait
```

If you have installed Cilium with the `cilium` cli then Hubble component is not enabled by default (nor is Hubble Relay), and you can enbale Hubble using the following `cilium` CLI command:

```bash
cilium hubble enable
```

and then wait until Hubble is enabled:

```
🔑 Found existing CA in secret cilium-ca
✨ Patching ConfigMap cilium-config to enable Hubble...
♻️  Restarted Cilium pods
⌛ Waiting for Cilium to become ready before deploying other Hubble component(s)...
🔑 Generating certificates for Relay...
✨ Deploying Relay from quay.io/cilium/hubble-relay:v1.10.5...
⌛ Waiting for Hubble to be installed...
✅ Hubble was successfully enabled!
```

When you have a look at your running pods with `kubectl get pod -A` you should see a pod with a name starting with `hubble-relay`:

```
kubectl get pod -A                                                                         
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE
default       backend-56787b4bd7-dmzdh           1/1     Running   0          114m
default       frontend-7cbdcb86fd-gdb4q          1/1     Running   0          114m
default       not-frontend-5cf6d96558-gj4np      1/1     Running   0          114m
kube-system   cilium-28kmn                       1/1     Running   0          73s
kube-system   cilium-operator-8dd4dc946-n9ght    1/1     Running   0          149m
kube-system   coredns-558bd4d5db-xzvc9           1/1     Running   0          150m
kube-system   etcd-minikube                      1/1     Running   0          150m
kube-system   hubble-relay-f6d85866c-csthd       1/1     Running   0          41s
kube-system   kube-apiserver-minikube            1/1     Running   0          150m
kube-system   kube-controller-manager-minikube   1/1     Running   0          150m
kube-system   kube-proxy-bqs4d                   1/1     Running   0          150m
kube-system   kube-scheduler-minikube            1/1     Running   0          150m
kube-system   storage-provisioner                1/1     Running   1          150m
```

Cilium agents are restarting, and a new Hubble Relay pod is now present. We can wait for Cilium and Hubble to be ready by running:

```bash
cilium status --wait
```

which should give you an output similar to this:

```
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         OK
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

DaemonSet         cilium             Desired: 1, Ready: 1/1, Available: 1/1
Deployment        cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
Deployment        hubble-relay       Desired: 1, Ready: 1/1, Available: 1/1
Containers:       cilium             Running: 1
                  cilium-operator    Running: 1
                  hubble-relay       Running: 1
Cluster Pods:     9/9 managed by Cilium
Image versions    cilium             quay.io/cilium/cilium:v1.11.0@sha256:ea677508010800214b0b5497055f38ed3bff57963fa2399bcb1c69cf9476453a: 1
                  cilium-operator    quay.io/cilium/operator-generic:v1.11.0@sha256:b522279577d0d5f1ad7cadaacb7321d1b172d8ae8c8bc816e503c897b420cfe3: 1
                  hubble-relay       quay.io/cilium/hubble-relay:v1.11.0@sha256:306ce38354a0a892b0c175ae7013cf178a46b79f51c52adb5465d87f14df0838: 1
```

So the Hubble component is now enabled.

Once ready, we can locally port-forward to the Hubble pod:

```bash
cilium hubble port-forward&
```

{{% alert title="Note" color="primary" %}}
Note the `&` after the command which puts the process in the background so we can continue working in the shell.
{{% /alert %}}

And then check Hubble status via the Hubble CLI (which uses the port-forwarding just openend):

```bash
hubble status
```

```
Healthcheck (via localhost:4245): Ok
Current/Max Flows: 947/4095 (23.13%)
Flows/s: 3.84
Connected Nodes: 1/1
```

The Hubble CLI is now primed for observing network traffic within the cluster.


## Task {{% param sectionnumber %}}.4: Observing flows with Hubble

We now want to use the `hubble` cli to observe some network flows in out Kubernetes Cluster. Lets have a look at the following command:

```bash
hubble observe
```

which gives you a list on network flows:

```
Nov 23 14:49:03.030: 10.0.0.113:46274 <- kube-system/hubble-relay-f6d85866c-csthd:4245 to-stack FORWARDED (TCP Flags: ACK, PSH)
Nov 23 14:49:03.030: 10.0.0.113:46274 -> kube-system/hubble-relay-f6d85866c-csthd:4245 to-endpoint FORWARDED (TCP Flags: RST)
Nov 23 14:49:04.011: 10.0.0.113:44840 <- 10.0.0.114:4240 to-stack FORWARDED (TCP Flags: ACK)
Nov 23 14:49:04.011: 10.0.0.113:44840 -> 10.0.0.114:4240 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:49:04.226: 10.0.0.113:32898 -> kube-system/coredns-558bd4d5db-xzvc9:8080 to-endpoint FORWARDED (TCP Flags: SYN)
Nov 23 14:49:04.226: 10.0.0.113:32898 <- kube-system/coredns-558bd4d5db-xzvc9:8080 to-stack FORWARDED (TCP Flags: SYN, ACK)
Nov 23 14:49:04.227: 10.0.0.113:32898 -> kube-system/coredns-558bd4d5db-xzvc9:8080 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:49:04.227: 10.0.0.113:32898 -> kube-system/coredns-558bd4d5db-xzvc9:8080 to-endpoint FORWARDED (TCP Flags: ACK, PSH)
Nov 23 14:49:04.227: 10.0.0.113:32898 <- kube-system/coredns-558bd4d5db-xzvc9:8080 to-stack FORWARDED (TCP Flags: ACK, PSH)
Nov 23 14:49:04.227: 10.0.0.113:32898 -> kube-system/coredns-558bd4d5db-xzvc9:8080 to-endpoint FORWARDED (TCP Flags: ACK, FIN)
Nov 23 14:49:04.227: 10.0.0.113:32898 <- kube-system/coredns-558bd4d5db-xzvc9:8080 to-stack FORWARDED (TCP Flags: ACK, FIN)
Nov 23 14:49:04.227: 10.0.0.113:32898 -> kube-system/coredns-558bd4d5db-xzvc9:8080 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:49:04.842: 10.0.0.113:34716 -> kube-system/coredns-558bd4d5db-xzvc9:8181 to-endpoint FORWARDED (TCP Flags: SYN)
Nov 23 14:49:04.842: 10.0.0.113:34716 <- kube-system/coredns-558bd4d5db-xzvc9:8181 to-stack FORWARDED (TCP Flags: SYN, ACK)
Nov 23 14:49:04.842: 10.0.0.113:34716 -> kube-system/coredns-558bd4d5db-xzvc9:8181 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:49:04.842: 10.0.0.113:34716 -> kube-system/coredns-558bd4d5db-xzvc9:8181 to-endpoint FORWARDED (TCP Flags: ACK, PSH)
Nov 23 14:49:04.842: 10.0.0.113:34716 <- kube-system/coredns-558bd4d5db-xzvc9:8181 to-stack FORWARDED (TCP Flags: ACK, PSH)
Nov 23 14:49:04.843: 10.0.0.113:34716 <- kube-system/coredns-558bd4d5db-xzvc9:8181 to-stack FORWARDED (TCP Flags: ACK, FIN)
Nov 23 14:49:04.843: 10.0.0.113:34716 -> kube-system/coredns-558bd4d5db-xzvc9:8181 to-endpoint FORWARDED (TCP Flags: ACK, FIN)
Nov 23 14:49:05.971: kube-system/hubble-relay-f6d85866c-csthd:40844 -> 192.168.49.2:4244 to-stack FORWARDED (TCP Flags: ACK, PSH)

```

with

```bash
hubble observe -f
```

you can observe and follow the current active flows in your Kubernetes Cluster. Stop the command with `CTRL+C`.

Let us produce some traffic:

```bash
for i in {1..10}; do
  kubectl exec -ti ${FRONTEND} -- curl -I --connect-timeout 5 backend:8080
  kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
done
```

We can now use the `hubble` cli to filter traffic we are interested in. Here are some examples to specifically retrieve the network activity between our frontends and backend:

```bash
hubble observe --to-pod backend
hubble observe --namespace default --protocol tcp --port 8080
hubble observe --verdict DROPPED
```

```
hubble observe --to-pod backend
Nov 23 14:54:27.091: default/frontend-7cbdcb86fd-gdb4q:58842 -> default/backend-56787b4bd7-dmzdh:8080 L3-Only FORWARDED (TCP Flags: SYN)
Nov 23 14:54:27.091: default/frontend-7cbdcb86fd-gdb4q:58842 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: SYN)
Nov 23 14:54:27.091: default/frontend-7cbdcb86fd-gdb4q:58842 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:54:27.091: default/frontend-7cbdcb86fd-gdb4q:58842 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK, PSH)
Nov 23 14:54:27.092: default/frontend-7cbdcb86fd-gdb4q:58842 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK, FIN)
Nov 23 14:54:27.092: default/frontend-7cbdcb86fd-gdb4q:58842 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:54:27.267: default/not-frontend-5cf6d96558-gj4np:53766 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
Nov 23 14:54:27.267: default/not-frontend-5cf6d96558-gj4np:53766 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
Nov 23 14:54:28.295: default/not-frontend-5cf6d96558-gj4np:53766 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
Nov 23 14:54:28.295: default/not-frontend-5cf6d96558-gj4np:53766 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
Nov 23 14:54:30.311: default/not-frontend-5cf6d96558-gj4np:53766 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
Nov 23 14:54:30.311: default/not-frontend-5cf6d96558-gj4np:53766 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
Nov 23 14:54:32.430: default/frontend-7cbdcb86fd-gdb4q:58894 -> default/backend-56787b4bd7-dmzdh:8080 L3-Only FORWARDED (TCP Flags: SYN)
Nov 23 14:54:32.430: default/frontend-7cbdcb86fd-gdb4q:58894 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: SYN)
Nov 23 14:54:32.430: default/frontend-7cbdcb86fd-gdb4q:58894 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:54:32.430: default/frontend-7cbdcb86fd-gdb4q:58894 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK, PSH)
Nov 23 14:54:32.431: default/frontend-7cbdcb86fd-gdb4q:58894 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK, FIN)
Nov 23 14:54:32.431: default/frontend-7cbdcb86fd-gdb4q:58894 -> default/backend-56787b4bd7-dmzdh:8080 to-endpoint FORWARDED (TCP Flags: ACK)
Nov 23 14:54:32.603: default/not-frontend-5cf6d96558-gj4np:53820 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
Nov 23 14:54:32.603: default/not-frontend-5cf6d96558-gj4np:53820 <> default/backend-56787b4bd7-dmzdh:8080 Policy denied DROPPED (TCP Flags: SYN)
```

Note that Hubble tells us the reason a packet was `DROPPED` (in our case, denied by the network policies applied above). This is really handy when developing / debugging network policies in the later chapters.
