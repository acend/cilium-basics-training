---
title: "Exporting Events"
weight: 124
OnlyWhenNot: techlab
---


## {{% task %}} Export Network Events

Edit your `cilium-enterprise-values.yaml` file and include export-file-path field to export network events:

```yaml
cilium:
  (...)
  extraConfig:
    # Enable network event export
    export-file-path: "/var/run/cilium/hubble/hubble.log"
  (...)
hubble-enterprise:
  enabled: true

```

Then, run helm upgrade command to apply the new configuration:

```bash
helm upgrade cilium-enterprise isovalent/cilium-enterprise --version {{% param "ciliumVersion.enterprise" %}}
  --namespace kube-system -f cilium-enterprise-values.yaml --wait
```

and restart cilium daemonset for the new filters to take effect:

```yaml
kubectl rollout restart -n kube-system ds/cilium
```


## {{% task %}} Export Process Events


Edit your `cilium-enterprise-values.yaml` file and include exportFilename field to export process events:

```yaml
cilium:
  (...)
hubble-enterprise:
  enabled: true
  enterprise:
    # Enable process event export
    exportFilename: "fgs.log"
  (...)
```

Then, run helm upgrade command to apply the new configuration:

```bash
helm upgrade cilium-enterprise isovalent/cilium-enterprise --version {{% param "ciliumVersion.enterprise" %}}
  --namespace kube-system -f cilium-enterprise-values.yaml --wait
```


## {{% task %}} Observe Exported Events

Run the following command to observe exported events in `export-stdout` container logs:

```bash
kubectl logs -n kube-system -l app.kubernetes.io/name=hubble-enterprise -c export-stdout -f
```

Those exported events can now be sent to Splunk, Elasticsearch or similar.
