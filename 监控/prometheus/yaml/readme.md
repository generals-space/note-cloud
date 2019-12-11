在创建deploy之前, 需要先创建secret资源, 不然deploy会一直pending.

在独立的etcd集群部署中, 下面这条命令可以直接执行.

```
kubectl -n monitoring create secret generic etcd-certs --from-file=/etc/etcd/ssl/ca.pem --from-file=/etc/etcd/ssl/etcd.pem --from-file=/etc/etcd/ssl/etcd-key.pem
```

修改`/etc/kubernetes/manifests/`目录下`kube-controller-manager.yaml`与`kube-scheduler.yaml`, 将其中的`--address=127.0.0.1`修改为`--address=0.0.0.0`, 然后重启kubelet服务.

注意: 所有master节点止的`controller manager`和`scheduler`的配置都要修改.

然后还要为ta们创建service资源, 以便prometheus能够访问到.
