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

{"Items":[{"Pod":"db-57b9cc6d77-gqmqv","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"redis-7d4d7d9d86-4cqq6","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"result-9c45559fd-hdxwf","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"vote-76b878ddf8-dn88r","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"vote-76b878ddf8-mh9cm","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"worker-59bdcb645d-ct7zl","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"worker-687bb7596b-h9cg4","Requirements":{"limits":{"cpu":"500m","memory":"500Mi"},"requests":{"cpu":"50m","memory":"52428800"}}}]}
```

Note: the first run will be a bit slower as Go dependencies are fetched.

### Outside Kubernetes

If you prefer to run this program outside a Kubernetes cluster, you need to specify the path to your kubeconfig (which you can download from your Okteto Cloud admin console):

```
$ go run main.go -kubeconfig=./kube.config -namespace=pedro-gutierrez

{"Items":[{"Pod":"db-57b9cc6d77-gqmqv","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"redis-7d4d7d9d86-4cqq6","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"result-9c45559fd-hdxwf","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"vote-76b878ddf8-dn88r","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"vote-76b878ddf8-mh9cm","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"worker-59bdcb645d-ct7zl","Requirements":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"50m","memory":"20971520"}}},{"Pod":"worker-687bb7596b-h9cg4","Requirements":{"limits":{"cpu":"500m","memory":"500Mi"},"requests":{"cpu":"50m","memory":"52428800"}}}]}
```