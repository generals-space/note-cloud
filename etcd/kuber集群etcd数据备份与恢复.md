# kuber集群etcd数据备份与恢复

参考文章

1. [Kubernetes探秘—etcd状态数据及其备份](https://my.oschina.net/u/2306127/blog/2979019)
    - `etcdctl snapshot`子命令进行备份和恢复
    - 使用kuber的`Cronjob`实现定期自动化备份

```
snapshot save           Stores an etcd node backend snapshot to a given file
snapshot restore        Restores an etcd member snapshot to an etcd directory
```

1. 可以直接备份/etc/kubernetes/pki/etcd和/var/lib/etcd下的文件内容。
    - 对于多节点的etcd服务，不能使用直接备份和恢复目录文件的方法。
    - 备份之前先使用docker stop停止相应的服务，然后再启动即可(如果停止etcd服务，备份过程中服务会中断).
    - 缺省配置情况下，每隔10000次改变，etcd将会产生一个snap(如果只备份/var/lib/etcd/member/snap下的文件，不需要停止服务).
2. 通过etcd-client客户端备份。如下(注意，snapshot是在API3里支持的，cert/key/cacert 三个参数名称与API2的命令不同)：

```
sudo ETCDCTL_API=3 etcdctl snapshot save /home/supermap/k8s-backup/data/etcd-snapshot/$(date +%Y%m%d_%H%M%S)_snapshot.db
```
