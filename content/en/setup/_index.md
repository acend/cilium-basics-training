---
title: "Setup"
weight: 1
type: docs
menu:
  main:
    weight: 1
---

## Install Minikube

This training uses [Minikube](https://minikube.sigs.k8s.io/docs/) to provide a Kubernetes Cluster on your local machine.

Check the [Minikube start Guide](https://minikube.sigs.k8s.io/docs/start/) for instructions on how to install minikube on your system.


## Install a Kubernetes Cluster

We are going to spin up a new Kubernetes cluster with the following command:

```bash
minikube start --network-plugin=cni --cni=false --kubernetes-version=1.23.0 -p cluster1 
```

This will install a new Kubernetes Cluster without any Container Network Interface (CNI). The CNI will be installed later in the labs.

Minikube added a new conext into your Kubernetes config file and set this as your default context.

Check that everythis is up and running using the following command:

```bash
kubectl get node           
```

This should produce an output similar to the following:

```
NAME       STATUS   ROLES                  AGE   VERSION
minikube   Ready    control-plane,master   75s   v1.22.3
```

Check also

```bash
kubectl get pod -A
```

which prodocues the following output

```
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE
kube-system   coredns-78fcd69978-kqs2c           1/1     Running   0          2m1s
kube-system   etcd-minikube                      1/1     Running   0          2m13s
kube-system   kube-apiserver-minikube            1/1     Running   0          2m13s
kube-system   kube-controller-manager-minikube   1/1     Running   0          2m13s
kube-system   kube-proxy-9jxxj                   1/1     Running   0          2m1s
kube-system   kube-scheduler-minikube            1/1     Running   0          2m13s
kube-system   storage-provisioner                1/1     Running   0          2m12s
```


{{% alert title="Note" color="primary" %}}
You should not see any CNI related pods!
{{% /alert %}}
