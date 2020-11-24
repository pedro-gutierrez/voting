# Okteto playground

Usage:

1. Deploy the voting app in your Okteto cloud cluster (mine is at https://vote-pedro-gutierrez.cloud.okteto.net)
2. Bring up the development container, eg:

```
$ okteto up
 ✓  Development container activated
 ✓  Connected to your development container
 ✓  Files synchronized
    Namespace: pedro-gutierrez
    Name:      vote
    Forward:   2345 -> 2345
               8080 -> 80
```

3. Get all pods resource requirements in your namespace, eg:

```
pedro-gutierrez:vote app> go run main.go -namespace=pedro-gutierrez
{"Items":[{"Pod":"db-7c894b85ff-8gjr9","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"redis-86b47f8678-mnzx4","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"result-75c97ddf5d-sjmzl","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"vote-795fcddb76-sg99r","Requirements":{"limits":{"cpu":"1500m","memory":"3Gi"},"requests":{"cpu":"100m","memory":"214748364800m"}}},{"Pod":"worker-5764d777cd-phpgz","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}}]}
```

Note: the first run will be a bit slower as Go dependencies are fetched.

### Outside Kubernetes

If you prefer to run this program outside a Kubernetes cluster, you need to specify the path to your kubeconfig (which you can download from your Okteto Cloud admin console):

```
$ go run main.go -kubeconfig=./kube.config -namespace=pedro-gutierrez

{"Items":[{"Pod":"db-7c894b85ff-8gjr9","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"redis-86b47f8678-mnzx4","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"result-75c97ddf5d-sjmzl","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"vote-598c68df78-6w98m","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"vote-598c68df78-lwdhv","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}},{"Pod":"worker-5764d777cd-phpgz","Requirements":{"limits":{"cpu":"500m","memory":"1Gi"},"requests":{"cpu":"50m","memory":"107374182400m"}}}]}
```