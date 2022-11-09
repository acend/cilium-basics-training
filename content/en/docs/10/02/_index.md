---
title: "Kubernetes Without kube-proxy"
weight: 102
OnlyWhenNot: techlab
---

In this lab, we are going to provision a new Kubernetes cluster without `kube-proxy` to use Cilium as a full replacement for it.


## {{% task %}} Deploy a new Kubernetes Cluster without `kube-proxy`


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


## {{% task %}} Deploy Cilium and enable the Kube Proxy replacement

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

{{% alert title="Note" color="primary" %}}
Having a cluster running with kubeProxyReplacement set to partial breaks other minikube clusters running on the same host. If you still want to play around with `cluster1` after this chapter, you need to reboot your maching and start only cluster1 with `minikube start --profile cluster1`
{{% /alert %}}

We can now compare the running Pods on `cluster1` and `kubeless` in the `kube-system` Namespace.

```bash
kubectl --context cluster1 -n kube-system get pod
```

Here we see the running `kube-proxy` pod:

```
NAME                                    READY   STATUS    RESTARTS       AGE
cilium-operator-cb65bcb9b-cnxnq          1/1     Running   0             19m
cilium-tq9kk                             1/1     Running   0             8m42s
clustermesh-apiserver-67fd99fd9b-x2svr   2/2     Running   0             61m
coredns-6d4b75cb6d-fd6vk                 1/1     Running   1 (82m ago)   97m
etcd-cluster1                            1/1     Running   1 (82m ago)   98m
hubble-relay-84b4ddb556-nvftg            1/1     Running   0             19m
hubble-ui-579fdfbc58-t6xst               2/2     Running   0             19m
kube-apiserver-cluster1                  1/1     Running   1 (81m ago)   98m
kube-controller-manager-cluster1         1/1     Running   1 (82m ago)   98m
kube-proxy-5j84l                         1/1     Running   1 (82m ago)   97m
kube-scheduler-cluster1                  1/1     Running   1 (81m ago)   98m
storage-provisioner                      1/1     Running   2 (82m ago)   98m
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


## {{% task %}} Deploy our simple app again to the new cluster

As this is a new cluster we want to deploy our `simple-app.yaml` from lab 03 again to run some experiments. Run the following command using the `simple-app.yaml` from lab 03:

```bash
kubectl apply -f simple-app.yaml
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
