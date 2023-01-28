---
title: "eBPF"
weight: 12
OnlyWhenNot: techlab
---

To deepen our understanding of eBPF we will write and compile a small eBPF app:


## {{% task %}} Hello World

ebpf-go is a pure Go library that provides utilities for loading, compiling, and debugging eBPF programs written by the cilium project.

We will use this library and add our own hello world app as an example to it:

```bash
git clone https://github.com/cilium/ebpf.git
cd ebpf/
git checkout v0.9.3
cd examples
mkdir helloworld
cd helloworld
```

In the `helloworld` directory create two files named `helloworld.bpf.c` (eBPF code) and `helloworld.go` (loading, user side):

helloworld.bpf.c:
{{< readfile file="/content/en/docs/12/helloworld/helloworld.bpf.c" code="true" lang="c" >}}

helloworld.go:
{{< readfile file="/content/en/docs/12/helloworld/helloworld.go" code="true" lang="go" >}}


To compile the C code into ebpf bytecode with the corresponding Go source files we use a tool named bpf2go along with clang.
For a stable outcome we use the toolchain inside a docker container:

```bash
docker pull "ghcr.io/cilium/ebpf-builder:1666886595"
docker run -it --rm -v "$(pwd)/../..":/ebpf \
   -w /ebpf/examples/helloworld \
  --env MAKEFLAGS \
  --env CFLAGS="-fdebug-prefix-map=/ebpf=." \
  --env HOME="/tmp" \
  "ghcr.io/cilium/ebpf-builder:1666886595" /bin/bash
```
Now in the container we generate the ELF and go files:
```bash
GOPACKAGE=main go run github.com/cilium/ebpf/cmd/bpf2go -cc clang-14 -cflags '-O2 -g -Wall -Werror' bpf helloworld.bpf.c -- -I../headers
```

Let us examine the newly created files: `bpf_bpfel.go`/`bpf_bpfeb.go` contain the go code for the user state side of our app.
The `bpf_bpfel.o`/`bpf_bpfeb.o` files are ELF files and can be examined using readelf:

```bash
readelf --section-details --headers bpf_bpfel.o
```

We see two things:

* that Machine reads "Linux BPF" and
* our tracepoint sys_enter_execve in the sections part (tracepoint/syscalls/sys_enter_execve).


{{% alert title="Note" color="primary" %}}
There are always two files created: bpf_bpfel.o for little endian systems (like x86) and bpfen.o for big endian systems.
{{% /alert %}}

Now we have everything in place to build our app:


```bash
go mod tidy
go build helloworld.go bpf_bpfel.go
exit #exit container
```

We run our app in the background, cat tracepipe in the background and see if we get a hello world for every command:

```bash
sudo cat /sys/kernel/debug/tracing/trace_pipe &
sudo ./helloworld &
ls #example command
```

Now we can see, that for each programm called in linux, our code is executed and writes "Hello world" to trace_pipe.

Close now the running apps in the background

```bash
sudo kill $(jobs -p)
```
