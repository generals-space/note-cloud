# kuber-核心概念总结(转)

原文链接

[Kubernetes核心概念总结](https://www.cnblogs.com/zhenyuyaodidiao/p/6500720.html)

## 1. 基础架构

![](https://gitee.com/generals-space/gitimg/raw/master/9a3ce4b73794277bac82ec69dcca592d.png)

## 4. Service

### 4.2 Service代理外部服务

Service不仅可以代理Pod, 还可以代理任意其他后端, 比如运行在Kubernetes外部Mysql、Oracle等. 这是通过定义两个同名的service和endPoints来实现的. 示例如下: 

`redis-service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-service
spec:
  ports:
  - port: 6379
    targetPort: 6379
    protocol: TCP
```

`redis-endpoints.yaml`

```yaml
apiVersion: v1
kind: Endpoints
metadata:
  name: redis-service
subsets:
  - addresses:
    - ip: 10.0.251.145
    ports:
    - port: 6379
      protocol: TCP
```

基于文件创建完Service和Endpoints之后, 在Kubernetes的Service中即可查询到自定义的Endpoints:

```
[root@k8s-master demon]# kubectl describe service redis-service
Name:            redis-service
Namespace:        default
Labels:            <none>
Selector:        <none>
Type:            ClusterIP
IP:            10.254.52.88
Port:            <unset>    6379/TCP
Endpoints:        10.0.251.145:6379
Session Affinity:    None
No events.
[root@k8s-master demon]# etcdctl get /skydns/sky/default/redis-service
{"host":"10.254.52.88","priority":10,"weight":10,"ttl":30,"targetstrip":0}
```

### 4.3 Service内部负载均衡

当Service的Endpoints包含多个IP的时候, 及服务代理存在多个后端, 将进行请求的负载均衡. 默认的负载均衡策略是轮训或者随机(有kube-proxy的模式决定). 同时, Service上通过设置`Service.spec.sessionAffinity=ClientIP`, 来实现基于源IP地址的会话保持. 

### 4.4 发布Service

Service的虚拟IP是由Kubernetes虚拟出来的内部网络, 外部是无法寻址到的. 但是有些服务又需要被外部访问到, 例如web前段. 这时候就需要加一层网络转发, 即外网到内网的转发. Kubernetes提供了`NodePort`、`LoadBalancer`、`Ingress`三种方式. 

1. `NodePort`: 在之前的Guestbook示例中, 已经延时了`NodePort`的用法. `NodePort`的原理是, Kubernetes会在每一个Node上暴露出一个端口: nodePort, 外部网络可以通过(任一Node)[NodeIP]:[NodePort]访问到后端的Service. 
2. `LoadBalancer`: 在`NodePort`基础上, Kubernetes可以请求底层云平台创建一个负载均衡器, 将每个Node作为后端, 进行服务分发. 该模式需要底层云平台(例如GCE)支持. 
3. `Ingress`: 是一种HTTP方式的路由转发机制, 由`Ingress Controller`和HTTP代理服务器组合而成. `Ingress Controller`实时监控Kubernetes API, 实时更新HTTP代理服务器的转发规则. HTTP代理服务器有GCE Load-Balancer、HaProxy、Nginx等开源方案. 

### 4.5 servicede 自发性机制

Kubernetes中有一个很重要的服务自发现特性. 一旦一个service被创建, 该service的service IP和service port等信息都可以被注入到pod中供它们使用. Kubernetes主要支持两种service发现机制: 环境变量和DNS. 

**环境变量方式**

Kubernetes创建Pod时会自动添加所有可用的service环境变量到该Pod中, 如有需要. 这些环境变量就被注入Pod内的容器里. 需要注意的是, 环境变量的注入只发送在Pod创建时, 且不会被自动更新. 这个特点暗含了service和访问该service的Pod的创建时间的先后顺序, 即任何想要访问service的pod都需要在service已经存在后创建, 否则与service相关的环境变量就无法注入该Pod的容器中, 这样先创建的容器就无法发现后创建的service. 

**DNS方式**

Kubernetes集群现在支持增加一个可选的组件——DNS服务器. 这个DNS服务器使用Kubernetes的watchAPI, 不间断的监测新的service的创建并为每个service新建一个DNS记录. 如果DNS在整个集群范围内都可用, 那么所有的Pod都能够自动解析service的域名. Kube-DNS搭建及更详细的介绍请见: 基于Kubernetes集群部署skyDNS服务.

## 5. Deployment

**关于多重升级**

举例, 当你创建了一个`nginx1.7`的Deployment, 要求副本数量为5之后, `Deployment Controller`会逐步的将5个1.7的Pod启动起来; 当启动到3个的时候, 你又发出更新`Deployment`中Nginx到1.9的命令; 这时`Deployment Controller`会立即将已启动的3个1.7Pod杀掉, 然后逐步启动1.9的Pod. Deployment Controller不会等到1.7的Pod都启动完成之后, 再依次杀掉1.7, 启动1.9. 

### 6.3 Persistent Volume和Persistent Volume Claim

理解每个存储系统是一件复杂的事情, 特别是对于普通用户来说, 有时候并不需要关心各种存储实现, 只希望能够安全可靠地存储数据. Kubernetes中提供了`Persistent Volume`和`Persistent Volume Claim`机制, 这是存储消费模式. 

`Persistent Volume`是由系统管理员配置创建的一个数据卷(目前支持HostPath、GCE Persistent Disk、AWS Elastic Block Store、NFS、iSCSI、GlusterFS、RBD), 它代表了某一类存储插件实现; 

而对于普通用户来说, 通过Persistent Volume Claim可请求并获得合适的Persistent Volume, 而无须感知后端的存储实现. 

Persistent Volume和Persistent Volume Claim的关系其实类似于Pod和Node, Pod消费Node资源, Persistent Volume Claim则消费Persistent Volume资源. Persistent Volume和Persistent Volume Claim相互关联, 有着完整的生命周期管理: 

1. 准备: 系统管理员规划或创建一批Persistent Volume; 
2. 绑定: 用户通过创建Persistent Volume Claim来声明存储请求, Kubernetes发现有存储请求的时候, 就去查找符合条件的Persistent Volume(最小满足策略). 找到合适的就绑定上, 找不到就一直处于等待状态; 
3. 使用: 创建Pod的时候使用Persistent Volume Claim; 
4. 释放: 当用户删除绑定在Persistent Volume上的Persistent Volume Claim时, Persistent Volume进入释放状态, 此时Persistent Volume中还残留着上一个Persistent Volume Claim的数据, 状态还不可用; 
5. 回收: 是否的Persistent Volume需要回收才能再次使用. 回收策略可以是人工的也可以是Kubernetes自动进行清理(仅支持NFS和HostPath)

## 7. Pet Sets/StatefulSet

K8s在1.3版本里发布了Alpha版的`PetSet`功能. 在云原生应用的体系里, 有下面两组近义词; 

1. 无状态(stateless)、牲畜(cattle)、无名(nameless)、可丢弃(disposable); 
2. 有状态(stateful)、宠物(pet)、有名(having name)、不可丢弃(non-disposable)

RC和RS主要是控制提供无状态服务的, 其所控制的Pod的名字是随机设置的, 一个Pod出故障了就被丢弃掉, 在另一个地方重启一个新的Pod, 名字变了、名字和启动在哪儿都不重要, 重要的只是Pod总数; 而PetSet是用来控制有状态服务, PetSet中的每个Pod的名字都是事先确定的, 不能更改. PetSet中Pod的名字的作用, 是用来关联与该Pod对应的状态. 

对于RC和RS中的Pod, 一般不挂载存储或者挂载共享存储, 保存的是所有Pod共享的状态, Pod像牲畜一样没有分别; 对于PetSet中的Pod, 每个Pod挂载自己独立的存储, 如果一个Pod出现故障, 从其他节点启动一个同样名字的Pod, 要挂在上原来Pod的存储继续以它的状态提供服务. 

适合于PetSet的业务包括数据库服务MySQL和PostgreSQL, 集群化管理服务Zookeeper、etcd等有状态服务. PetSet的另一种典型应用场景是作为一种比普通容器更稳定可靠的模拟虚拟机的机制. 传统的虚拟机正是一种有状态的宠物, 运维人员需要不断地维护它, 容器刚开始流行时, 我们用容器来模拟虚拟机使用, 所有状态都保存在容器里, 而这已被证明是非常不安全、不可靠的. 使用PetSet, Pod仍然可以通过漂移到不同节点提供高可用, 而存储也可以通过外挂的存储来提供高可靠性, PetSet做的只是将确定的Pod与确定的存储关联起来保证状态的连续性. 
