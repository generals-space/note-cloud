# kuber集群etcd数据备份与恢复

参考文章

1. [Operating etcd clusters for Kubernetes](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/)
    - 官方文档
2. [Kubernetes探秘—etcd状态数据及其备份](https://my.oschina.net/u/2306127/blog/2979019)
    - `etcdctl snapshot`子命令进行备份和恢复
    - 使用kube的`Cronjob`实现定期自动化备份
3. [openshift 4.5.9 etcd损坏+脑裂修复过程](https://zhangguanzhang.github.io/2021/06/08/ocp4.5.9-restore-etcd/#/%E6%93%8D%E4%BD%9C)
4. [ETCD集群故障应急恢复-从snapshot恢复](https://blog.csdn.net/weixin_43845924/article/details/124975494)
5. [ETCD 常用操作，及数据备份和恢复](https://blog.csdn.net/zyw_0813/article/details/128912914)
    - 单节点, 及多节点集群的数据恢复方式

## 备份

```
export ETCDCTL_API=3
export ETCDCTL_CACERT=/etc/kubernetes/pki/etcd/ca.crt
export ETCDCTL_CERT=/etc/kubernetes/pki/etcd/server.crt
export ETCDCTL_KEY=/etc/kubernetes/pki/etcd/server.key
etcdctl snapshot save snapshot.db
```

> `snapshot.db`为生成的备份文件, 名称和路径随意.

## 恢复

**恢复操作需要先将 etcd 服务停掉, 如果是集群, 则集群所有节点都停掉**.

### 单节点

如果etcd是单节点, 可直接执行如下命令.

```
export ETCDCTL_API=3
export ETCDCTL_CACERT=/etc/kubernetes/pki/etcd/ca.crt
export ETCDCTL_CERT=/etc/kubernetes/pki/etcd/server.crt
export ETCDCTL_KEY=/etc/kubernetes/pki/etcd/server.key
etcdctl snapshot restore snapshot.db --data-dir <data-dir-location>
```

官网并没有解释`--data-dir`的涵义, 最初以为是备份文件的路径.

其实`restore`类似一个解压命令, 把`snapshot.db`文件解压成一个数据目录, 我们可以将当前etcd节点的`/var/lib/etcd`目录删掉, `--data-dir`就可以写成`/var/lib/etcd`(相当于直接替换掉原来的数据目录), 然后启动etcd进程即可.

### 多节点

见参考文章5, 好像还挺复杂.
