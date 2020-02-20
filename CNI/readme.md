参考文章

1. [Kubernetes指南 - 网络模型](https://feisky.gitbooks.io/kubernetes/network/network.html)
    - Kubernetes 支持的两种插件: kubenet 和 CNI
    - 列举各种第三方插件: Flannel, Weave, Calico, OVN等, 介绍了ta们各自的功能及优缺点.
    - 还有一些其他的kubernetes网络机制: ipvs, Canel插件, kube-router等(都不重要).
    - `Canal`是`Flannel`和`Calico`联合发布的一个统一网络插件, 提供`CNI`网络插件, 并支持`network policy`.
2. [浅谈k8s cni 插件](https://segmentfault.com/a/1190000017182169)
    - [containernetworking/plugins](https://github.com/containernetworking/plugins)工程
    - 介绍了使用CNI插件需要做的配置: 
        1. `kubelet`启动参数`--network-plugin=cni`
        2. `/etc/cni/net.d/`目录下增加CNI插件的配置文件
        3. `/opt/cni/bin`目录下存在CNI可执行文件
    - 介绍了`plugins`工程中的各种插件, 按照其功能分为3类: `main`, `ipam`和`meta`, 并分别介绍ta们工作的具体流程.
3. [Flannel是如何工作的](https://cloud.tencent.com/developer/article/1096997)
    - vxlan, hostgw, udp才是有真正使用场景的网络模型, 其他都是实验性的, 不建议上生产.
    - [containernetworking/plugins/plugins/meta/flannel/README.md](https://github.com/containernetworking/plugins/blob/master/plugins/meta/flannel/README.md)工程才是真正的cni插件.
4. [K8S 网络插件（CNI）超过 10Gbit/s 的基准测试结果](https://zhuanlan.zhihu.com/p/53296042)
    - CNI性能测试结果, 分析和建议...不过好像没说方法?

CNI: 容器网络接口.

官方的CNI插件在[containernetworking/plugins](https://github.com/containernetworking/plugins)工程, 通过`yum install kubernetes-cni`, 安装在`/opt/cni/bin`目录下, ta们的配置文件存放在`/etc/cni/net.d/`目录.

## 1. CNI的工作流程

一般来说, CNI插件由`kubelet`组件调用, 当然, 这需要在setup集群时为`kubelet`指定`--network-plugin=cni`(如果你使用的是`kubeadm`工具, 这应该是个默认值).

`kubelet`在接收到调度器新建Pod的指令后, 先创建`pause`容器, 完成后将**容器ID**, **netns路径**等信息, 再加上`/etc/cni/net.d/`目录下的配置(这里称为`netconf`)当作参数传给指定的CNI插件, 由CNI插件创建`veth pair`连接宿主机与Pod, 同时按照`netconf`中的`ipam`配置为此Pod赋予IP. 完成后再由`kubelet`继续创建业务Pod, 共享此`pause`容器的网络空间.

连接宿主机与Pod的方式基本上全都靠`bridge`设备, 就如同`docker0`, 在CNI中, 默认为`cni0`.

而为Pod赋予IP, 目前官方提供了3种方式: 

1. static: 需要在CNI配置中添加静态IP字符串, 这个基本不会用到.
2. host-local: 在CNI配置中写入指定网段, 由`host-local`插件从中选择合适的IP.
3. dhcp: dhcp插件其实是一个客户端, ta会向集群网络中发布广播请求, 并返回申请到的IP.

更详细的内容可以见参考文章2.

> 可以说, 官方提供的CNI所做的事情, 完全可以通过`ip`, `brctl`模拟出来.

## 2. flannel的工作流程

我们熟知的Kubernetes的网络插件有`flannel`, `calico`等, 但ta们其实不算CNI, 因为ta们没有实现CNI接口, 不过ta们都基于CNI.

以`flannel`为例, 其实`flannel`分为两个工程:

1. [coreos/flannel](https://github.com/generals-space/flannel) (下称`coreos/flannel`)
2. [containernetworking/plugins](https://github.com/containernetworking/plugins)中的子项目, 位于`plugins/meta/flannel`. (下称`cni/flannel`)

后者才是真正的CNI插件(你为发现`/opt/cni/bin`目录下还有一个`flannel`), 见参考文章3.

按照上面给出的CNI工作流程, `cni/flannel`看起来只做了两件事: 连接宿主机与Pod, 为Pod赋予IP(当然不是像说得这么简单, 不过大概就是这个流程).

那么前者做了什么呢? 

首先, 集群在setup的时候, Pod的网段`pod cidr`就已经划分好了, 假设为`10.254.0.0/16`. `coreos/flannel`会监控node节点的变动, 每增加一个节点, 就为其划分一个子网, 如`10.254.1.0/24`, 之后在此node上创建的Pod就会在这个网段中分配IP. 当然, node移除时也需要归还该网段.

但这样只是为Pod分配了IP而已, 根本无法通信. 

```
+----------------------------------+                      +----------------------------------+
|  +----------------+              |                      |              +----------------+  |
|  | 10.254.1.1/24  |   Pod11      |                      |      Pod21   | 10.254.2.1/24  |  |
|  |   +--------+   |              |                      |              |   +--------+   |  |
|  |   |  eth0  |   |              |                      |              |   |  eth0  |   |  |
|  |   +----↑---+   |              |                      |              |   +----↑---+   |  |
|  +--------|-------+              |                      |              +--------|-------+  |
|    +------↓-----+                |                      |                +------↓-----+    |
|    |    cni0    |                |                      |                |    cni0    |    |
|    +------------+                |                      |                +------------+    |
|                      +--------+  |                      |  +--------+                      |
|    192.168.0.201/24  |  eht0  |  |                      |  |  eht0  |  192.168.0.202/24    |
|                      +----┬---+  |                      |  +----┬---+                      |
+---------------------------|------+                      +-------|--------------------------+
                            |        +------------------+         |                             
                            └───────>|  192.168.0.1/24  |<────────┘                             
                                     +------------------+
                                           网关/路由器
```

flannel提供了几个网络模型, 通过NAT, 路由等方式, 以实现Pod与Pod之间, Pod与宿主机之间能够通信.

这不是本文的重点, 不过多讲解.

