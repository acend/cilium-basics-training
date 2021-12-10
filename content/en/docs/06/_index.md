---
title: "6. Troubleshooting"
weight: 6
sectionnumber: 6
---


## Component & Cluster Health

An initial overview of Cilium can be retrieved by listing all pods to verify whether all pods have the status `Running`:

```bash
kubectl -n kube-system get pods -l k8s-app=cilium
```

```
NAME           READY     STATUS    RESTARTS   AGE
cilium-2hq5z   1/1       Running   0          4d
cilium-6kbtz   1/1       Running   0          4d
cilium-klj4b   1/1       Running   0          4d
cilium-zmjj9   1/1       Running   0          4d
```

If Cilium encounters a problem that it cannot recover from, it will automatically report the failure state via cilium status which is regularly queried by the Kubernetes liveness probe to automatically restart Cilium pods. If a Cilium pod is in state CrashLoopBackoff then this indicates a permanent failure scenario.

If a particular Cilium pod is not in running state, the status and health of the agent on that node can be retrieved by running cilium status in the context of that pod:

```bash
kubectl -n kube-system exec cilium-2hq5z -- cilium status
```

```

```

Detailed information about the status of Cilium can be inspected with the cilium status --verbose command. Verbose output includes detailed IPAM state (allocated addresses), Cilium controller status, and details of the Proxy status.


## Logs

To retrieve log files of a cilium pod, run:

```bash
kubectl -n kube-system logs --timestamps <pod-name>
```

The `<pod-name>` can be determined with the following command and selecting the name of one of the pods:

```bash
kubectl -n kube-system get pods -l k8s-app=cilium
```

If the cilium pod was already restarted due to the liveness problem after encountering an issue, it can be useful to retrieve the logs of the pod before the last restart:

```bash
kubectl -n kube-system logs --timestamps -p <pod-name>
```

