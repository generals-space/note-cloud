参考文章

1. [Kubernetes 通过statefulset部署redis cluster集群](https://www.cnblogs.com/kuku0223/p/10906003.html)
    - statefulset+NFS存储部署redis-cluster集群示例
2. [官方文档 Headless Services](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services)
3. [kubernetes之StatefulSet详解](https://blog.51cto.com/newfly/2140004)
    - 为什么需要 headless service 无头服务？
    - 为什么需要volumeClaimTemplate？
    - ss的几个重点: Pod管理策略, 更新策略.

20200620更新

今天终于想到 headlesss service 有什么用了...

像 elasticsearch, etcd 这种分布式服务, 在集群初期 setup 时, 配置文件中就要写上集群中所有节点的IP(或是域名).

比如 es

```yaml
network.host: 0.0.0.0
http.port: 9200
cluster.initial_master_nodes: ["es-01"]
```

再像`etcd`

```yaml
listen-peer-urls: https://172.16.43.101:2380
initial-advertise-peer-urls: https://172.16.43.101:2380
initial-cluster: k8s-master-43-101=https://172.16.43.101:2380,k8s-master-43-102=https://172.16.43.102:2380,k8s-master-43-103=https://172.16.43.103:2380
```

但是由于`kuber`集群的特性, Pod 是没有固定IP的, 所以配置文件里不能写IP. 但是用 Service 也不合适, 因为 Service 作为 Pod 前置的 LB, 一般是为一组后端 Pod 提供访问入口的, 而且 Service 的`selector`也没有办法区别同一组 Pod, 没有办法为每个 Pod 创建单独的 Serivce.

于是有了 Statefulset. ta为每个 Pod 做一个编号, 就是为了能在这一组服务内部区别各个 Pod, 各个节点的角色不会变得混乱. 

同时创建所谓的 headless service 资源, 这个 headless service 不分配 ClusterIP, 因为根本不会用到. 集群内的节点是通过`Pod名称+序号.Service名称`确定彼此进行通信的, 只要序号不变, 访问就不会出错.

------

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
