## Usage

The following assumes you have the plugin installed via:

```shell
kubectl krew install kubectl-slowdrain
```

### Drain a node slowly

To drain a node while removing application pods one by one with a delay:

```shell
kubectl slowdrain <node_name>
```

By default, the plugin will prompt for confirmation before deleting pods.

### Set a custom delay between pod deletions

Use the `-d` or `--delay` flag to specify the delay (in seconds) between each pod deletion:

```shell
kubectl slowdrain <node_name> --delay 60
```

### Auto-confirm pod deletion

Use the `-y` flag to skip the confirmation prompt and automatically proceed:

```shell
kubectl slowdrain <node_name> -y
```

### Set log level

Use `--log-level` to adjust the logging verbosity. Available levels: `info`, `warn`, `error`, `debug`.

```shell
kubectl slowdrain <node_name> --log-level debug
```

## How it works

`kubectl-slowdrain` is a Kubernetes plugin designed to safely drain a node by removing application pods one by one with a configurable delay. It:

1. Cordon the target node to prevent new pods from being scheduled.
2. Lists the pods running on the node, categorizing them into infrastructure (`kube-*`, `openshift-*`, `infra-*`) and application pods.
3. Prompts for confirmation before proceeding unless `-y` is used.
4. Deletes application pods randomly, respecting the configured delay.
5. Logs operations at the selected verbosity level.

This plugin is useful in environments where careful pod eviction is required to maintain service availability and avoid cascading failures.
