# rback

[![CI](https://github.com/tmojzes/rback/actions/workflows/ci.yml/badge.svg)](https://github.com/tmojzes/rback/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tmojzes/rback)](https://goreportcard.com/report/github.com/tmojzes/rback)
[![GitHub Release](https://img.shields.io/github/v/release/tmojzes/rback)](https://github.com/tmojzes/rback/releases)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

A simple "RBAC in Kubernetes" visualizer. No matter how complex the setup, `rback` queries all RBAC related information of a Kubernetes cluster in constant time and generates a graph representation of service accounts, (cluster) roles, and the respective access rules in [Graphviz dot](https://www.graphviz.org/doc/info/lang.html) format.

For example, here is an Amazon EKS cluster as seen by `rback`:

![EKS cluster](examples/eks.dot.png)

Another example would be a local K3S cluster:

![K3S cluster](examples/k3s.dot.png)

See the [examples/](examples/) directory for more samples.

---

## Installation

`rback` requires access to a Kubernetes cluster (e.g., EKS, GKE, AKS, k3s, kind, Minikube) and `kubectl` configured locally.

### Pre-built Binaries

Download the latest pre-compiled binary for your architecture (macOS Intel & Apple Silicon, Linux amd64 & arm64, Windows) from the [Releases page](https://github.com/tmojzes/rback/releases):

```sh
# Example for macOS (Apple Silicon / arm64)
curl -sL https://github.com/tmojzes/rback/releases/latest/download/rback_darwin_arm64.tar.gz | tar -xz
chmod +x rback && sudo mv rback /usr/local/bin

# Example for Linux (x86_64 / amd64)
curl -sL https://github.com/tmojzes/rback/releases/latest/download/rback_linux_amd64.tar.gz | tar -xz
chmod +x rback && sudo mv rback /usr/local/bin
```

### Via `go install`

If you have Go installed:

```sh
go install github.com/tmojzes/rback@latest
```

### Building From Source

Clone the repository and build using [Task](https://taskfile.dev):

```sh
git clone https://github.com/tmojzes/rback.git
cd rback
task build
# or install directly to $GOPATH/bin:
task install
```

> [!NOTE]
> If you do not have `task` installed globally, you can invoke it directly via Go's tool mechanism:
>
> ```sh
> go tool task build
> ```

---

## Using `rback` Directly

Query the cluster using `kubectl` and pipe the JSON output to `rback`:

```sh
kubectl get sa,roles,rolebindings,clusterroles,clusterrolebindings --all-namespaces -o json | rback > result.dot
```

Alternatively, read from a saved file using `-f`:

```sh
rback -f rbac-dump.json > result.dot
```

### Rendering the Graph

#### Render Locally

Install [Graphviz](https://www.graphviz.org/) (`brew install graphviz` on macOS or `sudo apt-get install graphviz` on Linux):

```sh
# On macOS:
kubectl get sa,roles,rolebindings,clusterroles,clusterrolebindings --all-namespaces -o json | rback | dot -Tpng > /tmp/rback.png && open /tmp/rback.png

# On Linux:
kubectl get sa,roles,rolebindings,clusterroles,clusterrolebindings --all-namespaces -o json | rback | dot -Tpng > /tmp/rback.png && xdg-open /tmp/rback.png
```

#### Render Online

You can paste the generated `.dot` content directly into browser-based Graphviz viewers such as:

- [magjac.com/graphviz-visual-editor/](http://magjac.com/graphviz-visual-editor/)
- [dreampuf.github.io/GraphvizOnline](https://dreampuf.github.io/GraphvizOnline/)

---

## Using `rback` as a `kubectl` Plugin

A plugin script is provided in [`kubectl-plugin/kubectl-rback`](kubectl-plugin/kubectl-rback). To install:

1. Copy or symlink `kubectl-plugin/kubectl-rback` into your `PATH` (e.g. `/usr/local/bin/kubectl-rback`).
2. Ensure it is executable (`chmod +x /usr/local/bin/kubectl-rback`).
3. Ensure both `rback` and Graphviz `dot` are in your `PATH`.

Now run:

```sh
kubectl rback
```

This will automatically query the cluster, generate the graph, render a PNG, and open it in your default image viewer (supporting both macOS and Linux).

---

## Usage Examples

### Filtering by Namespace

Focus on a single namespace or comma-separated list of namespaces:

```sh
kubectl rback -n my-namespace
kubectl rback -n my-namespace1,my-namespace2
```

### Focusing on Specific Resources

To focus on a specific resource and its direct relationships:

```sh
# ServiceAccounts:
kubectl rback serviceaccount my-service-account
kubectl rback sa my-service-account

# Roles & ClusterRoles:
kubectl rback role my-role
kubectl rback clusterrole my-cluster-role

# RoleBindings & ClusterRoleBindings:
kubectl rback rolebinding my-role-binding
kubectl rback clusterrolebinding my-cluster-role-binding
```

Short aliases are supported: `sa`, `r`, `cr`, `rb`, `crb`, `u` (user), `g` (group).

You can also specify multiple resource names:

```sh
kubectl rback r my-role1 my-role2
```

### Checking Permissions (`who-can`)

Find out who can perform an action in the cluster:

```sh
kubectl rback who-can create pods
kubectl rback who-can get secrets
kubectl rback who-can list services
```

This renders matched `(Cluster)Roles`, related bindings, and subjects (`ServiceAccounts`, `Users`, `Groups`). The matched rule is highlighted in bold.

To only show the matched rules instead of all rules in the role:

```sh
kubectl rback --show-matched-rules-only who-can create pods
```

### Toggling Rules & Legend

Hide access rules completely:

```sh
kubectl rback --show-rules=false
```

Hide the graph legend:

```sh
kubectl rback --show-legend=false
```

---

## Development

All development tasks are managed via [Task](https://taskfile.dev):

```sh
# List available tasks
task --list

# Build local binary
task build

# Run unit tests with race detection
task test

# Generate test coverage report
task test:coverage

# Run linter / go vet
task lint

# Install binary to $GOPATH/bin
task install

# Clean build artifacts
task clean
```

> [!NOTE]
> If you don't have `task` installed, prepend `go tool` (e.g. `go tool task test`).

---

## Attribution & License

`rback` was originally created by [Michael Hausenblas](https://github.com/mhausenblas) and maintained under `team-soteria/rback`.

This fork is actively maintained by [Tamás Mojzes](https://github.com/tmojzes).

Distributed under the [Apache License 2.0](LICENSE).
