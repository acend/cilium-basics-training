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
minikube delete
minikube start --network-plugin=cni --cni=false --kubernetes-version=1.23.0
```

This will install a new Kubernetes Cluster without any Container Network Interface (CNI). The CNI will be installed later in the labs.

Minikube added a new context into your Kubernetes config file and set this as your default context. Check it with the following command:

```bash
kubectl config current-context
```

This should show 'minikube'. Now check that everything is up and running using the following command:

```bash
kubectl get node           
```

This should produce an output similar to the following:

```
NAME       STATUS   ROLES                  AGE   VERSION
minikube   Ready    control-plane,master   86s   v1.23.0
```
Depending on your minikube version and environment your node might stay NotReady because no CNI exists. It will become ready after the cilium installation.

Check if all pods are running with

```bash
kubectl get pod -A
```

which produces the following output

```
NAMESPACE     NAME                               READY   STATUS              RESTARTS   AGE
kube-system   coredns-558bd4d5db-wdlrh           0/1     ContainerCreating   0          3m1s
kube-system   etcd-minikube                      1/1     Running             0          3m7s
kube-system   kube-apiserver-minikube            1/1     Running             0          3m16s
kube-system   kube-controller-manager-minikube   1/1     Running             0          3m7s
kube-system   kube-proxy-9bjbq                   1/1     Running             0          3m1s
kube-system   kube-scheduler-minikube            1/1     Running             0          3m7s
kube-system   storage-provisioner                1/1     Running             0          3m11s
```


{{% alert title="Note" color="primary" %}}
Depending on your minikube version coredns might start or not. In both cases we can proceed.
You should not see any CNI related pods!
{{% /alert %}}
