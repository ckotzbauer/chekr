# chekr

A inspection utility for the maintenance of Kubernetes clusters.

[![build](https://github.com/ckotzbauer/chekr/actions/workflows/test.yml/badge.svg)](https://github.com/ckotzbauer/chekr/actions/workflows/test.yml)

Chekr is a cli-tool to handle some common cases of ongoing maintenance tasks when operating Kubernetes clusters.


## Features

* Generate a high-availability report of your workload depending on your configurations.
* Generate a resource-usage overview by querying Prometheus metrics and compare them with the requests and limits beeing set.
* List deprecated API obbjects
* Output as table, json or html.
* Bash completions (Bash, Zsh, Fish, PowerShell)


## Compatibility

* Operating systems:
  * Linux (amd64, arm64, i386)
  * Darwin (amd64, arm64)
  * Windows (amd64, i386)
* Kubernetes: 1.20, 1.21, 1.22


## Installation

### Download binary

Download the latest binary for your platform from the [release page](https://github.com/ckotzbauer/chekr/releases) and store it on your machine.


### Docker

```
docker pull ckotzbauer/chekr:latest
docker pull ghcr.io/ckotzbauer/chekr:latest

docker run --rm \
    -v ~/.kube/config:/home/chekr/.kube/config:ro \
    ckotzbauer/chekr:latest \
    deprecation
```


## Usage

Chekr tries to be as most compatible to `kubectl` as possible. Most of the global flags are directly passed to the kubernetes-client of golang
to give the user full control. By default, the local `kubeconfig` is used or the respective `ServiceAccount` credentials, if a cluster environment
is detected. By default the namespace of the `kubeconfig context` is used for the request.

### Global flags

```
    --as string                      Username to impersonate for the operation
    --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
    --certificate-authority string   Path to a cert file for the certificate authority
    --client-certificate string      Path to a client certificate file for TLS
    --client-key string              Path to a client key file for TLS
    --cluster string                 The name of the kubeconfig cluster to use
    --context string                 The name of the kubeconfig context to use
    --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
-n, --namespace string               If present, the namespace scope for this CLI request
-o, --output string                  Output-Format. Valid values are [table, json, html] (default "table")
    --output-file string             File to write to output to.
    --password string                Password for basic authentication to the API server
    --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
    --server string                  The address and port of the Kubernetes API server
    --tls-server-name string         If provided, this name will be used to validate server certificate. If this is not provided, hostname used to contact the server is used.
    --token string                   Bearer token for authentication to the API server
    --user string                    The name of the kubeconfig user to use
    --username string                Username for basic authentication to the API server
```

### High-availability

This feature generates a report of all pods about their resiliency according their configuration.
To change the default namespace use global `-n` flag. You can filter by pod-selectors with `-l`. All Pods are categorized by multiple factors:
* Type (Deployment, Statefulset, DaemonSet, Job, CRD, standalone)
* Replica count
* Deployment- / Updatestrategy
* PVCs (ROX, RWO, RWX)
* Pod-Anti-Affinity

The pods are ranked with the following categories:
* **0**: undefinied ranking
* **1**: high-available (failure-resilient, zero-downtime-deployment capable)
* **2**: zero-downtime-deployment capable (non failure-resilient)
* **3**: single-point-of-failure
* **4**: standalone pod

**Note:** This only checks common settins which are usually responsible for failure- and deployment-behaviors in a default Kubernetes environment. This
cannot detect special Kubernetes configurations/addons or the application behavior itself in the pod!

```
Creates high-availability report of your workload.

Usage:
  chekr ha [flags]

Flags:
  -h, --help              help for ha
  -l, --selector string   Label-Selector

```

### Resource usage

To generate a report about the resource-consumption from pods you can use this subcommand. It queries the given Prometheus server
(`--prometheus-url` is mandatory) and compares the `Requests`, `Limits` and `Usage` metrics for memory and cpu of the last **30 days**.
To change the default namespace use global `-n` flag. You can filter by pod-selectors with `-l`. The `--timeout` duration applies to the
Prometheus client.

[See used prometheus metrics](https://github.com/ckotzbauer/chekr/blob/master/pkg/resources/metrics.go)

```
Analyze resource requests and limits of pods.

Usage:
  chekr resources [flags]

Flags:
  -h, --help                         help for resources
  -P, --prometheus-password string   Prometheus-Password
  -u, --prometheus-url string        Prometheus-URL (mandatory)
  -U, --prometheus-username string   Prometheus-Username
  -l, --selector string              Label-Selector
  -t, --timeout duration             Timeout (default 30s)
```

### Deprecated API objects

To get an overview of api-objects in your cluster which are deprecated you can use this feature. It scans all objects and gives you a list
of objects which are deprecated and will be removed in a future version. You can ignore kinds with `-i`. To only view deprecations until a given
version you can specify `-V`. This will hide all items, which are deprecated in a never version than specified. Increase the burst with `-t` to bypass throttling
from the Kubernetes server.

**Note:**: This command always scans all namespaces and cannot be filtered with the global `-n` flag.

```
List deprected objects in your cluster.

Usage:
  chekr deprecation [flags]

Flags:
  -h, --help                    help for deprecation
  -i, --ignored-kinds strings   All kinds you want to ignore (e.g. Deployment,DaemonSet)
  -V, --k8s-version string      Highest K8s major.minor version to show deprecations for (e.g. 1.21)
  -t, --throttle-burst int      Burst used for throttling of Kubernetes discovery-client (default 100)
```


[Contributing](https://github.com/ckotzbauer/chekr/blob/master/CONTRIBUTING.md)
--------
[License](https://github.com/ckotzbauer/chekr/blob/master/LICENSE)
--------
[Changelog](https://github.com/ckotzbauer/chekr/blob/master/CHANGELOG.md)
--------