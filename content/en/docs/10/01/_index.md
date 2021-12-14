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
cilium install --config cluster-pool-ipv4-cidr=10.3.0.0/16  --kube-proxy-replacement strict --cluster-name cluster3 --cluster-id 3
```
