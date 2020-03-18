参考文章

1. [Kubernetes 通过statefulset部署redis cluster集群](https://www.cnblogs.com/kuku0223/p/10906003.html)
    - statefulset+NFS存储部署redis-cluster集群示例

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
