---
title: "2. Install Cilium"
weight: 2
sectionnumber: 2
---


Cilium can be installed using multiple ways:

* Cilium CLI
* Using Helm

In this lab we are going to use the [Cilium command line](https://github.com/cilium/cilium-cli/) tool (Cilium CLI)


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

