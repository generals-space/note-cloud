# 内置etcd的etcdctl操作示例(如何指定证书)

```
etcdctl --endpoints=https://192.168.7.13:2379 --cert-file=/etc/kubernetes/pki/etcd/server.crt --key-file=/etc/kubernetes/pki/etcd/server.key --ca-file=/etc/kubernetes/pki/etcd/ca.crt cluster-health
```

查看目录下的文件列表

```
ETCDCTL_API=3 etcdctl --endpoints=https://192.168.7.13:2379 --cert=/etc/etcd/ssl/etcd.pem --key=/etc/etcd/ssl/etcd-key.pem --cacert=/etc/etcd/ssl/ca.pem get / --prefix --keys-only

ETCDCTL_API=3 etcdctl --endpoints=https://192.168.7.13:2379 --cert=/etc/etcd/ssl/etcd.pem --key=/etc/etcd/ssl/etcd-key.pem --cacert=/etc/etcd/ssl/ca.pem get /registry/pods/kube-system/ --prefix --keys-only
```
