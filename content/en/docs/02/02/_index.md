---
title: "Upgrade Cilium"
weight: 22
OnlyWhenNot: techlab
---

In the previous lab, we intentionally installed version `v{{% param "ciliumVersion.preUpgrade" %}}` of Cilium. In this lab, we show you how to upgrade this installation.


## {{% task %}} Running pre-flight check

When rolling out an upgrade with Kubernetes, Kubernetes will first terminate the Pod followed by pulling the new image version and then finally spin up the new image. In order to reduce the downtime of the agent and to prevent `ErrImagePull` errors during the upgrade, the pre-flight check pre-pulls the new image version. If you are running in "Kubernetes Without kube-proxy" mode you must also pass on the Kubernetes API Server IP and/or the Kubernetes API Server Port when generating the `cilium-preflight.yaml` file.

```bash
helm install cilium-preflight cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace=kube-system \
  --set preflight.enabled=true \
  --set agent=false \
  --set operator.enabled=false \
  --wait
```


## {{% task %}} Clean up pre-flight check


To check the preflight Pods we check if the pods are `READY` using:

```bash
kubectl get pod -A | grep cilium-pre-flight
```

and you should get an output like this:

```
kube-system   cilium-pre-flight-check-84f67b54f6-hz57g   1/1     Running   0               63s
kube-system   cilium-pre-flight-check-skglp              1/1     Running   0               63s
```

The pods are `READY` with a value of `1/1` and therefore we can delete the `cilium-preflight` release again with:

```bash
helm delete cilium-preflight --namespace=kube-system
```


## {{% task %}} Upgrading Cilium

During normal cluster operations, all Cilium components should run the same version. Upgrading just one of them (e.g., upgrading the agent without upgrading the operator) could result in unexpected cluster behavior. The following steps will describe how to upgrade all of the components from one stable release to a later stable release.

When upgrading from one minor release to another minor release, for example 1.x to 1.y, it is recommended to upgrade to the latest patch release for a Cilium release series first. The latest patch releases for each supported version of Cilium are [here](https://github.com/cilium/cilium#stable-releases). Upgrading to the latest patch release ensures the most seamless experience if a rollback is required following the minor release upgrade. The upgrade guides for previous versions can be found for each minor version at the bottom left corner.

Helm can be used to either upgrade Cilium directly or to generate a new set of YAML files that can be used to upgrade an existing deployment via kubectl. By default, Helm will generate the new templates using the default values files packaged with each new release. You still need to ensure that you are specifying the equivalent options as used for the initial deployment, either by specifying them at the command line or by committing the values to a YAML file.

To minimize datapath disruption during the upgrade, the `upgradeCompatibility` option should be set to the initial Cilium version which was installed in this cluster.

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDRList={10.1.0.0/16} \
  --set cluster.name=cluster1 \
  --set cluster.id=1 \
  --set operator.replicas=1 \
  --set kubeProxyReplacement=disabled \
  --set upgradeCompatibility=1.11 \
  --wait
```
{{% alert title="Note" color="primary" %}}
When upgrading from one minor release to another minor release using `helm upgrade`, do not use Helm’s `--reuse-values` flag. The  `--reuse-values` flag ignores any newly introduced values present in the new release and thus may cause the Helm template to render incorrectly. Instead, if you want to reuse the values from your existing installation, save the old values in a values file, check the file for any renamed or deprecated values, and then pass it to the `helm upgrade` command as described above. You can retrieve and save the values from an existing installation with the following command:

```bash
helm get values cilium --namespace=kube-system -o yaml > old-values.yaml
```

The `--reuse-values` flag may only be safely used if the Cilium chart version remains unchanged, for example when `helm upgrade` is used to apply configuration changes without upgrading Cilium.
{{% /alert %}}


## {{% task %}} Explore your installation after the upgrade

We can run:

```bash
cilium status --wait
```

again to verify the upgrade to the new version succeded

```
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         disabled
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

DaemonSet         cilium             Desired: 1, Ready: 1/1, Available: 1/1
Deployment        cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
Containers:       cilium-operator    Running: 1
                  cilium             Running: 1
Cluster Pods:     1/1 managed by Cilium
Image versions    cilium             quay.io/cilium/cilium:v{{% param "ciliumVersion.postUpgrade" %}}:: 1
                  cilium-operator    quay.io/cilium/operator-generic:v{{% param "ciliumVersion.postUpgrade" %}}@: 1
```

And we see the right version in the `cilium` and `cilium-operator` images.


### Nice to know

In Cilium release 1.11.0 automatic mount of eBPF maps in the host filesystem were enabled. These eBPF maps are basically very efficient key-value stores used by Cilium. Having them mounted in the filesystem, allows the datapath to continue operating even if the `cilium-agent` is restarting. We can verify that Cilium created global traffic control eBPF maps on the node in /sys/fs/bpf/tc/globals/:

```bash
docker exec cluster1 ls /sys/fs/bpf/tc/globals/
```


## Rolling Back

Occasionally, it may be necessary to undo the rollout because a step was missed or something went wrong during the upgrade. To undo the rollout run:

```
helm history cilium --namespace=kube-system
# helm rollback cilium [REVISION] --namespace=kube-system
```

This will revert the latest changes to the Cilium DaemonSet and return Cilium to the state it was in prior to the upgrade.
