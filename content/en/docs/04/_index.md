---
title: "4. Metrics"
weight: 4
sectionnumber: 4
---

Hubble and its UI did allow us to see traffic flow inside our cluster. Both Cilium and Hubble can be configured to serve Prometheus metrics independently of each other.
With metrics displayed in grafana or another UI, we can get a quick overview of our cluster state and its traffic.

Cilium metrics show us the state of Cilium itself, namely of the cilium-agent, cilium-envoy, and cilium-operator processes.
Hubble metrics give us information about the traffic of our applications.


## Task {{% param sectionnumber %}}.1:  Enable metrics

```bash
helm upgrade -i cilium cilium/cilium \
   --namespace kube-system \
   --reuse-values \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.enabled=true \
   --set hubble.metrics.enabled="{dns,drop:destinationContext=pod;sourceContext=pod,tcp,flow,port-distribution,icmp,http:destinationContext=pod}"
```


### Verify cilium metrics

We now see that the cilium agent has different metric endpoints:

* hubble port 9091
* cilium agent port 9090
* cilium envoy port 9095

```bash
CILIUM_AGENT_IP=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].status.hostIP}")
kubectl run -n kube-system -it --env="CILIUM_AGENT_IP=${CILIUM_AGENT_IP}" --rm curl --image=curlimages/curl -- sh
echo ${CILIUM_AGENT_IP}
curl -s ${CILIUM_AGENT_IP}:9090/metrics | grep cilium_nodes_all_num #show total number of cilium nodes
curl -s ${CILIUM_AGENT_IP}:9091/metrics | grep hubble_tcp_flags_total # show total number of TCP flags
exit
```
{{% alert title="Note" color="primary" %}}
The Cilium agent pods run as daemonset on the HostNetwork. If you have curl installed you can also directly call a node.
```bash
NODE=$(kubectl get nodes --selector=kubernetes.io/role!=master -o jsonpath={.items[*].status.addresses[?\(@.type==\"InternalIP\"\)].address})
curl -s $NODE:9090/metrics | grep cilium_nodes_all_num
```
{{% /alert %}}

{{% alert title="Note" color="primary" %}}
It is not yet possible to get metrics from Cilium envoy (port 9095). Envoy only starts on a node if there is at least one pod with a layer 7 networkpolicy.
{{% /alert %}}

You should see now an output like this.
```bash
If you don't see a command prompt, try pressing enter.
echo ${CILIUM_AGENT_IP}
192.168.49.2
/ $ curl -s ${CILIUM_AGENT_IP}:9090/metrics | grep cilium_nodes_all_num #show total number of cilium nodes
# HELP cilium_nodes_all_num Number of nodes managed
# TYPE cilium_nodes_all_num gauge
cilium_nodes_all_num 1
/ $ curl -s ${CILIUM_AGENT_IP}:9091/metrics | grep hubble_tcp_flags_total # show total number of TCP flags
# HELP hubble_tcp_flags_total TCP flag occurrences
# TYPE hubble_tcp_flags_total counter
hubble_tcp_flags_total{family="IPv4",flag="FIN"} 2704
hubble_tcp_flags_total{family="IPv4",flag="RST"} 388
hubble_tcp_flags_total{family="IPv4",flag="SYN"} 1609
hubble_tcp_flags_total{family="IPv4",flag="SYN-ACK"} 1549
```


## Task {{% param sectionnumber %}}.2:  Store and visualize metrics

Install prometheus grafana into cilium-monitoring namespace to store and visualize cilium and hubble metrics.
```bash
kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/v1.11/examples/kubernetes/addons/prometheus/monitoring-example.yaml
```

Make sure grafana and prometheus pods are up and running before continuing with the next step.

```bash
kubectl -n cilium-monitoring get pod
```
you should see both Pods in state `Running`:

```
NAME                          READY   STATUS    RESTARTS   AGE
grafana-6c7d4c9fd8-2xdp2      1/1     Running   0          41s
prometheus-55777f54d9-hkpkq   1/1     Running   0          41s
```


Generate some traffic for some minutes in the background
```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
i=0; while [ $i -le 300 ]; do kubectl exec -ti ${FRONTEND} -- curl -Is backend:8080; sleep 1; ((i++)); done &
```


In a second terminal access Grafana with kubectl proxy-forward
```bash
kubectl -n cilium-monitoring port-forward service/grafana --address 0.0.0.0 --address :: 3000:3000 &
```

Now open your browser and go to http://localhost:3000/dashboards. Open the Hubble Dashboard and browse through its graphs. For a better view, you can change the timespan to the last 5 minutes. Verify that you see the generated traffic under Network, Forwarded vs Dropped Traffic.

Not all graphs have data available. This is because we have not yet used network policies or any layer 7 components. This will be done in the later chapters.

In grafana use the left side menu: `Dashboard`, click on `Manage`, then click on `Cilium Metrics`. Here we see information about Cilium itself. Again not all graphs contain data as we have not used all features of cilium yet.

Browse through the graphs, try to find the number of IPs allocated and the number of cilium endpoints.

You can close the grafana browser tab or leave it open for later use.
