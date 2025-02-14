# slowdrain - A `kubectl` Plugin for Controlled Node Draining

`slowdrain` is a `kubectl` plugin that allows you to drain a node in Kubernetes gradually, ensuring a controlled removal of application pods while preserving infrastructure pods.

## Features

- Drains a Kubernetes node deleting application pods one by one with a delay.
- Randomized pod deletion in different namespaces to avoid service disruptions.
- Configurable delay between pod terminations (default: 20 seconds)
- Optional auto-confirmation mode for scripting and automation

## ⚠️ Warning: No Eviction API Usage

This plugin **does not use the Kubernetes Eviction API**, meaning that **Pod Disruption Budgets (PDBs) will not be respected**.

Pods are deleted **directly**, without considering potential implications such as the number of replicas in the associated Deployment or other safeguards that Kubernetes provides for controlled draining.

Use this tool with caution, especially in production environments, as it may cause unintended service disruptions.

## Installation

```sh
kubectl krew install slowdrain
```

## Usage

Run the plugin with:

```sh
kubectl slowdrain NODE_NAME
```

For detailed usage instructions and options, see [USAGE.md](doc/USAGE.md).

## Contributing

If you'd like to contribute, please check out the repository and submit pull requests or issues.
