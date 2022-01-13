---
title: "5. Network Policies"
weight: 5
sectionnumber: 5
---

## Network Policies

One of the most basic CNI functions is the ability to enforce network policies and implement an in-cluster zero-trust container strategy. Network policies are a default Kubernetes object for controlling network traffic, but a CNI such as Cilium is required to enforce them. We will demonstrate traffic blocking with our simple app.

{{% alert title="Note" color="primary" %}}
If you are not yet familiar with Kubernetes Network Policies we suggest to go to the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
{{% /alert %}}


## Task {{% param sectionnumber %}}.1: Cilium Endpoint

Each Pod from our simple application is represented in Cilium as an [Endpoint](https://docs.cilium.io/en/stable/concepts/terminology/#endpoint). We can use the `cilium` tool inside a Cilium pod to list them.

First get all Cilium pods with:

```bash
kubectl -n kube-system get pods -l k8s-app=cilium
```

```NAME           READY   STATUS    RESTARTS        AGE
cilium-mvh65   1/1     Running   1 (6h26m ago)   6h30m
```

and then run:

```bash
kubectl -n kube-system exec <podname> -- cilium endpoint list
```

{{% alert title="Note" color="primary" %}}
Or you can also use some jsonpath magic and execute the command in one line:

```bash
kubectl -n kube-system exec $(kubectl -n kube-system get pods -l k8s-app=cilium -o jsonpath='{.items[0].metadata.name}') -- cilium endpoint list
```
{{% /alert %}}


## Task {{% param sectionnumber %}}.2: Verify connectivity

Make your your `FRONTEND` and `NOT_FRONTEND` environment variable are still set. Otherwise set them again:

```bash
FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}')
echo ${FRONTEND}
NOT_FRONTEND=$(kubectl get pods -l app=not-frontend -o jsonpath='{.items[0].metadata.name}')
echo ${NOT_FRONTEND}
```

Now we generate some traffic as a baseline test.

```bash
kubectl exec -ti ${FRONTEND} -- curl -I --connect-timeout 5 backend:8080
kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
```

This will execute a simple `curl` call from the `frontend` and `not-frondend` application to the `backend` application:

```
# Frontend
HTTP/1.1 200 OK
X-Powered-By: Express
Vary: Origin, Accept-Encoding
Access-Control-Allow-Credentials: true
Accept-Ranges: bytes
Cache-Control: public, max-age=0
Last-Modified: Sat, 26 Oct 1985 08:15:00 GMT
ETag: W/"83d-7438674ba0"
Content-Type: text/html; charset=UTF-8
Content-Length: 2109
Date: Tue, 23 Nov 2021 12:50:44 GMT
Connection: keep-alive

# Not Frontend
HTTP/1.1 200 OK
X-Powered-By: Express
Vary: Origin, Accept-Encoding
Access-Control-Allow-Credentials: true
Accept-Ranges: bytes
Cache-Control: public, max-age=0
Last-Modified: Sat, 26 Oct 1985 08:15:00 GMT
ETag: W/"83d-7438674ba0"
Content-Type: text/html; charset=UTF-8
Content-Length: 2109
Date: Tue, 23 Nov 2021 12:50:44 GMT
Connection: keep-alive
```

and we see, both applications can connect to the `backend` application.

Until now ingress and egress policy enforcement is still disabled on all of our pods because no network policy has been imported yet selecting any of the pods. Let us change this.


## Task {{% param sectionnumber %}}.3: Disallow traffic with a Network Policy

We block traffic by applying the following network policy:

{{< highlight yaml >}}{{< readfile file="content/en/docs/05/backend-ingress-deny.yaml" >}}{{< /highlight >}}

The policy will deny all ingress traffic as it is of type Ingress but specifies no allow rule, and will be applied to all pods with the `app=backend` label thanks to the podSelector.

Ok, then let's create the policy with:

```bash
kubectl apply -f backend-ingress-deny.yaml
```

and you can verify the created Network Policy with:

```bash
kubectl get netpol
```

which gives you an output similar to this:

```
                                                    
NAME                   POD-SELECTOR   AGE
backend-ingress-deny   app=backend    2s

```


## Task {{% param sectionnumber %}}.4: Verify connectivity again

We can now execute the connectivity check again:

```bash
kubectl exec -ti ${FRONTEND} -- curl -I --connect-timeout 5 backend:8080
kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
```

but this time you see that the `frontend` and `not-frontend` application cannot connect anymore to the `backend`:

```
# Frontend
curl: (28) Connection timed out after 5001 milliseconds
command terminated with exit code 28
# Not Frontend
curl: (28) Connection timed out after 5001 milliseconds
command terminated with exit code 28
```

The network policy correctly switched the default ingress behavior from default allow to default deny. We can also check this in grafana.

{{% alert title="Note" color="primary" %}}
Note: our earlier grafana port-forward should still be running (can be checked by running jobs or `ps aux | grep "cilium-monitoring"`). If it does not, the website http://localhost:3000/ will not be available.

```bash
kubectl -n cilium-monitoring port-forward service/grafana --address 0.0.0.0 --address :: 3000:3000 &
```

{{% /alert %}}


In grafana browse to the dasboard `Hubble Metrics`. You should see now data in more graphs. Check the graphs `Drop Reason`, `Forwarded vs Dropped` and `Top 10 Source Pods with Denied Packets`. You should find the pods from our simple application there.

Let's now selectively re-allow traffic again, but only from frontend to backend.


## Task {{% param sectionnumber %}}.5: Allow traffic from frontend to backend

We can do it by crafting a new network policy manually, but we can also use the Network Policy Editor to help us out:

* Go to https://networkpolicy.io/editor.
* Upload our initial backend-ingress-deny policy.

![Cilium editor with backend-ingress-deny Policy](cilium_editor_1.png)

* Rename the network policy to backend-allow-ingress-frontend (using the Edit button in the center).

![Cilium editor edit name](cilium_editor_edit_name.png)

* On the ingress side, add `app=frontend` as podSelector for pods in the same namespace.

![Cilium editor add rule](cilium_editor_add.png)

* Inspect the ingress flow colors: the policy will deny all ingress traffic to pods labelled `app=backend`, except for traffic coming from pods labelled `app=frontend`.

![Cilium editor backend allow rule](cilium_editor_backend-allow-ingress.png)


* Download the policy YAML file.

The file should look like this:

{{< highlight yaml >}}{{< readfile file="content/en/docs/05/backend-allow-ingress-frontend.yaml" >}}{{< /highlight >}}

Apply the new policy:

```bash
kubectl create -f backend-allow-ingress-frontend.yaml
```

and then execute the connectivity test again:

```bash
kubectl exec -ti ${FRONTEND} -- curl -I --connect-timeout 5 backend:8080
kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
```

This time, the `frontend` application is able to connect to the `backend` but the `not-frontend` application still cannot connect to the `backend`:

```
# Frontend
HTTP/1.1 200 OK
X-Powered-By: Express
Vary: Origin, Accept-Encoding
Access-Control-Allow-Credentials: true
Accept-Ranges: bytes
Cache-Control: public, max-age=0
Last-Modified: Sat, 26 Oct 1985 08:15:00 GMT
ETag: W/"83d-7438674ba0"
Content-Type: text/html; charset=UTF-8
Content-Length: 2109
Date: Tue, 23 Nov 2021 13:08:27 GMT
Connection: keep-alive

# Not Frontend
curl: (28) Connection timed out after 5001 milliseconds
command terminated with exit code 28

```

Note that this is working despite the fact we did not delete the previous `backend-ingress-deny` policy:

```bash
kubectl get netpol
```

```
NAME                             POD-SELECTOR   AGE
backend-allow-ingress-frontend   app=backend    2m7s
backend-ingress-deny             app=backend    12m

```

Network policies are additive. Just like with firewalls, it is thus a good idea to have default DENY policies and then add more specific ALLOW policies as needed.

We can verify our connection being blockend with hubble.

Generate some traffic.

```bash
kubectl exec -ti ${NOT_FRONTEND} -- curl -I --connect-timeout 5 backend:8080
```

With hubble observe you can now check the packet beeing dropped along with the cause (Policy denied).

{{% alert title="Note" color="primary" %}}
Note: our earlier cilium hubble port-forward should still be running (can be checked by running jobs or `ps aux | grep "cilium hubble port-forward"`). If it does not, hubble status will fail and we have to run it again:

```bash
cilium hubble port-forward&
hubble status
```

{{% /alert %}}

```bash
hubble observe --from-label app=not-frontend
```

And the output should look like this:
```bash
Jan 13 12:54:46.883: default/not-frontend-8f467ccbd-lh4w4:50802 -> kube-system/coredns-64897985d-7rjfv:53 to-endpoint FORWARDED (UDP)
Jan 13 12:54:46.883: default/not-frontend-8f467ccbd-lh4w4:50802 -> kube-system/coredns-64897985d-7rjfv:53 to-endpoint FORWARDED (UDP)
Jan 13 12:54:46.884: default/not-frontend-8f467ccbd-lh4w4:37134 <> default/backend-65f7c794cc-pj2tc:8080 Policy denied DROPPED (TCP Flags: SYN)
Jan 13 12:54:46.884: default/not-frontend-8f467ccbd-lh4w4:37134 <> default/backend-65f7c794cc-pj2tc:8080 Policy denied DROPPED (TCP Flags: SYN)
Jan 13 12:54:47.906: default/not-frontend-8f467ccbd-lh4w4:37134 <> default/backend-65f7c794cc-pj2tc:8080 Policy denied DROPPED (TCP Flags: SYN)
Jan 13 12:54:47.906: default/not-frontend-8f467ccbd-lh4w4:37134 <> default/backend-65f7c794cc-pj2tc:8080 Policy denied DROPPED (TCP Flags: SYN)
Jan 13 12:54:49.922: default/not-frontend-8f467ccbd-lh4w4:37134 <> default/backend-65f7c794cc-pj2tc:8080 Policy denied DROPPED (TCP Flags: SYN)
Jan 13 12:54:49.922: default/not-frontend-8f467ccbd-lh4w4:37134 <> default/backend-65f7c794cc-pj2tc:8080 Policy denied DROPPED (TCP Flags: SYN)
```


## Task {{% param sectionnumber %}}.6: Inspecting the cilium endpoints again

We can now check the cilium endpoints again.

```bash
kubectl -n kube-system exec -it ds/cilium -- cilium endpoint list
```

And now we see that the pods with the label `app=backend` now have ingress policy enforcement enabled.
