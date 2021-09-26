# node节点宕机时的pod驱散行为[kcm kubelet eviction]

参考文章

1. [谈谈 K8S 的 pod eviction](http://wsfdl.com/kubernetes/2018/05/15/node_eviction.html)
    - 非常透彻, 值得一看
    - K8S `pod eviction` 机制, 某些场景下如节点`NotReady`, 资源不足时, 把 pod 驱逐至其它节点
    - 有两个组件可以发起`pod eviction`: `kube-controller-manager`和`kubelet`, 及这两个场景的具体介绍.
    - `kube-controller-manager`发起的驱逐, 效果需要商榷
2. [kubernetes之node 宕机，pod驱离问题解决](https://www.cnblogs.com/cptao/p/10911959.html)
    - [ ] `kube-controller-manager`的`--pod-eviction-timeout`选项不起作用
    - [x] 部署文件的污点设置`tolerations`, `tolerationSeconds: 10`有效, 当`NotReady`时间过长, 就会重新调度.

理想的情况下, 驱逐对无状态且设计良好的业务方影响很小. 但是并非所有的业务方都是无状态的, 也并非所有的业务方都针对 Kubernetes 优化其业务逻辑. 

例如

1. 对于有状态的业务, 如果没有共享存储, 异地重建后的 pod 完全丢失原有数据; 即使数据不丢失, 对于 Mysql 类的应用, 如果出现双写, 重则破坏数据. 
2. 对于关心 IP 层的业务, 异地重建后的 pod IP 往往会变化, 虽然部分业务方可以利用 service 和 dns 来解决问题, 但是引入了额外的模块和复杂性. 

Node 对于 Pod 的驱逐往往是因为资源压力, 比如磁盘空间不足(CPU不足应该只会降低请求处理速度, 而内存不足则可能引发OOM, 姑且认为只有磁盘吧).

`Evicted`的Pod的事件描述可能如下

```
Events:
  Type     Reason       Age     From                    Message
  ----     ------       ----    ----                    -------
  Nomarl   Scheduled    90m     default-scheduler       Successfully assigned PodXXX to NodeXXX
  Warning  Evicted      90m     kubelet, 主机名          The node was low on resource: [DiskPressure].
```

相关的配置在 Node 上 kubelet 的配置文件中, 如`/var/lib/kubelet/config.yaml`.

```yaml
evictionHard:
  imagefs.available: 15%
  memory.available: 100Mi
  nodefs.available: 10%
  imagefs.available: 5%
evictionPressureTransitionPeriod: 5m0s
```
