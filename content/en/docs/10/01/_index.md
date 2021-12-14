---
title: "10.1 Kubernetes Without kube-proxy"
weight: 101
sectionnumber: 10.1
---

In this lab we are going to provision a new Kubernetes cluster without `kube-proxy` to use Cilium as a fully replacement for it.


## Task {{% param sectionnumber %}}.1: Deploy a new Kubernetes Cluster without `kube-proxy`


Create a new Kubernetes Cluster using the `minikube`. As `minikube` uses `kubeadm` we can skip the phase where `kubeadm` installs the `kube-proxy` addon. Execute the following command to create a third cluster:

```
minikube start --network-plugin=cni --cni=false --kubernetes-version=1.23.0 --extra-config=kubeadm.skip-phases=addon/kube-proxy -p cluster3
```

```
😄  [cluster3] minikube v1.24.0 on Ubuntu 20.04
✨  Automatically selected the docker driver. Other choices: virtualbox, ssh
❗  With --network-plugin=cni, you will need to provide your own CNI. See --cni flag as a user-friendly alternative
👍  Starting control plane node cluster3 in cluster cluster3
🚜  Pulling base image ...
🔥  Creating docker container (CPUs=2, Memory=8000MB) ...
🐳  Preparing Kubernetes v1.23.0 on Docker 20.10.8 ...
    ▪ kubeadm.skip-phases=addon/kube-proxy
    ▪ Generating certificates and keys ...
    ▪ Booting up control plane ...
    ▪ Configuring RBAC rules ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "cluster3" cluster and "default" namespace by default
```


## Task {{% param sectionnumber %}}.1: Deploy Cilium and enable the Kube Proxy replacement

Install Cilium with the following command:


```bash
cilium install --config cluster-pool-ipv4-cidr=10.3.0.0/16  --kube-proxy-replacement strict --cluster-name cluster3 --cluster-id 3 --wait false
```

{{% alert title="Note" color="primary" %}}

As the `cilium` and `cilium-operator` by default tries to communicate with the Kubernetes API using the default `kubernetes` service ip, they cannot do this with disabled `kube-proxy`. We therefore need to set the `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT` environment variable to tell the two pods how to connect to the Kubernetes API. Unfortunatly as of writing this, the `cilium` CLI does not do this when setting `--kube-proxy-replacement strict`.

First, get the correct IP address with:

```bash
kubectl config view -o jsonpath='{.clusters[?(@.name == "cluster3")].cluster.server}'
```

and then set these two environment variables with using the output from the previous command:

```bash
kubectl -n kube-system set env daemonset/cilium KUBERNETES_SERVICE_HOST=<API-IPADDRESS> KUBERNETES_SERVICE_PORT=<API-PORT>
kubectl -n kube-system set env deployment/cilium-operator KUBERNETES_SERVICE_HOST=<API-IPADDRESS> KUBERNETES_SERVICE_PORT=<API-PORT>
```
{{% /alert %}}

We can now compare the running pods on `cluster2` and `cluster3` in the `kube-system` namespace.

```bash
kubectl --context cluster2 -n kube-system get pod
```

Here we see the running `kube-proxy` pod:

```
NAME                                    READY   STATUS    RESTARTS       AGE
cilium-mvh65                            1/1     Running   1 (100m ago)   104m
cilium-operator-56f9689f68-twv8k        1/1     Running   2 (100m ago)   21h
clustermesh-apiserver-5cc7c7b64-vbnwk   2/2     Running   4 (100m ago)   21h
coredns-64897985d-f2s5c                 1/1     Running   2 (100m ago)   21h
etcd-cluster2                           1/1     Running   2 (100m ago)   21h
hubble-relay-6486ddd7cc-2c9p2           1/1     Running   2 (100m ago)   21h
kube-apiserver-cluster2                 1/1     Running   2 (99m ago)    21h
kube-controller-manager-cluster2        1/1     Running   2 (100m ago)   21h
kube-proxy-gd8w4                        1/1     Running   2 (100m ago)   21h
kube-scheduler-cluster2                 1/1     Running   2 (99m ago)    21h
storage-provisioner                     1/1     Running   4 (100m ago)   21h

```

while on `cluster3` there is no `kube-proxy` pod anymore:

```
NAME                               READY   STATUS    RESTARTS       AGE
cilium-operator-68bfb94678-785dk   1/1     Running   0              17m
cilium-vrqms                       1/1     Running   0              17m
coredns-64897985d-fk5lj            1/1     Running   0              59m
etcd-cluster3                      1/1     Running   0              59m
kube-apiserver-cluster3            1/1     Running   0              59m
kube-controller-manager-cluster3   1/1     Running   0              59m
kube-scheduler-cluster3            1/1     Running   0              59m
storage-provisioner                1/1     Running   13 (17m ago)   59m
```


## Task {{% param sectionnumber %}}.2: Deploy our sample-app again to the new cluster

As this is a new cluster we want to deploy our `simple-app.yaml` from lab 03 again to run some experiments. Run the following command using the `simple-app.yaml` from lab 03:

```bash
kubectl create -f simple-app.yaml
```

Now lets redo the task from lab 03.

Let's make life again a bit easier by storing the Pod's name into an environment variable so we can reuse it later again:

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
echo ${FRONTEND}
NOT_FRONTEND=$(kubectl get pods -l app=not-frontend -o jsonpath='{.items[0].metadata.name}')
echo ${NOT_FRONTEND}
```

Then execute:

```bash
kubectl exec -ti ${FRONTEND} -- curl -I --connect-timeout 5 backend:8080
kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
```
and then with the result you see that altought we have no `kube-proxy` running, the backend service can still be reached.

```
HTTP/1.1 200 OK
X-Powered-By: Express
Vary: Origin, Accept-Encoding
Access-Control-Allow-Credentials: true
Accept-Ranges: bytes
Cache-Control: public, max-age=0
Last-Modified: Sat, 26 Oct 1985 08:15:00 GMT
ETag: W/"83d-7438674ba0"
Content-Type: text/html; charset=UTF-8
Content-Length: 2109
Date: Tue, 14 Dec 2021 10:01:16 GMT
Connection: keep-alive

HTTP/1.1 200 OK
X-Powered-By: Express
Vary: Origin, Accept-Encoding
Access-Control-Allow-Credentials: true
Accept-Ranges: bytes
Cache-Control: public, max-age=0
Last-Modified: Sat, 26 Oct 1985 08:15:00 GMT
ETag: W/"83d-7438674ba0"
Content-Type: text/html; charset=UTF-8
Content-Length: 2109
Date: Tue, 14 Dec 2021 10:01:16 GMT
Connection: keep-alive
```
