# kubectl run创建临时pod

```console
kubectl run --rm -i --tty fun --image quay.io/coreos/etcd --restart=Never -- /bin/sh
```
