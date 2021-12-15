---
title: "9. Metrics"
weight: 9
sectionnumber: 9
---

{{% alert title="Note" color="primary" %}}
This lab should be done on your `cluster1`, make sure to switch to `cluster1` with `minikube profile cluster1`
{{% /alert %}}

Cilium and Hubble can both be configured to serve Prometheus metrics independently of each other.

Cilium metrics provide insights into the state of Cilium itself, namely of the cilium-agent, cilium-envoy, and cilium-operator processes.
Hubble metrics provide insight into the services.


## Task {{% param sectionnumber %}}.1:  Enable metrics

```bash
helm upgrade -i cilium cilium/cilium \
   --namespace kube-system \
   --reuse-values \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.enabled=true \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"
```


### Verify cilium metrics

We now see that the cilium agent has two different metric endpoints:

* hubble port 6942
* cilium port 9090

```bash
CILIUM_AGENT_IP=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].status.hostIP}")
kubectl run -n kube-system -it --env="CILIUM_AGENT_IP=${CILIUM_AGENT_IP}" --rm curl --image=curlimages/curl -- sh
echo ${CILIUM_AGENT_IP}
curl -s ${CILIUM_AGENT_IP}:6942/metrics
curl -s ${CILIUM_AGENT_IP}:9090/metrics
exit
```


## Task {{% param sectionnumber %}}.2:  Visualize metrics

Install grafana into cilium-monitoring namespace to visualize cilium and hubble metrics.
```bash
kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/v1.11/examples/kubernetes/addons/prometheus/monitoring-example.yaml
```

Generate some traffic
```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
kubectl exec -ti ${FRONTEND} -- curl -Is backend:8080
```


Access Grafana with kubectl proxy-forward
```bash
kubectl -n cilium-monitoring port-forward service/grafana --address 0.0.0.0 --address :: 3000:3000 &
```

Now open your browser and go to http://localhost:3000/dashboards. After you have finished you can stop port-forwarding with `kill %1`
