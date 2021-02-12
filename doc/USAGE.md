
## Usage
The following assumes you have the plugin installed via

```shell
kubectl krew install cyclonus
```

### Lint your network policies in kubernetes

```shell
kubectl cyclonus --explain=false --lint=true --all-namespaces
```

### Analyze network policies from files

```shell
kubectl cyclonus --policy-path $PATH
```
