---
title: "Metrics"
weight: 4
OnlyWhenNot: techlab
---

With metrics displayed in Grafana or another UI, we can get a quick overview of our cluster state and its traffic.

Both Cilium and Hubble can be configured to serve Prometheus metrics independently of each other. Cilium metrics show us the state of Cilium itself, namely of the `cilium-agent`, `cilium-envoy`, and `cilium-operator` processes.
Hubble metrics on the other hand give us information about the traffic of our applications.


## {{% task %}} Enable metrics

We start by enabling different metrics, for dropped and HTTP traffic we also want to have metrics specified by pod.

```bash
helm upgrade -i cilium cilium/cilium --version {{% param "ciliumVersion.postUpgrade" %}} \
  --namespace kube-system \
  --set ipam.operator.clusterPoolIPv4PodCIDRList={10.1.0.0/16} \
  --set cluster.name=cluster1 \
  --set cluster.id=1 \
  --set operator.replicas=1 \
  --set upgradeCompatibility=1.11 \
  --set kubeProxyReplacement=disabled \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --set hubble.ui.enabled=true \
  `# enable metrics:` \
  --set prometheus.enabled=true \
  --set operator.prometheus.enabled=true \
  --set hubble.metrics.enabled="{dns,drop:destinationContext=pod;sourceContext=pod,tcp,flow,port-distribution,icmp,http:destinationContext=pod}"
```


### Verify Cilium metrics

We now verify that the Cilium agent has different metric endpoints exposed and list some of them:

* hubble port 9965
* cilium agent port 9962
* cilium envoy port 9095

```bash
CILIUM_AGENT_IP=$(kubectl get pod -n kube-system -l k8s-app=cilium -o jsonpath="{.items[0].status.hostIP}")
kubectl run -n kube-system -it --env="CILIUM_AGENT_IP=${CILIUM_AGENT_IP}" --rm curl --image=curlimages/curl -- sh
```
```bash
echo ${CILIUM_AGENT_IP}
curl -s ${CILIUM_AGENT_IP}:9962/metrics | grep cilium_nodes_all_num #show total number of cilium nodes
curl -s ${CILIUM_AGENT_IP}:9965/metrics | grep hubble_tcp_flags_total # show total number of TCP flags
exit
```
You should see now an output like this.
```
If you don't see a command prompt, try pressing enter.
echo ${CILIUM_AGENT_IP}
192.168.49.2
/ $ curl -s ${CILIUM_AGENT_IP}:9962/metrics | grep cilium_nodes_all_num #show total number of cilium nodes
# HELP cilium_nodes_all_num Number of nodes managed
# TYPE cilium_nodes_all_num gauge
cilium_nodes_all_num 1
/ $ curl -s ${CILIUM_AGENT_IP}:9965/metrics | grep hubble_tcp_flags_total # show total number of TCP flags
# HELP hubble_tcp_flags_total TCP flag occurrences
# TYPE hubble_tcp_flags_total counter
hubble_tcp_flags_total{family="IPv4",flag="FIN"} 2704
hubble_tcp_flags_total{family="IPv4",flag="RST"} 388
hubble_tcp_flags_total{family="IPv4",flag="SYN"} 1609
hubble_tcp_flags_total{family="IPv4",flag="SYN-ACK"} 1549
```
{{% alert title="Note" color="primary" %}}
The Cilium agent pods run as DaemonSet on the HostNetwork. This means you could also directly call a node.
```bash
NODE=$(kubectl get nodes --selector=kubernetes.io/role!=master -o jsonpath={.items[*].status.addresses[?\(@.type==\"InternalIP\"\)].address})
curl -s $NODE:9962/metrics | grep cilium_nodes_all_num
```
{{% /alert %}}

{{% alert title="Note" color="primary" %}}
It is not yet possible to get metrics from Cilium Envoy (port 9095). Envoy only starts on a node if there is at least one Pod with a layer 7 networkpolicy.
{{% /alert %}}


## {{% task %}} Store and visualize metrics

To make sense of metrics, we store them in Prometheus and visualize them with Grafana dashboards.
Install both into `cilium-monitoring` Namespace to store and visualize Cilium and Hubble metrics.
```bash
kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/v1.12/examples/kubernetes/addons/prometheus/monitoring-example.yaml
```

Make sure Prometheus and Grafana pods are up and running before continuing with the next step.

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


In a second terminal access Grafana with kubectl proxy-forward (for those in the webshell: don't forget to connect to the VM first)
```bash
kubectl -n cilium-monitoring port-forward service/grafana --address ::,0.0.0.0 3000:3000 &
echo "http://$(curl -s ifconfig.me):3000/dashboards"
```

Now open a new tab in your browser and go to URL from the output (for those working on their localmachine use http://localhost:3000/dashboards). In Grafana use the left side menu: `Dashboard`, click on `Manage`, then click on `Hubble`. For a better view, you can change the timespan to the last 5 minutes.

Verify that you see the generated traffic under Network, Forwarded vs Dropped Traffic. Not all graphs will have data available. This is because we have not yet used network policies or any layer 7 components. This will be done in the later chapters.

Change to the  `Cilium Metrics` Dashboard. Here we see information about Cilium itself. Again not all graphs contain data as we have not used all features of Cilium yet.

Try to find the number of IPs allocated and the number of Cilium endpoints.

Leave the Grafana Tab open, we will use it in the later chapters.
