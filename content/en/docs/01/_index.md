---
title: "1. Introduction"
weight: 1
sectionnumber: 1
---

## Title 1

{{% alert title="Note" color="primary" %}}
Sample Note
{{% /alert %}}

Sample code block:
```bash
echo "Hello World!"
```

{{% onlyWhen variant1 %}}
This is only rendered when `enabledModule` in `config.toml` contains `variant1`.
{{% /onlyWhen %}}

{{% onlyWhen variant2 %}}
This is only rendered when `enabledModule` in `config.toml` contains `variant2`.
{{% /onlyWhen %}}

{{% onlyWhen variant1 variant2 %}}
This is only rendered when `enabledModule` in `config.toml` contains `variant1` or `variant2`.
{{% /onlyWhen %}}

{{% onlyWhen variant9 %}}
This is only rendered when `enabledModule` in `config.toml` contains `variant9`.
{{% /onlyWhen %}}

{{% onlyWhenNot variant1 %}}
This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1`.
{{% /onlyWhen %}}

{{% onlyWhenNot variant2 %}}
This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant2`.
{{% /onlyWhen %}}

{{% onlyWhenNot variant1 variant2 %}}
This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1` **nor** `variant2`.
{{% /onlyWhen %}}

{{% onlyWhenNot variant9 %}}
This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant9`.
{{% /onlyWhen %}}


## Title 2


```yaml
foo: bar
```


## Task 1.1: Fix Deployment


```yaml
foo: bar
```


## Task 1.2: Fix Release


```yaml
foo: bar
```


## Task 1.3: Fix Release again


```yaml
foo: bar
```


## Task 1.4: Fix Release again and again


```yaml
foo: bart
```
