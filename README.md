# kubectl cyclonus

A `kubectl` plugin to work with network policies.  This plugin is based on [the cyclonus tool](https://github.com/mattfenwick/cyclonus), which provides the underlying functionality.


## Quick Start

Make sure you have krew installed; then:

```
kubectl krew install cyclonus

kubectl cyclonus -h
```

## Install a few network policies in your cluster

If you don't have any network policies installed on your cluster, here's a few examples to get started:

```shell
$ kubectl create -f networkpolicies/simple-example

$ kubectl get netpol -A
  NAMESPACE   NAME                        POD-SELECTOR   AGE
  y           allow-all-egress-by-label   pod in (a,b)   96s
  y           allow-all-for-label         pod=b          96s
  y           allow-by-ip                 pod=c          96s
  y           allow-label-to-label        pod=a          96s
  y           deny-all                    <none>         96s
  y           deny-all-egress             <none>         96s
  y           deny-all-for-label          pod=a          96s
```

## Explain policies

Groups policies by target, divides rules into egress and ingress, and gives a basic explanation of the combined
policies.  This clarifies the interactions between "denies" and "allows" from multiple policies.

```
$ kubectl cyclonus --mode explain -A

+---------+---------------+------------------------+---------------------+--------------------------+
|  TYPE   |    TARGET     |      SOURCE RULES      |        PEER         |      PORT/PROTOCOL       |
+---------+---------------+------------------------+---------------------+--------------------------+
| Ingress | namespace: y  | y/allow-label-to-label | no ips              | no ports, no protocols   |
|         | Match labels: | y/deny-all-for-label   |                     |                          |
|         |   pod: a      |                        |                     |                          |
+         +               +                        +---------------------+--------------------------+
|         |               |                        | namespace: y        | all ports, all protocols |
|         |               |                        | pods: Match labels: |                          |
|         |               |                        |   pod: c            |                          |
+         +---------------+------------------------+---------------------+                          +
|         | namespace: y  | y/allow-all-for-label  | all pods, all ips   |                          |
|         | Match labels: |                        |                     |                          |
|         |   pod: b      |                        |                     |                          |
+         +---------------+------------------------+---------------------+--------------------------+
|         | namespace: y  | y/allow-by-ip          | ports for all IPs   | no ports, no protocols   |
|         | Match labels: |                        |                     |                          |
|         |   pod: c      |                        |                     |                          |
+         +               +                        +---------------------+--------------------------+
|         |               |                        | 0.0.0.0/24          | all ports, all protocols |
|         |               |                        | except []           |                          |
|         |               |                        |                     |                          |
+         +               +                        +---------------------+--------------------------+
|         |               |                        | no pods             | no ports, no protocols   |
|         |               |                        |                     |                          |
|         |               |                        |                     |                          |
+         +---------------+------------------------+---------------------+                          +
|         | namespace: y  | y/deny-all             | no pods, no ips     |                          |
|         | all pods      |                        |                     |                          |
+---------+---------------+------------------------+---------------------+--------------------------+
```

## Which policy rules apply to a pod?

This takes the previous command a step further: it combines the rules from all the targets that apply
to a pod. 

```
$ kubectl cyclonus \
  --mode query-target \
  --target-pod-path ./examples/targets.json

Combined rules for pod {Namespace:y Labels:map[pod:a]}:
+---------+---------------+-----------------------------+---------------------+--------------------------+
|  TYPE   |    TARGET     |        SOURCE RULES         |        PEER         |      PORT/PROTOCOL       |
+---------+---------------+-----------------------------+---------------------+--------------------------+
| Ingress | namespace: y  | y/allow-label-to-label      | no ips              | no ports, no protocols   |
|         | Match labels: | y/deny-all-for-label        |                     |                          |
|         |   pod: a      | y/deny-all                  |                     |                          |
+         +               +                             +---------------------+--------------------------+
|         |               |                             | namespace: y        | all ports, all protocols |
|         |               |                             | pods: Match labels: |                          |
|         |               |                             |   pod: c            |                          |
+---------+---------------+-----------------------------+---------------------+--------------------------+
|         |               |                             |                     |                          |
+---------+---------------+-----------------------------+---------------------+--------------------------+
| Egress  | namespace: y  | y/deny-all-egress           | all pods, all ips   | all ports, all protocols |
|         | Match labels: | y/allow-all-egress-by-label |                     |                          |
|         |   pod: a      |                             |                     |                          |
+---------+---------------+-----------------------------+---------------------+--------------------------+
```


## Will policies allow or block traffic?

Given arbitrary traffic examples (from a source to a destination, including labels, over a port and protocol),
this command parses network policies and determines if the traffic is allowed or not.

```
$ kubectl cyclonus \
  --mode query-traffic \
  --traffic-path ./examples/traffic.json

Traffic:
+--------------------------+-------------+---------------+-----------+-----------+------------+
|      PORT/PROTOCOL       | SOURCE/DEST |    POD IP     | NAMESPACE | NS LABELS | POD LABELS |
+--------------------------+-------------+---------------+-----------+-----------+------------+
| 80 (serve-80-tcp) on TCP | source      | 192.168.1.99  | y         | ns: y     | app: c     |
+                          +-------------+---------------+           +           +------------+
|                          | destination | 192.168.1.100 |           |           | pod: b     |
+--------------------------+-------------+---------------+-----------+-----------+------------+

Is traffic allowed?
+-------------+--------+---------------+
|    TYPE     | ACTION |    TARGET     |
+-------------+--------+---------------+
| Ingress     | Allow  | namespace: y  |
|             |        | Match labels: |
|             |        |   pod: b      |
+             +--------+---------------+
|             | Deny   | namespace: y  |
|             |        | all pods      |
+-------------+--------+---------------+
|             |        |               |
+-------------+--------+---------------+
| Egress      | Deny   | namespace: y  |
|             |        | all pods      |
+-------------+--------+---------------+
| IS ALLOWED? | FALSE  |                
+-------------+--------+---------------+
```

## Simulated probe

Runs a simulated connectivity probe against a set of network policies, without using a kubernetes cluster.

```
$ kubectl cyclonus \
  --mode probe \
  --probe-path ./examples/probe.json

Combined:
+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
|     | X/A | X/B | X/C | Y/A | Y/B | Y/C | Z/A | Z/B | Z/C |
+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
| x/a | .   | .   | .   | X   | .   | X   | .   | .   | .   |
| x/b | .   | .   | .   | X   | .   | X   | .   | .   | .   |
| x/c | .   | .   | .   | X   | .   | X   | .   | .   | .   |
| y/a | .   | .   | .   | X   | .   | X   | .   | .   | .   |
| y/b | .   | .   | .   | X   | .   | X   | .   | .   | .   |
| y/c | X   | X   | X   | X   | X   | X   | X   | X   | X   |
| z/a | .   | .   | .   | X   | .   | X   | .   | .   | .   |
| z/b | .   | .   | .   | X   | .   | X   | .   | .   | .   |
| z/c | .   | .   | .   | X   | .   | X   | .   | .   | .   |
+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
```

## Linter

Checks network policies for common problems.

```
$ kubectl cyclonus \
  --mode lint \
  --lint=true

+-----------------+------------------------------+-------------------+-----------------------------+
| SOURCE/RESOLVED |             TYPE             |      TARGET       |       SOURCE POLICIES       |
+-----------------+------------------------------+-------------------+-----------------------------+
| Resolved        | CheckTargetAllEgressAllowed  | namespace: y      | y/allow-all-egress-by-label |
|                 |                              |                   |                             |
|                 |                              | pod selector:     |                             |
|                 |                              | matchExpressions: |                             |
|                 |                              | - key: pod        |                             |
|                 |                              |   operator: In    |                             |
|                 |                              |   values:         |                             |
|                 |                              |   - a             |                             |
|                 |                              |   - b             |                             |
|                 |                              |                   |                             |
+-----------------+------------------------------+-------------------+-----------------------------+
| Resolved        | CheckDNSBlockedOnTCP         | namespace: y      | y/deny-all-egress           |
|                 |                              |                   |                             |
|                 |                              | pod selector:     |                             |
|                 |                              | {}                |                             |
|                 |                              |                   |                             |
+-----------------+------------------------------+-------------------+-----------------------------+
| Resolved        | CheckDNSBlockedOnUDP         | namespace: y      | y/deny-all-egress           |
|                 |                              |                   |                             |
|                 |                              | pod selector:     |                             |
|                 |                              | {}                |                             |
|                 |                              |                   |                             |
+-----------------+------------------------------+-------------------+-----------------------------+
```

# How to release

This uses krew automation; just push a new tag and wait for [the goreleaser action](https://github.com/mattfenwick/kubectl-cyclonus/actions) to run.

# How to test a local krew release

1. Choose one of the following:

     - easy way: `kubectl krew install --manifest=./deploy/krew/cyclonus.yaml`

     - hard way:
        1. download a binary from [the github project release page](https://github.com/mattfenwick/kubectl-cyclonus/releases)

        2. run a `krew install` against the downloaded binary

            ```
            kubectl krew install --manifest=./deploy/krew/cyclonus.yaml --archive=/Users/mfenwick/Downloads/kubectl-cyclonus_darwin_amd64.tar.gz
            ```

2. test it

    ```
    kubectl cyclonus -h
    ```

3. clean up

    ```
    kubectl krew uninstall cyclonus
    ```
