# kube调度插件-topologySpreadConstraints[schedule plugin]

参考文章

1. [Introducing PodTopologySpread](https://kubernetes.io/blog/2020/05/introducing-podtopologyspread/)
    - IBM开发的调度插件, v1.18版本启用.
2. [PodTopologySpread介绍](https://cloud.tencent.com/developer/article/1631990)
    - 参考文章1的中文版
3. [Pod Topology Spread Constraints](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/)
4. [Pod 拓扑分布约束](https://kubernetes.io/zh-cn/docs/concepts/scheduling-eviction/topology-spread-constraints/)
5. [Pod Topology Spread Constraints介绍](https://cloud.tencent.com/developer/article/1639217)
    - PodTopologySpread 由 EvenPodsSpread 特性门控所控制，在 v1.16 版本第一次发布，并在 v1.18 版本进入 beta 阶段默认启用

## 引言

基于Pod亲和性与反亲和性实现的调度, 仍然过于粗放. 

举例如下, 假设Deployment副本数为6, 通过`nodeSelector`/`nodeAffinity`筛选出的, 符合条件的Node节点数为3.

我们一般希望该Deployment下的6个Pod均匀的分布在这3个节点上, 每个节点2个Pod.

但如果其中一个Node节点是新增的, 比较空闲. 那么原生调度器会结合资源条件给各节点打分, 会出现0:0:6的情况, 完全没有容灾能力.

而如果设置了`podAntiAffinity`, 又因为节点数量不足导致还有3个Pod无法被调度, 一直处于Pending状态.

`topologySpreadConstraints`属性就是为了处理这种**极端情况**出现的.

## 使用方法

下面的约束部分就是为了解决上述问题存在的.

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: general-sts
  labels:
    app: general-sts
  namespace: default
spec:
  replicas: 6
  serviceName: ""
  podManagementPolicy: Parallel
  selector:
    matchLabels:
      app: general-sts
  template:
    metadata:
      labels:
        app: general-sts
    spec:
      affinity:
        ## 主机亲和
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: general-label-key
                operator: In
                values:
                - general-label-value
      containers:
      - name: centos7
        image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7
        imagePullPolicy: IfNotPresent
        command: ["tail", "-f", "/etc/os-release"]
      topologySpreadConstraints:
      ## maxSkew 表示各Node节点上 Pod 数量的差值最大为1, 不能出现某些主机2个, 而某些主机0个的情况.
      ## 不均匀分布的最大程度, 该值必须大于0.
      - maxSkew: 1
        ## 对符合 nodeSelector/nodeAffinity 的主机, 按照 topologyKey 进行区分, 
        topologyKey: kubernetes.io/hostname
        ## 如果出现了无法实现 n:n:n-1:n-1 比例的调度场景, 比如某些节点资源不足, 余下的 Pod 的调度方式
        ## ScheduleAnyway 表示可以按照实际资源进行调度,
        ## 还有一个 DoNotSchedule 选项表示剩下的 Pod 就 Pending 着吧.
        whenUnsatisfiable: ScheduleAnyway
        ## 表示被此约束限制的 Pod, 用于计算 maxSkew 的值, 一般与 Pod 自身的 label 保持一致即可.
        labelSelector:
          matchLabels: 
            app: general-sts
```

## zone

参考文章1到4的示例中, 都提到一个`zone`标签, 有点抽象.

假设 Node 节点中存在两个zone的标签, `zone=北京`和`zone=香港`, 为了保证用户能够更快地访问到服务, 需要将Pod分别平均调度到这两个区域内的Node节点上(先不考虑用户如何访问到最近区域的Pod).

而`zone`标签与上面的`kubernetes.io/hostname`有一点不同, 拥有`zone=北京`标签的主机理论上不只一台, 而`kubernetes.io/hostname`则是唯一的, 所以有可能会出现如下情况

```
      北京              香港
+---------------+---------------+
|     zoneA     |     zoneB     |
+-------+-------+-------+-------+
| node1 | node2 | node3 | node4 |
+-------+-------+-------+-------+
|   P   |   P   |  P P  |       |
+-------+-------+-------+-------+
```

zoneA与zoneB中的Pod数量比例满足约束, 但是 Node 节点上的 Pod 数量才可能失衡, 需要修改`topologyKey`, 或是配合`podAntiAffinity`.
