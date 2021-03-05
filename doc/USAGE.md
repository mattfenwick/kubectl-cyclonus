
## Usage
The following assumes you have the plugin installed via

```shell
kubectl krew install cyclonus
```

### Lint your kubernetes network policies

```shell
kubectl cyclonus --mode lint -A
```

### Check whether your policies will allow or deny specific traffic

```shell
kubectl cyclonus --mode query-traffic --traffic-path ./examples/traffic.json
```

### Understand connectivity matrix allowed by your network policies

```shell
kubectl cyclonus --mode probe --probe-path ./examples/probe.json
```

### Check which policies apply to a pod

```shell
kubectl cyclonus --mode query-target --target-pod-path ./examples/targets.json
```

### Analyze network policies from files only

```shell
kubectl cyclonus --mode explain --policy-path $PATH
```
