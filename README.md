# Cilium Basics Training

Cilium Basics Training Description


## Content Sections

The training content resides within the [content](content) directory.

The main part are the labs, which can be found at [content/en/docs](content/en/docs).


## Hugo

This site is built using the static page generator [Hugo](https://gohugo.io/).

The page uses the [docsy theme](https://github.com/google/docsy) which is included as a Git Submodule.
Docsy is being enhanced using [docsy-plus](https://github.com/puzzle/docsy-plus/) as well as
[docsy-acend](https://github.com/puzzle/docsy-acend/) and [docsy-puzzle](https://github.com/puzzle/docsy-puzzle/)
for brand specific settings.

After cloning the main repo, you need to initialize the submodule like this:

```bash
git submodule update --init --recursive
```

The default configuration uses the acend setup from [config/_default](config/_default/config.toml).
Alternatively you can use the Puzzle setup from [config/puzzle](config/puzzle/config.toml), which is enabled with
`--environment puzzle`.


### Docsy theme usage

* [Official docsy documentation](https://www.docsy.dev/docs/)
* [Docsy Plus](https://github.com/puzzle/docsy-plus/)


### Update submodules for theme updates

Run the following command to update all submodules with their newest upstream version:

```bash
git submodule update --remote
```


### Shortcode usage


#### `onlyWhen` and `onlyWhenNot`

The `onlyWhen` and `onlyWhenNot` shortcodes allow text to be rendered if certain conditions apply.

* `{{% onlyWhen variant1 %}}`: This is only rendered when `enabledModule` in `config.toml` contains `variant1`
* `{{% onlyWhen variant1 variant2 %}}`: This is only rendered when `enabledModule` in `config.toml` contains `variant1` **or** `variant2`
* `{{% onlyWhenNot variant1 %}}`: This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1`
* `{{% onlyWhenNot variant1 variant2 %}}`: This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1` **or** `variant2`

In order to only render text if **all** of multiple conditions do not apply simply chain several `onlyWhenNot` shortcodes:

```
{{% onlyWhenNot variant1 %}}
{{% onlyWhenNot variant2 %}}
This is only rendered when `enabledModule` in `config.toml` **does not** contain `variant1` **nor** `variant2`.
{{% /onlyWhen %}}
{{% /onlyWhen %}}
```


## Build using Docker

Build the image:

```bash
docker build <--build-arg TRAINING_HUGO_ENV=...> -t quay.io/acend/cilium-basics-training .
```

Run it locally:

```bash
docker run -i -p 8080:8080 quay.io/acend/cilium-basics-training
```


## How to develop locally

To develop locally we don't want to rebuild the entire container image every time something changed, and it is also important to use the same hugo versions like in production.
We simply mount the working directory into a running container, where hugo is started in the server mode.

```bash
export HUGO_VERSION=$(grep "FROM klakegg/hugo" Dockerfile | sed 's/FROM klakegg\/hugo://g' | sed 's/ AS builder//g')
docker run --rm --interactive --publish 8080:8080 -v $(pwd):/src klakegg/hugo:${HUGO_VERSION} server -p 8080 --bind 0.0.0.0
```

use the following command to set the hugo environment

```bash
export HUGO_VERSION=$(grep "FROM klakegg/hugo" Dockerfile | sed 's/FROM klakegg\/hugo://g' | sed 's/ AS builder//g')
docker run --rm --interactive --publish 8080:8080 -v $(pwd):/src klakegg/hugo:${HUGO_VERSION} server --environment=<environment> -p 8080 --bind 0.0.0.0
```


## Linting of Markdown content

Markdown files are linted with <https://github.com/DavidAnson/markdownlint>.
Custom rules are in `.markdownlint.json`.
There's a GitHub Action `.github/workflows/markdownlint.yaml` for CI.
For local checks, you can either use Visual Studio Code with the corresponding extension, or the command line like this:

```shell script
npm install
npm run mdlint
```

Npm not installed? no problem

```bash
export HUGO_VERSION=$(grep "FROM klakegg/hugo" Dockerfile | sed 's/FROM klakegg\/hugo://g' | sed 's/ AS builder//g')
docker run --rm --interactive -v $(pwd):/src klakegg/hugo:${HUGO_VERSION}-ci /bin/bash -c "set -euo pipefail;npm install; npm run mdlint;"
```


## How to setup an entire new Training

* create an empty git repo
* Copy the contents of this repo to it
  * check git submodules
    * all of [.gitmodules](.gitmodules) needed?
    * if checkout is not working, add them manually:
      * `git submodule add https://github.com/google/docsy.git ./themes/docsy`
      * `git submodule add https://github.com/puzzle/docsy-plus.git ./themes/docsy-plus`
      * ...
  * replace all CHANGEME
    * `https://github.com/changeme/changeme-training` to your repo url
    * `quay.io/acend/hugo-training-template` to your image registry url
    * `acend/changeme-training` to your org and training
    * `changeme/changeme-training` to your org and training
    * `acend-hugo-training-template-prod` to your prod deployment namespace
    * `acend-hugo-training-template-test` to your test deployment namespace
    * `hugo-training-template` to your training
    * `changeme-training` to your training
    * `changeme Training` to your training name, eg. `Hugo Training`
    * `acend-hugo-template` to your org and training
    * check remaining `changeme`'s
  * Configure all names, URLs and so on in the [build actions](.github/workflows/) and [values.yaml](./helm-chart/values.yaml)
  * remove `How to setup an entire new Training` chapter from README.md
  * adapt or remove not needed variants in the config folder
* Create a container image Repo and make sure the secrets configured in the Github actions have access to the repo
* Create two namespaces on your k8s cluster, make sure the secrets configured in the Github actions have access to the k8s Cluster and namespace or project in case of rancher
  * Test namespace: used to deploy PR Environments
  * Prod namespace: prod deployment


### Quota on Testnamespace

Add the quota to the test namespace:

```bash
kubectl apply -f object-count-quota.yaml -n <namespace>
```


## Github Actions


### Build

The [build action](.github/workflows/build.yaml) is fired on Pull Requests does the following

* builds all PR Versions (Linting and Docker build)
* deploys the built container images to the container registry
* Deploys a PR environment in a k8s test namespace with helm
* Triggers a redeployment
* Comments in the PR where the PR Environments can be found


### PR Cleanup

The [pr-cleanup action](.github/workflows/pr-cleanup.yaml) is fired when Pull Requests are closed and does the following

* Uninstalls PR Helm Release


### Push Main

The [push main action](.github/workflows/push-main.yaml) is fired when a commit is pushed to the main branch (eg. a PR is merged) and does the following, it's very similar to the Build Action

* builds main Versions (Linting and Docker build)
* deploys the built container images to the container registry
* Deploys the main Version on k8s using helm
* Triggers a redeployment


## Helm

Manually deploy the training Release using the following command:

```bash
helm install --repo https://acend.github.io/helm-charts/  <release> acend-training-chart --values helm-chart/values.yaml -n <namespace>
```

For debugging purposes use the `--dry-run` parameter

```bash
helm install --dry-run --repo https://acend.github.io/helm-charts/  <release> acend-training-chart --values helm-chart/values.yaml -n <namespace>
```


## Contributions

If you find errors, bugs or missing information please help us improve and have a look at the [Contribution Guide](CONTRIBUTING.md).
