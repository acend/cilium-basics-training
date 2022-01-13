---
title: "3.1 Hubble UI"
weight: 31
sectionnumber: 3.1
---

Not only does Hubble allow to inspect flows from the command line, it also allows us to see them in real time on a graphical service map via Hubble UI. Again, this also is an optional component that is disabled by default.


## Task {{% param sectionnumber %}}.1: Enable the Hubble UI component

Enable the optional Hubble UI component with Helm looks like this:

```
```bash
helm repo add cilium https://helm.cilium.io/
helm upgrade -i cilium cilium/cilium --version 1.11.0 \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDR=10.1.0.0/16 \
  --reuse-values \
  --set hubble.ui.enabled=true \
  --wait
```

{{% alert title="Note" color="primary" %}}
When using the `cilium` CLI you can execute the following command to enable the Hubble UI:

```bash
cilium hubble enable --ui
```
{{% /alert %}}

Take a look at the pods again to see what happened under the hood:

```bash
kubectl get pods -A
```

We see, there is again a new Pod running for the `hubble-ui` component.

```
default       backend-56787b4bd7-dmzdh           1/1     Running   0          138m
default       frontend-7cbdcb86fd-gdb4q          1/1     Running   0          138m
default       not-frontend-5cf6d96558-gj4np      1/1     Running   0          138m
kube-system   cilium-nz77s                       1/1     Running   0          8m48s
kube-system   cilium-operator-8dd4dc946-n9ght    1/1     Running   0          174m
kube-system   coredns-558bd4d5db-xzvc9           1/1     Running   0          174m
kube-system   etcd-minikube                      1/1     Running   0          174m
kube-system   hubble-relay-f6d85866c-csthd       1/1     Running   0          25m
kube-system   hubble-ui-9b6d87f-zpvgp            3/3     Running   0          8m16s
kube-system   kube-apiserver-minikube            1/1     Running   0          174m
kube-system   kube-controller-manager-minikube   1/1     Running   0          174m
kube-system   kube-proxy-bqs4d                   1/1     Running   0          174m
kube-system   kube-scheduler-minikube            1/1     Running   0          174m
kube-system   storage-provisioner                1/1     Running   1          174m
```

Cilium agents are restarting, and a new Hubble UI pod is now present on top of the Hubble Relay pod. As above, we can wait for Cilium and Hubble to be ready by running:

```bash
cilium status --wait
```

```
cilium status --wait
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         OK
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

Deployment        hubble-ui          Desired: 1, Ready: 1/1, Available: 1/1
DaemonSet         cilium             Desired: 1, Ready: 1/1, Available: 1/1
Deployment        cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
Deployment        hubble-relay       Desired: 1, Ready: 1/1, Available: 1/1
Containers:       cilium-operator    Running: 1
                  hubble-relay       Running: 1
                  hubble-ui          Running: 1
                  cilium             Running: 1
Cluster Pods:     6/6 managed by Cilium
Image versions    cilium             quay.io/cilium/cilium:v1.10.5: 1
                  cilium-operator    quay.io/cilium/operator-generic:v1.10.5: 1
                  hubble-relay       quay.io/cilium/hubble-relay:v1.10.5: 1
                  hubble-ui          quay.io/cilium/hubble-ui:v0.7.9@sha256:e0e461c680ccd083ac24fe4f9e19e675422485f04d8720635ec41f2ba9e5562c: 1
                  hubble-ui          quay.io/cilium/hubble-ui-backend:v0.7.9@sha256:632c938ef6ff30e3a080c59b734afb1fb7493689275443faa1435f7141aabe76: 1
                  hubble-ui          docker.io/envoyproxy/envoy:v1.18.2@sha256:e8b37c1d75787dd1e712ff389b0d37337dc8a174a63bed9c34ba73359dc67da7: 1
```


And then check Hubble status:

```bash
hubble status
```

{{% alert title="Note" color="primary" %}}
Note: our earlier cilium hubble port-forward should still be running (can be checked by running jobs or `ps aux | grep "cilium hubble port-forward"`). If it does not, hubble status will fail and we have to run it again:

```bash
cilium hubble port-forward&
hubble status
```

{{% /alert %}}


To start Hubble UI execute

```bash
cilium hubble ui
```

The browser should automatically open http://localhost:12000/ (open it manually if not).

We can then access the graphical service map by selecting our default namespace

![Hubble UI Choose Namespace](../cilium_choose_ns.png)

Then you should see a spinning circle and the message "Waiting for service map data..."

Lets generate some network activity again:

```bash
for i in {1..10}; do
  kubectl exec -ti ${FRONTEND} -- curl -I --connect-timeout 5 backend:8080
  kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
done
```

and then you should see a service map in the hubble ui

![Hubble UI - Service Map](../hubble_ui_servicemap.png)

and also a Table with the already familiar Flow output previously seen in the `hubble observe` command:

![Hubble UI - Service Map](../hubble_ui_flows.png)

Hubble flows are displayed in real time at the bottom, with a visualization of the namespace objects in the center. Click on any flow, and click on any property from the right side panel: notice that the filters at the top of the UI have been updated accordingly.

Let's run a connectivity test again and see what happens in Hubble UI in the cilium-test namespace

```bash
cilium connectivity test
```

We can see that Hubble UI is not only capable of displaying flows within a namespace, it also helps visualizing flows going in or out.

![Hubble UI - Connectivity Test](../cilium_hubble_connectivity_test.png)
