# pod-qos服务质量[qosClass request limit burstable]

参考文章

1. [kubernetes 之QoS服务质量管理](https://www.cnblogs.com/tylerzhou/p/11043280.html)
    - QoS的3种选项: `Guaranteed` > `Burstable` > `BestEffort`
2. [kubenetes之配置pod的QoS](https://www.cnblogs.com/tylerzhou/p/11043282.html)
3. [Configure Quality of Service for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/)
4. [Configure Out of Resource Handling](https://kubernetes.io/docs/tasks/administer-cluster/out-of-resource/)
5. [Assign CPU Resources to Containers and Pods](https://kubernetes.io/docs/tasks/configure-pod-container/assign-cpu-resource/)
    - Pod 调度行为受到`request`的影响(而不是`limit`).
    - `cpu`取值的规则与示例
    - `apiserver`插件: `metrics-server`

## 1. 引言

`QoS`: Quality of Service, 服务质量.

在kubernetes中, 每个POD都有个QoS标记, 通过这个Qos标记来对POD进行服务质量管理. 它取决于用户对服务质量的预期, 也就是期望的服务质量. 

简单来说就是, 当你对一个 Pod 的性能表现的期望比较确定, 希望为ta分配指定大小的资源(也许在这个资源配比下项目组经过了严密的压力测试, 性能达标), 那么`kuber`会为你的服务提供更加可靠的资源保障. 而很多时候, 一些不那么重要的, 测试性的 Pod 没有那么多规矩, 一般都不会指定ta的`resources.{requests,limits}`, 那么`kuber`可能会在资源紧张的场景下挤压(甚至说掠夺)这些Pod的资源.

当你指定的资源分配越笃定, 越规范, kuber就越会保障这样的 Pod 更好地运行.

kuber中只有Pod拥有`QoS`属性, 且体现在两个指标: CPU和内存.

按照资源指定的严密性, kuber中的`QoS`分为3个级别: `Guaranteed` > `Burstable` > `BestEffort`.

```yaml
status:
  qosClass: BestEffort
```

## 2. QoS判定规则

### 2.1 `BestEffort`

`resources`啥都不设置就是`BestEffort`. 不过需要该Pod下所有 container 都不设置才行, 万一设置了某个 container 的`requests`或`limits`, 就不是了.

### 2.2 `Burstable`

Pod 中任意 container 设置了`requests`或`limits`.

## 2.3 `Guaranteed`

要成为`Guaranteed`, 需要达到的条件比较多, 如下

1. Pod 中所有 container **都**设置了`requests`和`limits`
2. 每个 container 自身的`requests.cpu`和`limits.cpu`值相同, `requests.memory`与`limits.memory`值相同(**`cpu`与`memory`都要有**).
3. container 之间的`cpu`和`memory`不必相同.

------

`BestEffort`的最简单, `Burstable`和`Guaranteed`的判断可能复杂一些, 接下来着重解释.

## 3. 示例

### 3.1 注意点1

对于`resources`资源配置, 如果只设置`limits`而不设置`requests`, kuber会自动将`requests`的取值与`limits`保持一致.

如下单 container 的 Pod 将为`Guaranteed`类型.

```yaml
      containers:
      - name: c1
        resources:
          limits:
            cpu: 0.5
            memory: 100Mi
```

### 3.2 注意点2

虽然如果不设置`requests`时会自动取`limits`的值, 但是如果`limits`下只设置了`cpu`或`memory`的**其中一个**, 也是不行的.

如下2种配置都将得到`Burstable`

```
    containers:        |    containers:           
    - name: c1         |    - name: c1            
      resources:       |      resources:          
        limits:        |        limits:           
          cpu: 0.5     |          memory: 100Mi   
        requests:      |        requests:         
          cpu: 0.5     |          memory: 100Mi   
```

当然, 多 container 中, 一个设置`limits.cpu`, 另一个设置`limits.memory`就更不行了.

### 3.3 注意点3

所以各 container 的如果不设置`requests`块, 只要`limits.cpu`和`limits.memory`相同也是会被标识为`Guaranteed`的.

如下配置将得到`Guaranteed`标记.

```yaml
      containers:
      - name: c1
        image: generals/centos7
        resources:
          limits:
            cpu: 1
            memory: 200Mi
      - name: c2
        image: generals/centos7
        resources:
          limits:
            cpu: 0.5
            memory: 100Mi
```

> 注意: `cpu`与`memory`都要有. 

## 3. QoS作用

### BestEffort

当NODE节点上内存资源充足的时候, QoS级别是`BestEffort`的POD可以使用节点上剩余的所有内存资源(因为没有设置`limit`, 所以理论上没有上限). 

当NODE节点上内存资源不够的时候, QoS级别是BestEffort的POD会最先被kill掉; 

### Burstable

当NODE节点上内存资源充足的时候, QoS级别是Burstable的POD会按照requests和limits的设置来使用. 

当NODE节点上内存资源不够的时候, 如果QoS级别是BestEffort的POD已经都被kill掉了, 那么会查找QoS级别是Burstable的POD, 并且这些POD使用的内存已经超出了requests设置的内存值, 这些POD会被kill掉; 

