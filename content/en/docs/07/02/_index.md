---
title: "7.2 HTTP-aware L7 Policy"
weight: 72
sectionnumber: 7.2
---


## Task {{% param sectionnumber %}}.1: Deploy a new Demo Application

In this Star Wars inspired example, there are three microservices applications: deathstar, tiefighter, and xwing. The deathstar runs an HTTP webservice on port 80, which is exposed as a Kubernetes Service to load balance requests to deathstar across two Pod replicas. The deathstar service provides landing services to the empire’s spaceships so that they can request a landing port. The tiefighter Pod represents a landing-request client service on a typical empire ship and xwing represents a similar service on an alliance ship. They exist so that we can test different security policies for access control to deathstar landing services.

The file `sw-app.yaml` contains a Kubernetes Deployment for each of the three services. Each deployment is identified using the Kubernetes labels (`org=empire`, `class=deathstar`), (`org=empire`, `class=tiefighter`), and (`org=alliance`, `class=xwing`). It also includes a deathstar-service, which load balances traffic to all pods with labels `org=empire` and `class=deathstar`.

{{< highlight yaml >}}{{< readfile file="content/en/docs/07/02/sw-app.yaml" >}}{{< /highlight >}}

Create and apply the file with:

```bash
kubectl apply -f sw-app.yaml
```

And as we have already some Network Policies in our Namespace the default ingress behavior is default deny. Therefore we need a new Network Policy to access services on `deathstar`:

Create a file `cnp.yaml` with the following content:

{{< highlight yaml >}}{{< readfile file="content/en/docs/07/02/cnp.yaml" >}}{{< /highlight >}}

Apply the `CiliumNetworkPolicy` with:

```bash
kubectl apply -f cnp.yaml
```

With this policy, our `tiefighter` has access to the `deathstar` application. You can verify this with:

```bash
kubectl exec tiefighter -- curl -m 2 -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
```

```
Ship landed
```

but the `xwing` does not have access:

```bash
kubectl exec xwing -- curl -m 2 -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
```

```
command terminated with exit code 28
```


## Task {{% param sectionnumber %}}.2: Apply and Test HTTP-aware L7 Policy

In the simple scenario above, it was sufficient to either give tiefighter / xwing full access to deathstar’s API or no access at all. But to provide the strongest security (i.e., enforce least-privilege isolation) between microservices, each service that calls deathstar’s API should be limited to making only the set of HTTP requests it requires for legitimate operation.

For example, consider that the deathstar service exposes some maintenance APIs that should not be called by random empire ships. To see this run:

```bash
kubectl exec tiefighter -- curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
```

```
Panic: deathstar exploded

goroutine 1 [running]:
main.HandleGarbage(0x2080c3f50, 0x2, 0x4, 0x425c0, 0x5, 0xa)
        /code/src/github.com/empire/deathstar/
        temp/main.go:9 +0x64
main.main()
        /code/src/github.com/empire/deathstar/
        temp/main.go:5 +0x85
```

Cilium is capable of enforcing HTTP-layer (i.e., L7) policies to limit what URLs the tiefighter is allowed to reach. Here is an example policy file that extends our original policy by limiting tiefighter to making only a POST /v1/request-landing API call, but disallowing all other calls (including PUT /v1/exhaust-port).

{{< highlight yaml >}}{{< readfile file="content/en/docs/07/02/cnp-l7.yaml" >}}{{< /highlight >}}

Update the existing rule to apply the L7-aware policy to protect deathstar using. Create a file `cnp-l7.yaml` with the above content and apply with:

```bash
kubectl apply -f cnp-l7.yaml
```

We can now re-run the same test as above, but we will see a different outcome:

```bash
kubectl exec tiefighter -- curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
```

```
Ship landed
```

and

```bash
kubectl exec tiefighter -- curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
```

```
Access denied
```

{{% alert title="Note" color="primary" %}}
You can now check the `Hubble` dashboard in Grafana again. The graphs under HTTP should soon show some data as well. To generate more data just request-landing on `deathstar` a few times with `tiefighter`
{{% /alert %}}
