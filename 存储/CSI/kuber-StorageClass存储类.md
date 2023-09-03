# kuber-StorageClass存储类

参考文章

1. [Kubernetes笔记（二十二）－－StorageClass设置默认存储后端及动态创建存储](https://blog.csdn.net/bigdata_mining/article/details/96973871)
    - 讲解了sc的作用及使用方法, 不过对于使用场景描述得有点模糊.
2. [CSI - Container Storage Interface（容器存储接口）](https://jimmysong.io/kubernetes-handbook/concepts/csi.html)
    - 很全面的文档, 不过太书面化, 比较难入门

## 引言

我们知道PV应该与PVC绑定使用, PV表示物理资源, 一般需要管理者手动创建才能允许PVC对象使用. 但是很多时候这样实在太不方便, 比如本地测试时, 我只想直接使用`hostPath`的形式外挂宿主机目录, 结果每建一个pod都要手动创建pv(没错, 因为pv中填写的物理机路径需要事先存在, 这个检查是没法避免的), 当然, NFS存储也是这样的道理.

那有没有一种方法, 可以事先划分一大块空间, 然后通过PVC挂载时自动在这块空间中创建需要的目录然后挂载呢?

当然有, 这种手段在kuber中称为`StorageClass(存储类)`.

## volumeBindingMode

`WaitForFirstConsumer`: 单纯创建pvc时不会自动创建pv, 也不会创建实际目录, 此时pvc的状态保持为`Pending`. 只有当pod引用此pvc时才会创建, 同时pvc的状态会变为`Bound`. 不过pod被删除后, 遗留的pvc仍然为`Bound`状态.
