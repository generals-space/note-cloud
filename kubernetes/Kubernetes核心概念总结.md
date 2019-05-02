# Kubernetes核心概念总结

原文链接

[Kubernetes核心概念总结](https://www.cnblogs.com/zhenyuyaodidiao/p/6500720.html)

## 1. 基础架构

![](https://gitee.com/generals-space/gitimg/raw/master/9a3ce4b73794277bac82ec69dcca592d.png)

### 1.1 Master

Master节点上面主要由四个模块组成: `APIServer`、`scheduler`、`controller manager`、`etcd`. 

`APIServer`: `APIServer`负责对外提供RESTful的Kubernetes API服务, 它是系统管理指令的统一入口, 任何对资源进行增删改查的操作都要交给APIServer处理后再提交给etcd. 如架构图中所示, `kubectl`（Kubernetes提供的客户端工具, 该工具内部就是对Kubernetes API的调用）是直接和APIServer交互的. 

`schedule`: `scheduler`的职责很明确, 就是负责调度`pod`到合适的`Node`上. 如果把`scheduler`看成一个黑匣子, 那么它的输入是pod和由多个Node组成的列表, 输出是Pod和一个Node的绑定, 即将这个pod部署到这个Node上. Kubernetes目前提供了调度算法, 但是同样也保留了接口, 用户可以根据自己的需求定义自己的调度算法. 

`controller manager`: 如果说APIServer做的是“前台”的工作的话, 那`controller manager`就是负责“后台”的. 每个资源一般都对应有一个控制器, 而`controller manager`就是负责管理这些控制器的. 比如我们通过APIServer创建一个pod, 当这个pod创建成功后, APIServer的任务就算完成了. 而后面保证Pod的状态始终和我们预期的一样的重任就由`controller manager`去保证了. 

`etcd`: `etcd`是一个高可用的键值存储系统, Kubernetes使用它来存储各个资源的状态, 从而实现了Restful的API. 

### 1.2 Node

每个Node节点主要由三个模块组成: `kubelet`、`kube-proxy`、`runtime`. 

`runtime`: `runtime`指的是容器运行环境, 目前Kubernetes支持docker和rkt两种容器. 

`kube-proxy`: 该模块实现了Kubernetes中的服务发现和反向代理功能. 反向代理方面: `kube-proxy`支持TCP和UDP连接转发, 默认基于`Round Robin`算法将客户端流量转发到与service对应的一组后端pod. 服务发现方面, kube-proxy使用etcd的watch机制, 监控集群中service和endpoint对象数据的动态变化, 并且维护一个service到endpoint的映射关系, 从而保证了后端pod的IP变化不会对访问者造成影响. 另外kube-proxy还支持session affinity. 

`kubelet`: Kubelet是Master在每个Node节点上面的**`agent`**, 是Node节点上面最重要的模块, 它负责维护和管理该Node上面的所有容器, 但是如果容器不是通过Kubernetes创建的, 它并不会管理. 本质上, 它负责使Pod得运行状态与期望的状态一致. 

至此, Kubernetes的Master和Node就简单介绍完了. 下面我们来看Kubernetes中的各种资源/对象. 

## 2. Pod

　Pod 是Kubernetes的基本操作单元, 也是应用运行的载体. 整个Kubernetes系统都是围绕着Pod展开的, 比如如何部署运行Pod、如何保证Pod的数量、如何访问Pod等. 另外, Pod是一个或多个相关容器的集合, 这可以说是一大创新点, 提供了一种容器的组合的模型. 

### 2.1 基本操作

|操作类型|方法|
|:-:|:-|
|创建| `kubectl create -f xxx.yaml` |
|查询| `kubectl get pod Pod名称` <br/> `kubectl describe pod Pod名称` |
|删除| `kubectl delete pod Pod名称` |
|更新| `kubectl replace /path/to/yourNewYaml.yaml` |

### 2.2 Pod与容器

在Docker中, 容器是最小的处理单元, 增删改查的对象是容器, 容器是一种虚拟化技术, 容器之间是隔离的, 隔离是基于Linux Namespace实现的. 

而在Kubernetes中, Pod包含一个或者多个相关的容器, Pod可以认为是容器的一种延伸扩展, 一个Pod也是一个隔离体, 而Pod内部包含的一组容器又是共享的（包括PID、Network、IPC、UTS）. 除此之外, Pod中的容器可以访问共同的数据卷来实现文件系统的共享. 

### 2.3 镜像

在kubernetes中, 镜像的下载策略为: 

- `Always`: 每次都下载最新的镜像

- `Never`: 只使用本地镜像, 从不下载

- `IfNotPresent`: 只有当本地没有的时候才下载镜像

Pod被分配到Node之后会根据镜像下载策略进行镜像下载, 可以根据自身集群的特点来决定采用何种下载策略. 无论何种策略, 都要确保Node上有正确的镜像可用. 

### 2.4 其他设置

通过yaml文件, 可以在Pod中设置: 

1. 启动命令, 如: `spec.containers.command`

2. 环境变量, 如: `spec.containers.env.name/value`

3. 端口桥接, 如: `spec.containers.ports.containerPort/protocol/hostIP/hostPort`（使用`hostPort`时需要注意端口冲突的问题, 不过Kubernetes在调度Pod的时候会检查宿主机端口是否冲突, 比如当两个Pod均要求绑定宿主机的80端口, Kubernetes将会将这两个Pod分别调度到不同的机器上）;

4. Host网络, 一些特殊场景下, 容器必须要以host方式进行网络设置（如接收物理机网络才能够接收到的组播流）, 在Pod中也支持host网络的设置, 如: `spec.hostNetwork=true`；

5. 数据持久化, 如: `spec.containers.volumeMounts.mountPath`;

6. 重启策略, 当Pod中的容器终止退出后, 重启容器的策略. 这里的所谓Pod的重启, 实际上的做法是容器的重建, 之前容器中的数据将会丢失, 如果需要持久化数据, 那么需要使用数据卷进行持久化设置. Pod支持三种重启策略: 
    - Always（默认策略, 当容器终止退出后, 总是重启容器）
    - OnFailure（当容器终止且异常退出时, 重启）
    - Never（从不重启）

### 2.5 Pod生命周期

Pod被分配到一个Node上之后, 就不会离开这个Node, 直到被删除. 当某个Pod失败, 首先会被Kubernetes清理掉, 之后ReplicationController将会在其它机器上（或本机）重建Pod, 重建之后Pod的ID发生了变化, 那将会是一个新的Pod. 所以, Kubernetes中Pod的迁移, 实际指的是在新Node上重建Pod. 以下给出Pod的生命周期图. 

![](https://gitee.com/generals-space/gitimg/raw/master/5b9bfedcca7b6eb04a9a9f5532af42df.png)

**生命周期回调函数**: 

- `PostStart`（容器创建成功后调用该回调函数）

- `PreStop`（在容器被终止前调用该回调函数）

以下示例中, 定义了一个Pod, 包含一个JAVA的web应用容器, 其中设置了`PostStart`和`PreStop`回调函数. 即在容器创建成功后, 复制`/sample.war`到`/app`文件夹中. 而在容器终止之前, 发送HTTP请求到`http://monitor.com:8080/waring`, 即向监控系统发送警告. 具体示例如下: 

```yml
## ...
containers:
- image: sample:v2  
    name: war
    lifecycle: 
      posrStart:
        exec:
          command:
          - “cp”
          - “/sample.war”
          - “/app”
      prestop:
       httpGet:
        host: monitor.com
        psth: /waring
        port: 8080
        scheme: HTTP
```

## 3. Replication Controller

`Replication Controller（RC）`是Kubernetes中的另一个核心概念, 应用托管在Kubernetes之后, Kubernetes需要保证应用能够持续运行, 这是RC的工作内容, 它会确保任何时间Kubernetes中都有指定数量的Pod在运行. 在此基础上, RC还提供了一些更高级的特性, 比如滚动升级、升级回滚等. 

### 3.1 RC与Pod的关联——Label

`RC`与`Pod`的关联是通过`Label`来实现的. `Label`机制是Kubernetes中的一个重要设计, 通过Label进行对象的弱关联, 可以灵活地进行分类和选择. 对于`Pod`, 需要设置其自身的`Label`来进行标识, `Label`是一系列的`Key/value`对, 在`Pod.metadata.labeks`中进行设置. 

Label的定义是任一的, 但是Label必须具有可标识性, 比如设置Pod的应用名称和版本号等. 另外Lable是不具有唯一性的, 为了更准确的标识一个Pod, 应该为Pod设置多个维度的label. 如下: 

- "release" : "stable", "release" : "canary"

- "environment" : "dev", "environment" : "qa", "environment" : "production"

- "tier" : "frontend", "tier" : "backend", "tier" : "cache"

- "partition" : "customerA", "partition" : "customerB"

- "track" : "daily", "track" : "weekly"

举例, 当你在`RC`的yaml文件中定义了该RC的`selector`中的`label`为`app:my-web`, 那么这个RC就会去关注`Pod.metadata.labeks`中label为app:my-web的Pod. 修改了对应Pod的Label, 就会使Pod脱离RC的控制. 同样, 在RC运行正常的时候, 若试图继续创建同样Label的Pod, 是创建不出来的. 因为RC认为副本数已经正常了, 再多起的话会被RC删掉的. 

### 3.2 弹性伸缩

弹性伸缩是指适应负载变化, 以弹性可伸缩的方式提供资源. 反映到Kubernetes中, 指的是可根据负载的高低动态调整Pod的副本数量. 调整Pod的副本数是通过修改RC中Pod的副本是来实现的, 示例命令如下: 

扩容Pod的副本数目到10

```
$ kubectl scale relicationcontroller yourRcName --replicas=10
```

缩容Pod的副本数目到1

```
$ kubectl scale relicationcontroller yourRcName --replicas=1
```

### 3.3 滚动升级

滚动升级是一种平滑过渡的升级方式, 通过逐步替换的策略, 保证整体系统的稳定, 在初始升级的时候就可以及时发现、调整问题, 以保证问题影响度不会扩大. Kubernetes中滚动升级的命令如下: 

```
$ kubectl rolling-update my-rcName-v1 -f my-rcName-v2-rc.yaml --update-period=10s
```

升级开始后, 首先依据提供的定义文件创建V2版本的RC, 然后每隔10s（`--update-period=10s`）逐步的增加V2版本的Pod副本数, 逐步减少V1版本Pod的副本数. 升级完成之后, 删除V1版本的RC, 保留V2版本的RC, 及实现滚动升级. 

升级过程中, 发生了错误中途退出时, 可以选择继续升级. Kubernetes能够智能的判断升级中断之前的状态, 然后紧接着继续执行升级. 当然, 也可以进行回退, 命令如下: 

```
$ kubectl rolling-update my-rcName-v1 -f my-rcName-v2-rc.yaml --update-period=10s --rollback
```

回退的方式实际就是升级的逆操作, 逐步增加V1.0版本Pod的副本数, 逐步减少V2版本Pod的副本数. 

> 连配置文件都一样, 只是指定了`--rollback`参数.

### 3.4 新一代副本控制器replica set

这里所说的`replica set`, 可以被认为 是“升级版”的`Replication Controller`. 也就是说. `replica set`也是用于保证与`label selector`匹配的pod数量维持在期望状态. 区别在于, `replica set`引入了对基于子集的`selector`查询条件, 而Replication Controller仅支持基于值相等的selecto条件查询. 这是目前从用户角度肴, 两者唯一的显著差异.  

社区引入这一API的初衷是用于取代vl中的`Replication Controller`, 也就是说．当v1版本被废弃时, `Replication Controller`就完成了它的历史使命, 而由`replica set`来接管其工作. 虽然`replica set`可以被单独使用, 但是目前它多被Deployment用于进行pod的创建、更新与删除. Deployment在滚动更新等方面提供了很多非常有用的功能, 关于Deployment的更多信息, 读者们可以在后续小节中获得. 

## 4. Service

为了适应快速的业务需求, 微服务架构已经逐渐成为主流, 微服务架构的应用需要有非常好的服务编排支持. Kubernetes中的核心要素Service便提供了一套简化的服务代理和发现机制, 天然适应微服务架构. 

### 4.1 原理

在Kubernetes中, 在受到RC调控的时候, Pod副本是变化的, 对于的虚拟IP也是变化的, 比如发生迁移或者伸缩的时候. 这对于Pod的访问者来说是不可接受的. Kubernetes中的Service是一种抽象概念, 它定义了一个Pod逻辑集合以及访问它们的策略, `Service`同`Pod`的关联同样是居于`Label`来完成的. Service的目标是提供一种桥梁,  它会为访问者提供一个固定访问地址, 用于在访问时重定向到相应的后端, 这使得非 Kubernetes原生应用程序, 在无须为Kubemces编写特定代码的前提下, 轻松访问后端. 

`Service`同RC一样, 都是通过Label来关联Pod的. 当你在Service的yaml文件中定义了该Service的selector中的label为app:my-web, 那么这个Service会将`Pod.metadata.labeks`中`label`为`app:my-web`的Pod作为分发请求的后端. 当Pod发生变化时（增加、减少、重建等）, Service会及时更新. 这样一来, `Service`就可以作为`Pod`的访问入口, 起到代理服务器的作用, 而对于访问者来说, 通过`Service`进行访问, 无需直接感知`Pod`. 

需要注意的是, `Kubernetes`分配给`Service`的固定IP是一个虚拟IP, 并不是一个真实的IP, 在外部是无法寻址的. 真实的系统实现上, Kubernetes是通过Kube-proxy组件来实现的虚拟IP路由及转发. 所以在之前集群部署的环节上, 我们在每个Node上均部署了`Proxy`这个组件, 从而实现了Kubernetes层级的虚拟转发网络. 

### 4.2 Service代理外部服务

Service不仅可以代理Pod, 还可以代理任意其他后端, 比如运行在Kubernetes外部Mysql、Oracle等. 这是通过定义两个同名的service和endPoints来实现的. 示例如下: 

`redis-service.yaml`

```yml
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

```yml
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

当Service的Endpoints包含多个IP的时候, 及服务代理存在多个后端, 将进行请求的负载均衡. 默认的负载均衡策略是轮训或者随机（有kube-proxy的模式决定）. 同时, Service上通过设置`Service.spec.sessionAffinity=ClientIP`, 来实现基于源IP地址的会话保持. 

### 4.4 发布Service

Service的虚拟IP是由Kubernetes虚拟出来的内部网络, 外部是无法寻址到的. 但是有些服务又需要被外部访问到, 例如web前段. 这时候就需要加一层网络转发, 即外网到内网的转发. Kubernetes提供了`NodePort`、`LoadBalancer`、`Ingress`三种方式. 

1. `NodePort`: 在之前的Guestbook示例中, 已经延时了`NodePort`的用法. `NodePort`的原理是, Kubernetes会在每一个Node上暴露出一个端口: nodePort, 外部网络可以通过（任一Node）[NodeIP]:[NodePort]访问到后端的Service. 

2. `LoadBalancer`: 在`NodePort`基础上, Kubernetes可以请求底层云平台创建一个负载均衡器, 将每个Node作为后端, 进行服务分发. 该模式需要底层云平台（例如GCE）支持. 

3. `Ingress`: 是一种HTTP方式的路由转发机制, 由`Ingress Controller`和HTTP代理服务器组合而成. `Ingress Controller`实时监控Kubernetes API, 实时更新HTTP代理服务器的转发规则. HTTP代理服务器有GCE Load-Balancer、HaProxy、Nginx等开源方案. 

### 4.5 servicede 自发性机制

Kubernetes中有一个很重要的服务自发现特性. 一旦一个service被创建, 该service的service IP和service port等信息都可以被注入到pod中供它们使用. Kubernetes主要支持两种service发现 机制: 环境变量和DNS. 

**环境变量方式**

Kubernetes创建Pod时会自动添加所有可用的service环境变量到该Pod中, 如有需要．这些环境变量就被注入Pod内的容器里. 需要注意的是, 环境变量的注入只发送在Pod创建时, 且不会被自动更新. 这个特点暗含了service和访问该service的Pod的创建时间的先后顺序, 即任何想要访问service的pod都需要在service已经存在后创建, 否则与service相关的环境变量就无法注入该Pod的容器中, 这样先创建的容器就无法发现后创建的service. 

**DNS方式**

Kubernetes集群现在支持增加一个可选的组件——DNS服务器. 这个DNS服务器使用Kubernetes的watchAPI, 不间断的监测新的service的创建并为每个service新建一个DNS记录. 如果DNS在整个集群范围内都可用, 那么所有的Pod都能够自动解析service的域名. Kube-DNS搭建及更详细的介绍请见: 基于Kubernetes集群部署skyDNS服务

### 4.6 多个service如何避免地址和端口冲突

此处设计思想是, Kubernetes通过为每个service分配一个唯一的ClusterIP, 所以当使用ClusterIP: port的组合访问一个service的时候, 不管port是什么, 这个组合是一定不会发生重复的. 另一方面, kube-proxy为每个service真正打开的是一个绝对不会重复的随机端口, 用户在service描述文件中指定的访问端口会被映射到这个随机端口上. 这就是为什么用户可以在创建service时随意指定访问端口. 

### 4.7 service目前存在的不足

Kubernetes使用`iptables`和`kube-proxy`解析service的人口地址, 在中小规模的集群中运行良好, 但是当service的数量超过一定规模时, 仍然有一些小问题. 首当其冲的便是service环境变量泛滥, 以及service与使用service的pod两者创建时间先后的制约关系. 目前来看, 很多使用者在使用Kubernetes时往往会开发一套自己的Router组件来替代service, 以便更好地掌控和定制这部分功能. 

## 5. Deployment

Kubernetes提供了一种更加简单的更新RC和Pod的机制, 叫做Deployment. 通过在Deployment中描述你所期望的集群状态, `Deployment Controller`会将现在的集群状态在一个可控的速度下逐步更新成你所期望的集群状态. Deployment主要职责同样是为了保证pod的数量和健康, 90%的功能与`Replication Controller`完全一样, 可以看做新一代的`Replication Controller`. 但是, 它又具备了`Replication Controller`之外的新特性: 

1. `RC`全部功能: Deployment继承了上面描述的`RC`的全部功能. 

2. 事件和状态查看: 可以查看Deployment的升级详细进度和状态. 

3. 回滚: 当升级pod镜像或者相关参数的时候发现问题, 可以使用回滚操作回滚到上一个稳定的版本或者指定的版本. 

4. 版本记录: 每一次对Deployment的操作, 都能保存下来, 给予后续可能的回滚使用. 

5. 暂停和启动: 对于每一次升级, 都能够随时暂停和启动. 

6. 多种升级方案: 
    - Recreate: 删除所有已存在的pod, 重新创建新的; 
    - RollingUpdate: 滚动升级, 逐步替换的策略, 同时滚动升级时, 支持更多的附加参数, 例如设置最大不可用pod数量, 最小升级间隔时间等等. 

### 5.1 滚动升级

相比于`RC`, `Deployment`直接使用`kubectl edit deployment/deploymentName`或者`kubectl set`方法就可以直接升级（原理是`Pod`的`template`发生变化, 例如更新`label`、更新镜像版本等操作会触发`Deployment`的滚动升级）. 操作示例——首先 我们同样定义一个`nginx-deploy-v1.yaml`的文件, 副本数量为2: 

```yml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

创建deployment

```
$ kubectl create -f nginx-deploy-v1.yaml --record
deployment "nginx-deployment" created
$ kubectl get deployments
NAME       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   3         0         0            0           1s
$ kubectl get deployments
NAME       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   3         3         3            3           18s
```

正常之后, 将nginx的版本进行升级, 从1.7升级到1.9. 

第一种方法, 直接set镜像: 

```
$ kubectl set image deployment/nginx-deployment2 nginx=nginx:1.9
deployment "nginx-deployment2" image updated
```

第二种方法, 直接edit: 

```
$ kubectl edit deployment/nginx-deployment
deployment "nginx-deployment2" edited
```

查看Deployment的变更信息（以下信息得以保存, 是创建时候加的“--record”这个选项起的作用）: 

```
$ kubectl rollout history deployment/nginx-deployment
deployments "nginx-deployment":
REVISION    CHANGE-CAUSE
          kubectl create -f docs/user-guide/nginx-deployment.yaml --record
          kubectl set image deployment/nginx-deployment nginx=nginx:1.9.1
          kubectl set image deployment/nginx-deployment nginx=nginx:1.91

$ kubectl rollout history deployment/nginx-deployment --revision=2
deployments "nginx-deployment" revision 2
  Labels:       app=nginx
          pod-template-hash=1159050644
  Annotations:  kubernetes.io/change-cause=kubectl set image deployment/nginx-deployment nginx=nginx:1.9.1
  Containers:
   nginx:
    Image:      nginx:1.9.1
    Port:       80/TCP
     QoS Tier:
        cpu:      BestEffort
        memory:   BestEffort
    Environment Variables:      <none>
  No volumes.
```

最后介绍下Deployment的一些基础命令. 

```
$ kubectl describe deployments  #查询详细信息, 获取升级进度
$ kubectl rollout pause deployment/nginx-deployment2  #暂停升级
$ kubectl rollout resume deployment/nginx-deployment2  #继续升级
$ kubectl rollout undo deployment/nginx-deployment2  #升级回滚
$ kubectl scale deployment nginx-deployment --replicas 10  #弹性伸缩Pod数量
```

**关于多重升级**

举例, 当你创建了一个`nginx1.7`的Deployment, 要求副本数量为5之后, `Deployment Controller`会逐步的将5个1.7的Pod启动起来；当启动到3个的时候, 你又发出更新`Deployment`中Nginx到1.9的命令；这时`Deployment Controller`会立即将已启动的3个1.7Pod杀掉, 然后逐步启动1.9的Pod. Deployment Controller不会等到1.7的Pod都启动完成之后, 再依次杀掉1.7, 启动1.9. 

## 6. Volume

在Docker的设计实现中, 容器中的数据是临时的, 即当容器被销毁时, 其中的数据将会丢失. 如果需要持久化数据, 需要使用Docker数据卷挂载宿主机上的文件或者目录到容器中. 在Kubernetes中, 当Pod重建的时候, 数据是会丢失的, Kubernetes也是通过数据卷挂载来提供Pod数据的持久化的. Kubernetes数据卷是对Docker数据卷的扩展, Kubernetes数据卷是Pod级别的, 可以用来实现Pod中容器的文件共享. 目前, Kubernetes支持的数据卷类型如下: 

1. EmptyDir
2. HostPath
3. GCE Persistent Disk
4. AWS Elastic Block Store
5. NFS
6. iSCSI
7. Flocker
8. GlusterFS
9. RBD
10. Git Repo
11. Secret
12. Persistent Volume Claim
13. Downward API

### 6.1本地数据卷

`EmptyDir`、`HostPath`这两种类型的数据卷, 只能最用于本地文件系统. 本地数据卷中的数据只会存在于一台机器上, 所以当Pod发生迁移的时候, 数据便会丢失. 该类型Volume的用途是: Pod中容器间的文件共享、共享宿主机的文件系统. 

#### 6.1.1 EmptyDir

如果Pod配置了EmpyDir数据卷, 在Pod的生命周期内都会存在, 当Pod被分配到 Node上的时候, 会在Node上创建EmptyDir数据卷, 并挂载到Pod的容器中. 只要Pod 存在, EmpyDir数据卷都会存在（容器删除不会导致EmpyDir数据卷丟失数据）, 但是如果Pod的生命周期终结（Pod被删除）, EmpyDir数据卷也会被删除, 并且永久丢失. 

EmpyDir数据卷非常适合实现Pod中容器的文件共享. Pod的设计提供了一个很好的容器组合的模型, 容器之间各司其职, 通过共享文件目录来完成交互, 比如可以通过一个专职日志收集容器, 在每个Pod中和业务容器中进行组合, 来完成日志的收集和汇总. 

#### 6.1.2 HostPath

HostPath数据卷允许将容器宿主机上的文件系统挂载到Pod中. 如果Pod需要使用宿主机上的某些文件, 可以使用HostPath. 

### 6.2网络数据卷

Kubernetes提供了很多类型的数据卷以集成第三方的存储系统, 包括一些非常流行的分布式文件系统, 也有在IaaS平台上提供的存储支持, 这些存储系统都是分布式的, 通过网络共享文件系统, 因此我们称这一类数据卷为网络数据卷. 

网络数据卷能够满足数据的持久化需求, Pod通过配置使用网络数据卷, 每次Pod创建的时候都会将存储系统的远端文件目录挂载到容器中, 数据卷中的数据将被水久保存, 即使Pod被删除, 只是除去挂载数据卷, 数据卷中的数据仍然保存在存储系统中, 且当新的Pod被创建的时候, 仍是挂载同样的数据卷. 网络数据卷包含以下几种: NFS、iSCISI、GlusterFS、RBD（Ceph Block Device）、Flocker、AWS Elastic Block Store、GCE Persistent Disk.

### 6.3 Persistent Volume和Persistent Volume Claim

理解每个存储系统是一件复杂的事情, 特别是对于普通用户来说, 有时候并不需要关心各种存储实现, 只希望能够安全可靠地存储数据. Kubernetes中提供了`Persistent Volume`和`Persistent Volume Claim`机制, 这是存储消费模式. 

`Persistent Volume`是由系统管理员配置创建的一个数据卷（目前支持HostPath、GCE Persistent Disk、AWS Elastic Block Store、NFS、iSCSI、GlusterFS、RBD）, 它代表了某一类存储插件实现；

而对于普通用户来说, 通过Persistent Volume Claim可请求并获得合适的Persistent Volume, 而无须感知后端的存储实现. 

Persistent Volume和Persistent Volume Claim的关系其实类似于Pod和Node, Pod消费Node资源, Persistent Volume Claim则消费Persistent Volume资源. Persistent Volume和Persistent Volume Claim相互关联, 有着完整的生命周期管理: 

1. 准备: 系统管理员规划或创建一批Persistent Volume；

2. 绑定: 用户通过创建Persistent Volume Claim来声明存储请求, Kubernetes发现有存储请求的时候, 就去查找符合条件的Persistent Volume（最小满足策略）. 找到合适的就绑定上, 找不到就一直处于等待状态；

3. 使用: 创建Pod的时候使用Persistent Volume Claim；

4. 释放: 当用户删除绑定在Persistent Volume上的Persistent Volume Claim时, Persistent Volume进入释放状态, 此时Persistent Volume中还残留着上一个Persistent Volume Claim的数据, 状态还不可用；

5. 回收: 是否的Persistent Volume需要回收才能再次使用. 回收策略可以是人工的也可以是Kubernetes自动进行清理（仅支持NFS和HostPath）

### 6.4信息数据卷

Kubernetes中有一些数据卷, 主要用来给容器传递配置信息, 我们称之为信息数据卷, 比如Secret（处理敏感配置信息, 密码、Token等）、Downward API（通过环境变量的方式告诉容器Pod的信息）、Git Repo（将Git仓库下载到Pod中）, 都是将Pod的信息以文件形式保存, 然后以数据卷方式挂载到容器中, 容器通过读取文件获取相应的信息. 

## 7. Pet Sets/StatefulSet

K8s在1.3版本里发布了Alpha版的`PetSet`功能. 在云原生应用的体系里, 有下面两组近义词；

1. 无状态（stateless）、牲畜（cattle）、无名（nameless）、可丢弃（disposable）；

2. 有状态（stateful）、宠物（pet）、有名（having name）、不可丢弃（non-disposable）

RC和RS主要是控制提供无状态服务的, 其所控制的Pod的名字是随机设置的, 一个Pod出故障了就被丢弃掉, 在另一个地方重启一个新的Pod, 名字变了、名字和启动在哪儿都不重要, 重要的只是Pod总数；而PetSet是用来控制有状态服务, PetSet中的每个Pod的名字都是事先确定的, 不能更改. PetSet中Pod的名字的作用, 是用来关联与该Pod对应的状态. 

对于RC和RS中的Pod, 一般不挂载存储或者挂载共享存储, 保存的是所有Pod共享的状态, Pod像牲畜一样没有分别；对于PetSet中的Pod, 每个Pod挂载自己独立的存储, 如果一个Pod出现故障, 从其他节点启动一个同样名字的Pod, 要挂在上原来Pod的存储继续以它的状态提供服务. 

适合于PetSet的业务包括数据库服务MySQL和PostgreSQL, 集群化管理服务Zookeeper、etcd等有状态服务. PetSet的另一种典型应用场景是作为一种比普通容器更稳定可靠的模拟虚拟机的机制. 传统的虚拟机正是一种有状态的宠物, 运维人员需要不断地维护它, 容器刚开始流行时, 我们用容器来模拟虚拟机使用, 所有状态都保存在容器里, 而这已被证明是非常不安全、不可靠的. 使用PetSet, Pod仍然可以通过漂移到不同节点提供高可用, 而存储也可以通过外挂的存储来提供高可靠性, PetSet做的只是将确定的Pod与确定的存储关联起来保证状态的连续性. 

## 8. ConfigMap

很多生产环境中的应用程序配置较为复杂, 可能需要多个config文件、命令行参数和环境变量的组合. 并且, 这些配置信息应该从应用程序镜像中解耦出来, 以保证镜像的可移植性以及配置信息不被泄露. 社区引入ConfigMap这个API资源来满足这一需求. 

ConfigMap包含了一系列的键值对, 用于存储被Pod或者系统组件（如controller）访问的信息. 这与secret的设计理念有异曲同工之妙, 它们的主要区别在于ConfigMap通常不用于存储敏感信息, 而只存储简单的文本信息. 

## 9. Horizontal Pod Autoscaler

自动扩展作为一个长久的议题, 一直为人们津津乐道. 系统能够根据负载的变化对计算资源的分配进行自动的扩增或者收缩, 无疑是一个非常吸引人的特征, 它能够最大可能地减少费用或者其他代价（如电力损耗）. 自动扩展主要分为两种, 其一为水平扩展, 针对于实例数目的增减；其二为垂直扩展, 即单个实例可以使用的资源的增减. Horizontal Pod Autoscaler（HPA）属于前者. 

### 9.1 Horizontal Pod Autoscaler如何工作

Horizontal Pod Autoscaler的操作对象是Replication Controller、ReplicaSet或Deployment对应的Pod, 根据观察到的CPU实际使用量与用户的期望值进行比对, 做出是否需要增减实例数量的决策. controller目前使用heapSter来检测CPU使用量, 检测周期默认是30秒. 

### 9.2 Horizontal Pod Autoscaler的决策策略

在HPA Controller检测到CPU的实际使用量之后, 会求出当前的CPU使用率（实际使用量与pod 请求量的比率). 然后, HPA Controller会通过调整副本数量使得CPU使用率尽量向期望值靠近．另外, 考虑到自动扩展的决策可能需要一段时间才会生效, 甚至在短时间内会引入一些噪声． 例如当pod所需要的CPU负荷过大, 从而运行一个新的pod进行分流, 在创建的过程中, 系统的CPU使用量可能会有一个攀升的过程. 所以, 在每一次作出决策后的一段时间内, 将不再进行扩展决策. 对于ScaleUp而言, 这个时间段为3分钟, Scaledown为5分钟. 再者HPA Controller允许一定范围内的CPU使用量的不稳定, 也就是说, 只有当aVg（CurrentPodConsumption／Target低于0.9或者高于1.1时才进行实例调整, 这也是出于维护系统稳定性的考虑. 