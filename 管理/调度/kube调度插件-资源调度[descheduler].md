# pod资源调度插件

参考文章

1. [资源调度](https://jimmysong.io/kubernetes-handbook/concepts/scheduling.html)
2. [解决 K8s 调度不均衡问题](https://www.cnblogs.com/fengjian2016/p/16408738.html)
    - 官方策略`BalancedResourceAllocation`多节点资源均衡调度, 是根据Pod requests 资源进行评分的, 而不是按 Node 当前资源水位进行调度.

向集群中新增worker节点时, 正在运行的pod不会自动调度到新节点上, 所以一开始新节点上会是空闲的.

想要让集群中的节点的资源利用率比较均衡一些, 想要将一些高负载的节点上的pod驱逐到新增节点上, 这是kuberentes的scheduler所不支持的, 需要使用如[descheduler](https://github.com/kubernetes-incubator/descheduler)这样的插件来实现.

想要运行一些大数据应用, 设计到资源分片, pod需要与数据分布达到一致均衡, 避免个别节点处理大量数据, 而其它节点闲置导致整个作业延迟, 这时候可以考虑使用[kube-batch](https://github.com/kubernetes-incubator/kube-batch).
