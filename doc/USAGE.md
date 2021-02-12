
## Usage
The following assumes you have the plugin installed via

```shell
kubectl krew install cyclonus
```

### Lint your kubernetes network policies

```shell
kubectl cyclonus --explain=false --lint=true
```

### Check whether your policies will allow or deny specific traffic

```shell
kubectl cyclonus --explain=false --traffic-path ./examples/traffic.json
```

### Understand connectivity matrix allowed by your network policies

```shell
kubectl cyclonus --explain=false --probe-path ./examples/probe.json
```

### Check which policies apply to a pod

```shell
kubectl cyclonus --explain=false --target-pod-path ./examples/targets.json
```

### Analyze network policies from files only

```shell
kubectl cyclonus --policy-path $PATH --all-namespaces=false
```
