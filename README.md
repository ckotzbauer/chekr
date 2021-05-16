# chekr

A inspection utility for the maintenance of Kubernetes clusters.

[![build](https://github.com/ckotzbauer/chekr/actions/workflows/test.yml/badge.svg)](https://github.com/ckotzbauer/chekr/actions/workflows/test.yml)

Chekr is a cli-tool to handle some common cases of ongoing maintenance tasks when operating Kubernetes clusters.


## Features

* Generate a high-availability report of your workload depending on your configurations.
* Generate a resource-usage overview by querying Prometheus metrics and compare them with the requests and limits beeing set.
* Output as table, json or html.
* Bash completions (Bash, Zsh, Fish, PowerShell)
* ... see the Roadmap


## Compatibility

* Operating systems:
  * Linux (amd64, arm64, i386)
  * Darwin (amd64, arm64)
  * Windows (amd64, i386)
* Kubernetes: 1.20, 1.21, 1.22


## Installation

### Download binary

Download the binary for your platform from the [release page](https://github.com/ckotzbauer/chekr/releases) and store it on your machine.


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

```
Creates high-availability report of your workload.

Usage:
  chekr ha [flags]

Flags:
  -h, --help              help for ha
  -l, --selector string   Label-Selector

```

### Resource usage

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

## Roadmap

* List the usage of deprecated api objects.


[Contributing](https://github.com/ckotzbauer/chekr/blob/master/CONTRIBUTING.md)
--------
[License](https://github.com/ckotzbauer/chekr/blob/master/LICENSE)
--------
[Changelog](https://github.com/ckotzbauer/chekr/blob/master/CHANGELOG.md)
--------