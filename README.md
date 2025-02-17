# slowdrain - A `kubectl` Plugin for Controlled Node Draining

`slowdrain` is a `kubectl` plugin that allows you to drain a node in Kubernetes gradually, ensuring a controlled removal of application pods while preserving infrastructure pods.

## Features

- Drains a Kubernetes node deleting application pods one by one with a delay.
- Randomized pod deletion in different namespaces to avoid service disruptions.
- Configurable delay between pod terminations (default: 20 seconds)
- Optional auto-confirmation mode for scripting and automation

## Who Might Find This Plugin Useful?

`kubectl-slowdrain` is **not** intended to replace the `kubectl drain` command. Instead, it provides an **alternative way** to drain a node by **removing application pods one by one**, introducing a **pause between each deletion**.

### **When to Use This Plugin?**
This approach can be useful in scenarios where:
- You want to **avoid multiple application pods being evicted at the same time**, which could lead to temporary unavailability.
- Your application is **sensitive to sudden pod removals**, and draining too many instances simultaneously could result in **HTTP 500 errors or degraded performance**.
- You have **misconfigured Pod Disruption Budgets (PDBs)** that do not properly regulate pod evictions.

By draining a node **gradually**, `kubectl-slowdrain` provides **a smoother transition** when performing maintenance or handling node failures.

## ⚠️ Warning: No Eviction API Usage

This plugin **does not use the Kubernetes Eviction API**, meaning that **Pod Disruption Budgets (PDBs) will not be respected**.

Pods are deleted **directly**, without considering potential implications such as the number of replicas in the associated Deployment or other safeguards that Kubernetes provides for controlled draining.

Use this tool with caution, especially in production environments, as it may cause unintended service disruptions.

## Installation

### Krew

```sh
kubectl krew install slowdrain
```

### Manual Installation

You can also install the plugin manually by download the latest release from the repository and extracting it to a directory in your PATH.

```shell
# Adjust for your platform
OS_ARCH="linux_amd64"
# Get latest tag
SLOWDRAIN_TAG=$(curl -s https://api.github.com/repos/cturiel/kubectl-slowdrain/releases/latest | grep "tag_name" | sed -E 's/.*"([^"]+)".*/\1/')

# Download and unpack plugin
curl -sL "https://github.com/cturiel/kubectl-slowdrain/releases/download/${SLOWDRAIN_TAG}/kubectl-slowdrain_${SLOWDRAIN_TAG}_${OS_ARCH}.tar.gz" | sudo tar xzvf - -C /usr/local/bin

# Change permission to allow execution
sudo chmod +x /usr/local/bin/kubectl-slowdrain

# Check if plugin is detected
kubectl plugin list
```

## Usage

Run the plugin with:

```sh
kubectl slowdrain <node_name>
```

For detailed usage instructions and options, see [USAGE.md](doc/USAGE.md).

## Demo

[![asciicast](https://asciinema.org/a/703849.svg)](https://asciinema.org/a/703849)

## Contributing

If you'd like to contribute, please check out the repository and submit pull requests or issues.
