# node节点宕机时的pod驱散行为

参考文章

1. [谈谈 K8S 的 pod eviction](http://wsfdl.com/kubernetes/2018/05/15/node_eviction.html)
    - 非常透彻, 值得一看
    - K8S `pod eviction` 机制, 某些场景下如节点`NotReady`, 资源不足时, 把 pod 驱逐至其它节点
    - 有两个组件可以发起`pod eviction`: `kube-controller-manager`和`kubelet`, 及这两个场景的具体介绍.
    - `kube-controller-manager`发起的驱逐, 效果需要商榷
2. [kubernetes之node 宕机，pod驱离问题解决](https://www.cnblogs.com/cptao/p/10911959.html)
    - [ ] `kube-controller-manager`的`--pod-eviction-timeout`选项不起作用
    - [x] 部署文件的污点设置`tolerations`, `tolerationSeconds: 10`有效, 当`NotReady`时间过长, 就会重新调度.

> 理想的情况下, 驱逐对无状态且设计良好的业务方影响很小. 但是并非所有的业务方都是无状态的, 也并非所有的业务方都针对 Kubernetes 优化其业务逻辑. 
> 
> 例如, 对于有状态的业务, 如果没有共享存储, 异地重建后的 pod 完全丢失原有数据; 即使数据不丢失, 对于 Mysql 类的应用, 如果出现双写, 重则破坏数据. 对于关心 IP 层的业务, 异地重建后的 pod IP 往往会变化, 虽然部分业务方可以利用 service 和 dns 来解决问题, 但是引入了额外的模块和复杂性. 
