---
title: "2.1 Upgrade Cilium"
weight: 21
sectionnumber: 2.1
---

In the previous lab we intentionally installed version `v10.5.0` of Cilium. In this lab we show you how to upgrade the installation.


## Task {{% param sectionnumber %}}.1: Running pre-flight check

When rolling out an upgrade with Kubernetes, Kubernetes will first terminate the pod followed by pulling the new image version and then finally spin up the new image. In order to reduce the downtime of the agent and to prevent `ErrImagePull` errors during upgrade, the pre-flight check pre-pulls the new image version. If you are running in "Kubernetes Without kube-proxy" mode you must also pass on the Kubernetes API Server IP and / or the Kubernetes API Server Port when generating the cilium-preflight.yaml file.

```bash
helm install cilium-preflight cilium/cilium --version 1.11.0 \
  --namespace=kube-system \
  --set preflight.enabled=true \
  --set agent=false \
  --set operator.enabled=false
```


## Task {{% param sectionnumber %}}.2: Clean up pre-flight check

Once the number of `READY` for the preflight DaemonSet is the same as the number of cilium pods running and the preflight Deployment is marked as `READY 1/1` you can delete the cilium-preflight and proceed with the upgrade.

```bash
helm delete cilium-preflight --namespace=kube-system
```


## Task {{% param sectionnumber %}}.3: Upgrading Cilium

During normal cluster operations, all Cilium components should run the same version. Upgrading just one of them (e.g., upgrading the agent without upgrading the operator) could result in unexpected cluster behavior. The following steps will describe how to upgrade all of the components from one stable release to a later stable release.

When upgrading from one minor release to another minor release, for example 1.x to 1.y, it is recommended to upgrade to the latest patch release for a Cilium release series first. The latest patch releases for each supported version of Cilium are [here](https://github.com/cilium/cilium#stable-releases). Upgrading to the latest patch release ensures the most seamless experience if a rollback is required following the minor release upgrade. The upgrade guides for previous versions can be found for each minor version at the bottom left corner.

Helm can be used to either upgrade Cilium directly or to generate a new set of YAML files that can be used to upgrade an existing deployment via kubectl. By default, Helm will generate the new templates using the default values files packaged with each new release. You still need to ensure that you are specifying the equivalent options as used for the initial deployment, either by specifying them at the command line or by committing the values to a YAML file.

To minimize datapath disruption during the upgrade, the `upgradeCompatibility` option should be set to the initial Cilium version which was installed in this cluster. Valid options are:

```bash
helm upgrade -i cilium cilium/cilium --version 1.11.0 \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDR=10.1.0.0/16 \
  --set cluster.name=cluster1 \
  --set cluster.id=1 \
  --set operator.replicas=1 \
  --set upgradeCompatibility=1.10 \
  --wait
```
{{% alert title="Note" color="primary" %}}
When upgrading from one minor release to another minor release using helm upgrade, do not use Helm’s `--reuse-values` flag. The  `--reuse-values` flag ignores any newly introduced values present in the new release and thus may cause the Helm template to render incorrectly. Instead, if you want to reuse the values from your existing installation, save the old values in a values file, check the file for any renamed or deprecated values, and then pass it to the `helm upgrade` command as described above. You can retrieve and save the values from an existing installation with the following command:

```bash
helm get values cilium --namespace=kube-system -o yaml > old-values.yaml
```

The `--reuse-values` flag may only be safely used if the Cilium chart version remains unchanged, for example when `helm upgrade` is used to apply configuration changes without upgrading Cilium.
{{% /alert %}}


## Rolling Back

Occasionally, it may be necessary to undo the rollout because a step was missed or something went wrong during upgrade. To undo the rollout run:

```bash
helm history cilium --namespace=kube-system
helm rollback cilium [REVISION] --namespace=kube-system
```

This will revert the latest changes to the Cilium DaemonSet and return Cilium to the state it was in prior to the upgrade.
