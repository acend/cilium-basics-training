---
title: "10.2 Kubernetes Without kube-proxy"
weight: 102
sectionnumber: 10.2
OnlyWhenNot: techlab
---

In this lab, we are going to provision a new Kubernetes cluster without `kube-proxy` to use Cilium as a full replacement for it.


## Task {{% param sectionnumber %}}.1: Deploy a new Kubernetes Cluster without `kube-proxy`


Create a new Kubernetes cluster using `minikube`. As `minikube` uses `kubeadm` we can skip the phase where `kubeadm` installs the `kube-proxy` addon. Execute the following command to create a third cluster:

```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version={{% param "kubernetesVersion" %}} --extra-config=kubeadm.skip-phases=addon/kube-proxy -p kubeless
```

```
😄  [cluster3] minikube v{{% param "kubernetesVersion" %}} on Ubuntu 20.04
✨  Automatically selected the docker driver. Other choices: virtualbox, ssh
❗  With --network-plugin=cni, you will need to provide your own CNI. See --cni flag as a user-friendly alternative
👍  Starting control plane node cluster3 in cluster cluster3
🚜  Pulling base image ...
🔥  Creating docker container (CPUs=2, Memory=8000MB) ...
🐳  Preparing Kubernetes v{{% param "kubernetesVersion" %}} on Docker 20.10.8 ...
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

As the `cilium` and `cilium-operator` Pods by default try to communicate with the Kubernetes API using the default `kubernetes` service IP, they cannot do this with disabled `kube-proxy`. We, therefore, need to set the `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT` environment variables to tell the two Pods how to connect to the Kubernetes API.

To find the correct IP address execute the following command:

```bash
API_SERVER_IP=$(kubectl config view -o jsonpath='{.clusters[?(@.name == "kubeless")].cluster.server}' | cut -f 3 -d / | cut -f1 -d:)
API_SERVER_PORT=$(kubectl config view -o jsonpath='{.clusters[?(@.name == "kubeless")].cluster.server}' | cut -f 3 -d / | cut -f2 -d:)
echo "$API_SERVER_IP:$API_SERVER_PORT"
```

Use the shown IP address and port in the next Helm command to install Cilium:

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDRList={10.3.0.0/16} \
  --set cluster.name=kubeless \
  --set cluster.id=3 \
  --set operator.replicas=1 \
  --set kubeProxyReplacement=strict \
  --set k8sServiceHost=$API_SERVER_IP \
  --set k8sServicePort=$API_SERVER_PORT \
  --wait
```

We can now compare the running Pods on `cluster2` and `kubeless` in the `kube-system` Namespace.

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

On `kubeless` there is no `kube-proxy` Pod anymore:

```bash
kubectl --context kubeless -n kube-system get pod
```

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


## Task {{% param sectionnumber %}}.2: Deploy our simple app again to the new cluster

As this is a new cluster we want to deploy our `simple-app.yaml` from lab 03 again to run some experiments. Run the following command using the `simple-app.yaml` from lab 03:

```bash
kubectl create -f simple-app.yaml
```

Now let us redo the task from lab 03.

Let's make life again a bit easier by storing the Pod's name into an environment variable so we can reuse it later again:

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
echo ${FRONTEND}
NOT_FRONTEND=$(kubectl get pods -l app=not-frontend -o jsonpath='{.items[0].metadata.name}')
echo ${NOT_FRONTEND}
```

Then execute

```bash
kubectl exec -ti ${FRONTEND} -- curl -I --connect-timeout 5 backend:8080
```

and

```bash
kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
```

You see that altought we have no `kube-proxy` running, the backend service can still be reached.

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


## Task {{% param sectionnumber %}}.3: Cleanup

We don't need `kubeless` anymore. You can stop `kubeless` with:

```bash
minikube stop -p kubeless
minikube delete -p kubeless
```

to free up resources and speed up things.
