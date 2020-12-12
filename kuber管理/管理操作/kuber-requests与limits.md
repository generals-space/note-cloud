# kuber-requests与limits 

kuber: 1.16.2, 单节点集群, 4c8g

本文章中不涉及`ResourceQuota`配额限制, 全部都是默认状态.

## 1. 认知

首先, 我们知道`requests`与`limits`分别确定了一个Pod的下限和上限, kuber 会为设置了`requests`的Pod从这个Pod被调度到的节点上划分出`requests`指定大小的资源, 由其独享, 且会根据`limits`限制其资源总量.

不过要注意, 不管是`requests`还是`limits`, 都是预分配的块, 并不是指此Pod实际使用的资源, 我们可以称其为"逻辑资源".

### 1.1

一个node节点上所有Pod的`requests`资源总和不可大于(也没法等于, 即只能小于)该node本身实际的资源, 否则之后调度到该node上的Pod将处于`Pending`状态, `describe`信息时会发现`insufficient`事件.

### 1.2 

`limits`对于调度没有影响, 可以写任意大的值, 而且`limits`会真实地限制一个Pod能够使用的最大资源. 

不过如果`limits`都写了非常大, 大到都超过了物理机实际拥有的资源, 而Pod中的资源占用总和超过了node节点的物理资源, 可能会引发系统级的OOM.

### 1.3 

假设某个node节点上存在n个Pod, 他们的`requests`资源总和小于该node的物理资源, 但是`limits`值设置得非常大, 总和远远超过了物理资源值. 当某个Pod的实际资源占用值已经达到了该node的物理资源上限时, 其他Pod会不会被重新调度呢? 如果另一个Pod也开始发力, 资源占用上升, k8s会不会将这两个Pod中的其中一个挤出这个节点呢?

答案是"不会".

`requests`与`limits`的设置值不同时, k8s将此类Pod定义为`Burstable`(可爆发型), 可以通过Pod信息的`status.qosClass`字段查看. 上面那种情况中, 所有Pod将会共享该node上的物理资源, 发生抢占, 只是这样的话, Pod中的各业务进程没有办法利用全部的资源, 性能必然不佳.

### 1.4 

k8s优先保障`requests`请求满足, 当`requests`与`limits`的定义值相同时, k8s会将此类Pod定义为`Guaranteed`(性能保障型).

~~当一个node节点上存在`Guaranteed`的Pod时, 其他非`Guaranteed`的Pod无论如何抢占, 都只能抢占 [node物理资源 - `GuaranteedPod.requests`] 即剩下的资源了, 就算Guaranteed Pod上没运行什么耗资源的进程, 也没用, k8s 会为ta预留足够的资源.~~

不对, 当`Guaranteed`的Pod上没运行什么耗资源的进程时, 其他非`Guaranteed`的Pod可以抢占node节点上的所有资源. 一旦`Guaranteed`的Pod开始发力, k8s 会保障此类Pod能获得的最大资源, 其余Pod就只能在[node物理资源 - `GuaranteedPod.requests`] 即剩下的资源中抢占了.

### 1.5 (猜测)

当`Burstable` Pod 发生抢占时, 他们各自的资源上限貌似与`requests`相关?

比如: node资源 2c, Pod A requests 100m/100mi, Pod B requests 500m/500mi.

当Pod A与B 马力全开时, ta们各自的资源占用比例大概是 1:5

按照这个逻辑, 一个Pod可使用的最大资源其实是可以由`requests`保证的, 最低也不会低于这个值.

## 一些误解

### node.status.capacity

使用`k get node 节点名称 -o yaml`可以查看一个node节点的详细信息, 其中`status`块有如下信息.

```yaml
status:
  allocatable:
    cpu: "4"
    ephemeral-storage: "45398517276"
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 7887740Ki
    pods: "110"
  capacity:
    cpu: "4"
    ephemeral-storage: 48106Mi
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 7990140Ki
    pods: "110"
```

`capacity`和`allocatable`其实都算是逻辑资源的状态, 我最初以为`allocatable.cpu == capacity.cpu`是因为该node上的Pod没有指定`requests`块, 但后来我发现这两部分的数据总是相近, 根本不会发生变化...

### metrics-server

在做实验时, 我想看看多个Pod实际资源占用上升, 即将超过node节点的物理资源时, 会不会发生重新调度, 于是部署了`metrics-server`.

但是`metrics-server`查看的是Pod/Node的实际资源占用, 不是逻辑资源, 其实帮助不大.
