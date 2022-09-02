---
title: "Setup"
weight: 1
type: docs
menu:
  main:
    weight: 1
---

This training can be done in two ways:

* on a local machine (proceed with [Local Machine Setup](#local-machine-setup))
* on a provided virtual machine using the webshell (proceed with [Webshell access](#webshell-access))


## Local machine setup


### Technical prerequisites

To run this training on your local machine please make sure the following requirements are met:

* Operating System: Linux with Kernel >= 4.9.17 or MacOS
* Docker [installed](https://docs.docker.com/get-docker/)
* kubectl >= 1.24 [installed](https://kubernetes.io/docs/tasks/tools/#kubectl)
* minikube >= 1.26 installed
* helm installed
* Minimum 8GB RAM

A note on Windows with WSL2: As of August 2022 the default kernel in WSL is missing some Netfilter modules. You can compile it [yourself](https://github.com/cilium/cilium/issues/17745#issuecomment-1004299480), but the training staff cannot give you any support with cluster related issues.


## Install minikube

This training uses [minikube](https://minikube.sigs.k8s.io/docs/) to provide a Kubernetes Cluster.

Check the [minikube start Guide](https://minikube.sigs.k8s.io/docs/start/) for instructions on how to install minikube on your system. If you are using the provided virtual machine minikube is already installed.


## Install helm

For a complete overview refer to the helm installation [website](https://helm.sh/docs/intro/install/). If you have helm 3 already installed you can skip this step.

Use your package manager (`apt`, `yum`, `brew` etc), download the [latest Release](https://github.com/helm/helm/releases) or use the following command to install [helm](https://helm.sh/docs/intro/install/) helm:

```bash
curl -s https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```


## Webshell access

Your trainer will give you the necessary details.
