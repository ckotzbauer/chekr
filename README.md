# chekr

A inspection utility for the maintenance of Kubernetes clusters.

[![build](https://github.com/ckotzbauer/chekr/actions/workflows/test.yml/badge.svg)](https://github.com/ckotzbauer/chekr/actions/workflows/test.yml)

Chekr is a cli-tool to handle some common cases of ongoing maintenance tasks when operating Kubernetes clusters.


## Features

* Generate a high-availability report of your workload depending on your configurations.
* Generate a resource-usage overview by querying Prometheus metrics and compare them with the requests and limits beeing set.
* List deprecated API objects
* Generate Kyverno policies for API-Deprecations
* Output as table, json or html.
* Bash completions (Bash, Zsh, Fish, PowerShell)


## Compatibility

* Operating systems:
  * Linux (amd64, arm64, i386)
  * Darwin (amd64, arm64)
  * Windows (amd64, i386)
* Kubernetes: 1.21, 1.22, 1.23


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
  -h, --help                           help for chekr
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
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
  -v, --verbosity string               Log-level (debug, info, warn, error, fatal, panic) (default "warning")
```

### High-availability

This feature generates a report of all pods about their resiliency according their configuration.
To change the default namespace use global `-n` flag. You can filter by label-pod-selectors with `-l` or by annotation key-value-pairs (splitted by a comma) with `-a`. 
All Pods are categorized by multiple factors:
* Type (Deployment, Statefulset, DaemonSet, Job, CRD, standalone)
* Replica count
* Deploymentstrategy / Updatestrategy
* PVCs (ROX, RWO, RWX)
* Pod-Anti-Affinity

The pods are ranked with the following categories:
* **0**: undefinied ranking
* **1**: high-available (failure-resilient, zero-downtime-deployment capable)
* **2**: zero-downtime-deployment capable (non failure-resilient)
* **3**: single-point-of-failure
* **4**: standalone pod

**Note:** This only checks common settings which are usually responsible for failure- and deployment-behaviors in a default Kubernetes environment. This
cannot detect special Kubernetes configurations/addons or the application behavior itself in the pod!

```
Creates high-availability report of your workload.

Usage:
  chekr ha [flags]

Flags:
  -a, --annotation string   Annotation-Selector
  -h, --help                help for ha
  -l, --selector string     Label-Selector
```

### Resource usage

To generate a report about the resource-consumption from pods you can use this subcommand. It queries the given Prometheus server
(`--prometheus-url` is mandatory) and compares the `Usage` metrics for memory and cpu of the last **30 days** with the the resource-requests and -limits.
To change the default namespace use global `-n` flag and for the count of days `-d`. You can filter by label-pod-selectors with `-l` or by annotation key-value-pairs (splitted by a comma) with `-a`. To customize the queried metrics for cpu or memory, use the `--cpu-metric` and `--memory-metric` flags.
The `--timeout` duration applies to the Prometheus client. Stopped Pods are ignored automatically.

[See used prometheus metrics](https://github.com/ckotzbauer/chekr/blob/main/pkg/resources/metrics.go)

**Note:** You can specify a URL to any Promtheus-API-compliant application. Mostly *[Thanos](https://thanos.io/)* is notable here.
You can either specify a HTTP(S) web-address or the name of the pod, its namespace and optional its port in the form `namespace/pod[:port]`. This will do
a port-forward under the hood to fetch the data. When no port is specified, `9090` is used (the default Prometheus port).

```
Analyze resource requests and limits of pods.

Usage:
  chekr resources [flags]

Flags:
  -a, --annotation string            Annotation-Selector
  -d, --count-days int               Count of days to analyze metrics from (until now). (default 30)
      --cpu-metric string            CPU-Usage metric to query (default "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate")
  -h, --help                         help for resources
      --memory-metric string         Memory-Usage metric to query (default "container_memory_working_set_bytes")
  -P, --prometheus-password string   Prometheus-Password
  -u, --prometheus-url string        Prometheus-URL
  -U, --prometheus-username string   Prometheus-Username
  -l, --selector string              Label-Selector
  -t, --timeout duration             Timeout (default 30s)
```

### Deprecated API objects

To stay up-to-date with API deprecations you can either list all deprecated objects which are currently present in your cluster or you can
generate a [Kyverno](https://kyverno.io/) policy for validation. The following flags are present on both subcommands: You can ignore kinds 
with `-i`. To only view deprecations until a given version you can specify `-V`. By default chekr will use the server-version of your cluster. 
This will hide all items, which are deprecated in a newer version than specified.

#### List deprecated API objects

The `list` subcommand scans all objects and gives you a list of then which are deprecated and will be removed in a future version. Increase 
the burst with `-t` to bypass throttling from the Kubernetes server. If deprecations were found, the command will exit with code `1` 
unless you specified the `--omit-exit-code` flag.

**Note:** This command always scans all namespaces and cannot be filtered with the global `-n` flag.
This feature was inspired by https://github.com/rikatz/kubepug.

```
List deprected objects in your cluster.

Usage:
  chekr deprecation list [flags]

Flags:
  -h, --help                    help for deprecation
  -i, --ignored-kinds strings   All kinds you want to ignore (e.g. Deployment,DaemonSet)
  -V, --k8s-version string      Highest K8s major.minor version to show deprecations for (e.g. 1.21)
      --omit-exit-code          Omits the non-zero exit code if deprecations were found.
  -t, --throttle-burst int      Burst used for throttling of Kubernetes discovery-client (default 100)
```

#### Generate Kyverno policy

When you have Kyverno installed in your cluster, you can generate a `ClusterPolicy` to find deprecations instead. You can specify several options
like `--background`, `--category`, `--subject` or `--validation-failure-action` to customize policy settings. With `--dry-run` in conjunction with 
`-o (json|yaml)` you can just generate the policy without applying it to the cluster.

**Note**: Chekr does not validate, that the Kyverno-CRDs are installed in your cluster.

```
Creates Kyverno validation policies for deprecated objects in cluster.

Usage:
  chekr deprecation kyverno-create [flags]

Flags:
      --background                         Whether background scans should be performed. (default true)
      --category string                    Category set for 'policies.kyverno.io/category' annotation. (default "Best Practices")
      --dry-run                            Whether or not the generated policy should be applied.
  -h, --help                               help for kyverno-create
  -i, --ignored-kinds strings              All kinds you want to ignore (e.g. Deployment,DaemonSet)
  -V, --k8s-version string                 Highest K8s major.minor version to show deprecations for (e.g. 1.21)
      --subject string                     Subject set for 'policies.kyverno.io/subject' annotation. (default "Kubernetes APIs")
      --validation-failure-action string   Validation-Failure-Action of the policy (audit or failure). (default "audit")
```

#### Delete Kyverno policy

To delete a previously applyed policy from chekr, simply execute the `kyverno-delete` subcommand.

```
Deletes Kyverno validation policies for deprecated objects in cluster.

Usage:
  chekr deprecation kyverno-delete [flags]

Flags:
  -h, --help   help for kyverno-delete
```


[Contributing](https://github.com/ckotzbauer/chekr/blob/main/CONTRIBUTING.md)
--------
[License](https://github.com/ckotzbauer/chekr/blob/main/LICENSE)
--------
[Changelog](https://github.com/ckotzbauer/chekr/blob/main/CHANGELOG.md)
--------