参考文章

1. [Kubernetes 通过statefulset部署redis cluster集群](https://www.cnblogs.com/kuku0223/p/10906003.html)
    - statefulset+NFS存储部署redis-cluster集群示例
2. [官方文档 Headless Services](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services)
3. [kubernetes之StatefulSet详解](https://blog.51cto.com/newfly/2140004)
    - 为什么需要 headless service 无头服务？
    - 为什么需要volumeClaimTemplate？
    - ss的几个重点: Pod管理策略, 更新策略.
4. [Kubernetes资源对象：StatefulSet](https://blog.csdn.net/fly910905/article/details/102092570)
    - 应用场景: 稳定的持久化存储, 稳定的网络标识, 有序部署与有序收缩.
    - 更新策略, 解释了`.spec.updateStrategy.rollingUpdate.partition`的作用
5. [Kubernetes指南 StatefulSet](https://feisky.gitbooks.io/kubernetes/concepts/statefulset.html)
  - 给出了更新策略中的`partition`和管理策略中的`parallel`的使用示例.

对于redis, mysql这种有状态的服务,我们使用`statefulset`方式为首选. 我们这边主要就是介绍`statefulset`这种方式, `statefulset`的设计原理模型:

1. 拓扑状态: 应用的多个实例之间不是完全对等的关系, 这个应用实例的启动必须按照某些顺序启动. 比如应用的主节点A要先于从节点B启动. 而如果你把A和B两个Pod删除掉, 他们再次被创建出来是也必须严格按照这个顺序才行. 并且, 新创建出来的Pod, 必须和原来的Pod的网络标识一样, 这样原先的访问者才能使用同样的方法, 访问到这个新的Pod. 
2. 存储状态: 应用的多个实例分别绑定了不同的存储数据. 对于这些应用实例来说, Pod A第一次读取到的数据, 和隔了十分钟之后再次读取到的数据, 应该是同一份, 哪怕在此期间Pod A被重新创建过. 一个数据库应用的多个存储实例.

## 关于headless service

在statefulset中, headless service也是非常重要的一个点. 其实headless service就是普通的`Service`资源, 且类型为`ClusterIP`, 只不过把`clusterIP`字段显示地设置为了`None`. 

```
$ kubectl get svc redis-service
NAME            TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
redis-service   ClusterIP   None         <none>        6379/TCP   53s
```

可以看到, 服务名称为`redis-service`, 其`CLUSTER-IP`为`None`, 表示这是一个"无头"服务.

`headless service`不具有普通Service资源的负载均衡能力, 因为没有`ClusterIP`, 所以`kube-proxy`组件不处理此类服务, 在访问此类服务的时候返回的是所有后端Pod的Endpoints列表.

访问一个普通的`Service`, kube-proxy会将请求重定向到后端的某个`Pod`, 多次请求虽然发送到的后端可能不同, 但是前端是无感知的, 因为Service本身有固定IP.

但是访问一个`headless service`, 其实是随机且直接访问到后端`Pod`, 比如多次`ping redis-service`, 你会发现解析出来的地址是不同的, 而这些地址都是Pod的地址.

```
$ ping redis-service
PING redis-service.default.svc.cluster.local (10.254.0.215) 56(84) bytes of data.
64 bytes from redis-app-5.redis-service.default.svc.cluster.local (10.254.0.215): icmp_seq=3 ttl=64 time=0.081 ms
...

$ ping redis-service
PING redis-service.default.svc.cluster.local (10.254.0.213) 56(84) bytes of data.
64 bytes from redis-app-3.redis-service.default.svc.cluster.local (10.254.0.213): icmp_seq=2 ttl=64 time=0.085 ms
```

## 更新策略

`statefulset`目前支持两种策略

- `OnDelete`: 当`.spec.template`更新时, 并不立即删除旧的`Pod`, 而是等待用户手动删除这些旧`Pod`后自动创建新Pod. 这是默认的更新策略, 兼容 v1.6 版本的行为.
- `RollingUpdate`: 当`.spec.template`更新时, 自动删除旧的`Pod`并创建新`Pod`替换. 在更新时, 这些`Pod`是按逆序的方式进行, 依次删除、创建并等待`Pod`变成`Ready`状态才进行下一个`Pod`的更新. 

其中`RollingUpdate`有一个`partition`选项, 只有序号大于或等于`partition`的`Pod`会在`.spec.template`更新的时候滚动更新, 而其余的`Pod`则保持不变(即便是删除后也是用以前的版本重新创建).

```yaml
spec:
  ## headless service名称
  serviceName: "redis-service"
  replicas: 6
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      partition: 4
```

> `partition`从0开始计数

这样, 在更新images版本后apply, 你会发现只有`redis-app-4`和`redis-app-5`会更新, 其他的`Pod`则保持不动.
