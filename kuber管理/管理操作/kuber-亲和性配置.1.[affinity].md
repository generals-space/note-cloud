# kuber-亲和性配置[affinity]

参考文章

1. [K8S高级调度——亲和性和反亲和性](https://www.jianshu.com/p/61725f179223)
    - 亲和性与反亲和性的应用场景
2. [Kubernetes调度之亲和性和反亲和性](https://johng.cn/kubernetes-affinity-anti-affinity/)
    - 运行时调度策略: `nodeAffinity（主机亲和性）`，`podAffinity（POD亲和性）`以及`podAntiAffinity（POD反亲和性）`
    - 3个示例, 但都不完整
    - `RequiredDuringSchedulingRequiredDuringExecution`还不支持, 看写文章的时间, 应该是`1.14`版本之前吧.
3. [k8s之pod亲和性与反亲和性的topologyKey](https://blog.csdn.net/asdfsadfasdfsa/article/details/106027367)
    - 亲和性/反亲和性中`topologyKey`字段的含义.

运行时调度策略有3种: `nodeAffinity（主机亲和性）`，`podAffinity（POD亲和性）`以及`podAntiAffinity（POD反亲和性）`

- `nodeAffinity`: 主要解决`Pod`要部署在哪些主机，以及`Pod`不能部署在哪些主机上的问题，处理的是`Pod`和主机之间的关系. 
- `podAffinity`: 主要解决`Pod`可以和哪些`Pod`部署在同一个拓扑域中的问题（拓扑域用主机标签实现，可以是单个主机，也可以是多个主机组成的cluster、zone等. ），
- `podAntiAffinity`: 主要解决`Pod`不能和哪些`Pod`部署在同一个拓扑域中的问题. 它们处理的是Kubernetes集群内部`Pod`和`Pod`之间的关系. 

每种亲和性和反亲和性策略都有三种规则可以设置

- `RequiredDuringSchedulingRequiredDuringExecution`: 在调度期间要求满足亲和性或者反亲和性规则，如果不能满足规则，则POD不能被调度到对应的主机上. 在之后的运行过程中，如果因为某些原因（比如修改label）导致规则不能满足，系统会尝试把POD从主机上删除(现在版本还不支持). 
- `RequiredDuringSchedulingIgnoredDuringExecution`: 在调度期间要求满足亲和性或者反亲和性规则，如果不能满足规则，则POD不能被调度到对应的主机上. 在之后的运行过程中，系统不会再检查这些规则是否满足. 
    - 比如最开始创建时时希望业务 Pod 部署在与 redis 相同的 node 上, 之后在运行期间该 node 上的 Pod 删除并被调度到其他 node 上了, 这种情况下业务 Pod 并不会随 redis 重新调度.
- `PreferredDuringSchedulingIgnoredDuringExecution`: 在调度期间尽量满足亲和性或者反亲和性规则，如果不能满足规则，POD也有可能被调度到对应的主机上. 在之后的运行过程中，系统不会再检查这些规则是否满足. 

示例

## nodeAffinity

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ## deploy 生成的 pod 的名称也是 centos-deploy-xxx
  name: centos-deploy
  labels:
    app: centos-deploy
spec:
  replicas: 3
  selector:
    matchLabels:
      ## 这里的 label 是与下面的 template -> metadata -> label 匹配的,
      ## 表示一种管理关系
      app: centos-pod
  template:
    metadata:
      labels:
        app: centos-pod
    spec:
      ## affinity 与 containers 平级
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/hostname
                operator: In
                values:
                - k8s-worker-01
                - k8s-worker-02
      containers:
      - name: centos
        image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7
        command: ['tail', '-f', '/etc/os-release']
```

> 除了`matchExpressions`运算之外, `nodeSelectorTerms`的规则其实还有`matchFields`运算, 目前还没有遇到过.

## podeAffinity

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ## deploy 生成的 pod 的名称也是 centos-deploy-xxx
  name: centos-deploy
  labels:
    app: centos-deploy
spec:
  replicas: 3
  selector:
    matchLabels:
      ## 这里的 label 是与下面的 template -> metadata -> label 匹配的,
      ## 表示一种管理关系
      app: centos-pod
  template:
    metadata:
      labels:
        app: centos-pod
    spec:
      ## affinity 与 containers 平级
      affinity:
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          ## topologyKey 用于筛选 node(通过 label), 目标 node 必须要有 kubernetes.io/hostname 标签.
          - topologyKey: kubernetes.io/hostname
            ## 这里其实还是筛选的目标 node, 
            ## 不过这个 node 的要求是**必须**要存在 k8s-app 标签值在 [kube-dns] 的 **Pod**
            ## 注意, key 字段的值应该是某类 Pod 的标签名.
            labelSelector:
              matchExpressions:
              - key: k8s-app
                operator: In
                values:
                - kube-dns
      containers:
      - name: centos
        image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7
        command: ['tail', '-f', '/etc/os-release']
```

> 除了`matchExpressions`运算之外, `labelSelector`其实还有`matchLabels`运算, 目前还没有遇到过.
