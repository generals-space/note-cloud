# pv删除失败处于Failed状态-host_path deleter only supports tmp.+ but received provided

kube: 1.23.4

## 问题描述

创建zk集群在删除后, 程序会清理落盘数据、pv、pvc等资源, 但在测试期间, 有几次存在pv残留的情况, 且处于 Failed 状态, 偶现.

```log
root@ly-xjf-r020801-gyt[~]kwd pv | grep zk-01-pv
zk-01-pv-0                20Gi       RWO            Delete           Failed        zjjpt-zk/zk-01-pvc-0                                  local-storage                            3m52s   Filesystem
zk-01-pv-1                20Gi       RWO            Delete           Failed        zjjpt-zk/zk-01-pvc-1                                  local-storage                            3m52s   Filesystem
zk-01-pv-2                20Gi       RWO            Delete           Failed        zjjpt-zk/zk-01-pvc-2                                  local-storage                            3m52s   Filesystem
```

describe pv 有如下信息.

```log
Events:
  Type     Reason                    Age   From                          Message
  ----     ------                    ----  ----                          -------
  Warning  OwnerRefInvalidNamespace  35m   garbage-collector-controller  ownerRef [zookeeper.middleware.com/v1/ZookeeperCluster, namespace: , name: zk-01, uid: fa883b5e-42be-4181-a6be-e8a49027bdb2] does not exist in namespace ""
  Warning  VolumeFailedDelete        32m   persistentvolume-controller   host_path deleter only supports /tmp/.+ but received provided /data/zk/zjjpt-zk/zk-01
```

## 解决方案

这个问题的原因在于, pv 的`persistentVolumeReclaimPolicy`字段设置为了`Delete`, 表示在其绑定的 pvc 被删除时, pv 也会自动被删除.

但是 hostPath 的回收器在高版本 k8s 中被设计只能**自动回收**指向`/tmp/*`目录的 pv, 而上面的出问题的 pv 则指向了`/data/zk/zjjpt-zk/zk-01`目录, 这是不被允许的.

我们的程序中设定了, 在集群删除时会手动删除关联的 pv, 理论上是没有问题的. 

但有可能出现在程序运行到这一步前, 由于先删除了 pvc(pvc是根据 ownerRef 自动回收的), 导致 k8s 根据 ReclaimPolicy 尝试自动回收 pv, 然后失败, 此时 deletionTimestamp 已被设置(只不过被 finalizer 阻塞住了), 程序之后再删除 pv 就没有了效果.

经测试这个问题在 1.17.2 版本下正常, 1.23.4 及 1.28.4 版本会出现.
