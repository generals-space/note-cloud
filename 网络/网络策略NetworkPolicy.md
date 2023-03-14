# 网络策略NetworkPolicy

参考文章

1. [网络策略](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
    - 官方文档
    - 网络策略通过网络插件来实现, 所以用户必须使用支持`NetworkPolicy`的网络解决方案(简单地创建资源对象, 而没有控制器来使它生效的话, 是没有任何作用的).
    - 有中文版, 但是关于对示例中策略的解释还是要看英文的...否则完全不明白在讲什么
    - 中文版页面的"网络策略入门指南"链接不存在, 英文版链接正常, 不过目标页面不是同一个了.
    - [Declare Network Policy](https://kubernetes.io/docs/tasks/administer-cluster/declare-network-policy/)
2. [借助 Calico，管窥 Kubernetes 网络策略](https://blog.fleeto.us/post/network-policy-basic-calico/)
    - 以Nginx为例简单介绍NetworkPolicy的基本使用方法
    - 该文文末给出的链接已经不存在了, 可以见参考文章1中的[Declare Network Policy]()链接
3. [calico 网络结合 k8s networkpolicy 实现租户隔离及部分租户下业务隔离](https://blog.csdn.net/qianggezhishen/article/details/80390598)
4. [Calico官方文档 Network policy](https://docs.projectcalico.org/v3.10/reference/resources/networkpolicy)

## 引言

默认情况下, Pod之间互相访问是没有任何阻拦的. 

一旦定义网络策略, 会使得只允许满足条件的双方可以通信, 其他都会被阻拦, 即 default deny.

所以, 网络策略声明一定要精准, 不然非常有可能出现"误伤".

但是网络策略生效规则还是比较拧巴, 跟Qos一样, 所以还是要根据实例学习一下.

## 示例1

假设存在如下2个Pod.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  namespace: default
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx
---
apiVersion: v1
kind: Pod
metadata:
  name: centos
  namespace: default
  labels:
    app: centos
spec:
  containers:
  - name: centos
    image: centos:7
    command:
    - sleep
    - "36000"
```

默认情况下, 双方通信没有任何限制. 

此时添加一个策略.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: default
spec:
  ## 当前策略影响到的Pod.
  ## 如果为空 {}, 则会影响到当前 ns 的所有 Pod, 一定要谨慎.
  podSelector:
    matchLabels:
      app: nginx
  policyTypes:
  - Ingress
  ingress:
  ## ingress 是一个数组, 可以声明多条.
  - from:
    ## 只允许如下 Pod, 访问 app=nginx 的 80 端口.
    - podSelector:
        matchLabels:
          app: centos
    ports:
    - protocol: TCP
      port: 80
```

有如下几个点可能容易误解:

1. centos访问`nginx:8080`是否受限?
    - 是, `ingress[].from[].podSelector`匹配到了 centos, ta只能访问nginx的80端口.
2. 当前ns下的其他容器访问`nginx:80`是否受限?
    - 是, 只限制了 centos, 没有限制其他容器来源.
3. 其他ns的 centos 容器访问`nginx:80`是否受限?
    - 是
4. nginx反过来访问centos是否受限?
    - 否, 该策略由`spec.podSelector`指定, 只生效于`app=nginx`的入流量, 其他容器不受影响.

## 示例2-多条件, 并集匹配

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: default
spec:
  ## 当前策略影响到的Pod.
  ## 如果为空 {}, 则会影响到当前 ns 的所有 Pod, 一定要谨慎.
  podSelector:
    matchLabels:
      app: nginx
  policyTypes:
  - Ingress
  ingress:
  ## ingress 是一个数组, 可以声明多条.
  - from:
    ## 只允许如下 Pod, 访问 app=nginx 的 80 端口.
    - podSelector:
        matchLabels:
          app: centos
    - namespaceSelector:
        matchLabels:
          ## 所有 ns 都有这个标签, 是 kube 自行添加的.
          kubernetes.io/metadata.name: xxx
    ports:
    - protocol: TCP
      port: 80
```

这里, `from`下有2条规则, 但并不是求交集, 比如只允许`xxx`命名空间下的`app=centos`容器去访问nginx的80端口, 但不是这样的(这样反而缩小了匹配范围, 会拦截更多的来源).

实际上应该是, **`xxx`空间下所有pod都可以访问`nginx:80`, 而`default`空间下还是只有`app=centos`才能访问, 其他的一率不可以**.

------

由于多条件求并集, 所以如下规则会导致, `default`下的所有容器都能访问`nginx:80`, 而不是只有`default`下的`app=centos`才可以访问.

```
    - podSelector:
        matchLabels:
          app: centos
    - namespaceSelector:
        matchLabels:
          kubernetes.io/metadata.name: default
```
