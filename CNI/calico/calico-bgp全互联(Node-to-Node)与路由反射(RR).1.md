# calico-bgp全互联(Node-to-Node)与路由反射(RR).1

参考文章

1. [k8s网络之calico学习](https://www.jianshu.com/p/106eb0a09765)
    - `calicoctl`查看集群中当前的连接模式
    - 将默认的 node-to-node 模式修改为 BGP Speaker RR模式
    - 本文的实验中使用docker运行`calico/routereflector:v0.6.1`镜像, 不够通用, 且未完成, **不适合做参考**.
2. [calico/routereflector](https://hub.docker.com/r/calico/routereflector)
    - 参考文章1应该就是借鉴的这个.
    - 这里写的`calicoctl bgp node-mesh off`关闭`node-to-node`网格, 但是`calicoctl`根本没有`bgp`这个子命令...差评
3. [Kubernetes网络组件之Calico策略实践(BGP、RR、IPIP)](https://blog.51cto.com/14143894/2463392)
    - 第6节: Route Reflector 模式（RR）（路由反射）
4. [Calico配置双RR架构](https://www.cnblogs.com/delacroix429/p/11718491.html)

## 0. 写在前面

本文还是作为一个实验性的试手, 正式的 rr 配置请见下一篇文章.

## 1. 查看

- km01: 192.168.80.121
- kw01: 192.168.80.124
- kw02: 192.168.80.125

使用`calicoctl`查看当前集群中使用的连接模式. 如下命令在`calicoctl`所在Pod内部执行, 该Pod运行在km01.

```console
$ calicoctl get node -o wide
NAME            ASN       IPV4                IPV6   
k8s-master-01   (64512)   192.168.80.121/24          
k8s-worker-01   (64512)   192.168.80.124/24          
k8s-worker-02   (64512)   192.168.80.125/24     
```

```console
$ calicoctl node status
Calico process is running.

IPv4 BGP status
+----------------+-------------------+-------+----------+-------------+
|  PEER ADDRESS  |     PEER TYPE     | STATE |  SINCE   |    INFO     |
+----------------+-------------------+-------+----------+-------------+
| 192.168.80.124 | node-to-node mesh | up    | 11:07:00 | Established |
| 192.168.80.125 | node-to-node mesh | up    | 11:06:58 | Established |
+----------------+-------------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.
```

默认calico使用`node-to-node`的全互联模式, 节点之间两两都会建立连接, 以进行路由交换. 适合规模不大的集群中运行, 一旦集群节点增大, mesh模式将形成一个巨大服务网格, 连接数暴增.

如果我们在 km01 宿主机上使用`netstat`查看与`bird`建立的连接, 有如下输出.

```console
$ netstat -nap | grep bird
tcp        0      0 0.0.0.0:179             0.0.0.0:*               LISTEN      31891/bird          
tcp        0      0 192.168.80.121:179      192.168.80.124:60680    ESTABLISHED 31891/bird          
tcp        0      0 192.168.80.121:179      192.168.80.125:36340    ESTABLISHED 31891/bird 
```

`km01`与`kw01`/`kw02`都建立了连接, 其他节点上也是如此.

## 2. 配置 BGP Speaker RR 模式

### 2.1 关闭 node-to-node 连接

首先创建`BGPConfiguration`资源, 这个资源在`02-crds.yaml`中已经声明过.

```yaml
apiVersion: projectcalico.org/v3
kind: BGPConfiguration
metadata:
  name: default
spec:
  logSeverityScreen: Info
  ## node-to-node 模式关闭
  nodeToNodeMeshEnabled: false
  asNumber: 61234
```

使用`calicoctl apply -f`应用此配置(类似于`kubectl apply -f`), 将会发生如下变化:

1. `calicoctl node status`之前能查看到的`node-to-node`连接将全部消失.
2. `netstat -anp | grep bird`之前与其他节点已经建立的连接也将消失.
3. 此时各 Pod 之间, Pod 与各宿主机节点之间的网络将会断开.

```console
$ calicoctl node status
Calico process is running.

IPv4 BGP status
No IPv4 peers found.

IPv6 BGP status
No IPv6 peers found.
```

### 2.2 选取 rr 节点

这里我们选取 km01 节点作为 rr 节点. 首先使用`calicoctl get node k8s-master-01 -o yaml`输出该节点的信息, 然后我们做一些修改, 修改的部分用注释写明.

```yaml
apiVersion: projectcalico.org/v3
kind: Node
metadata:
  annotations:
    projectcalico.org/kube-labels: '{"beta.kubernetes.io/arch":"amd64","beta.kubernetes.io/os":"linux","kubernetes.io/arch":"amd64","kubernetes.io/hostname":"k8s-master-01","kubernetes.io/os":"linux","node-role.kubernetes.io/master":"","route-reflector":"true"}'
  labels:
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/os: linux
    kubernetes.io/arch: amd64
    kubernetes.io/hostname: k8s-master-01
    kubernetes.io/os: linux
    node-role.kubernetes.io/master: ""
  name: k8s-master-01
  uid: dbf7cf79-d3f0-4950-b6f2-93f37d74d047
spec:
  bgp:
    ipv4Address: 192.168.80.121/24
    ## 这是 rr 节点的 id, 当集群中要部署多个 rr 节点时, 这些节点的 id 应该保持一致.
    routeReflectorClusterID: 224.0.0.1
  ## 这个字段貌似没什么用, 不过参考文章3,4都有写到.
  orchRefs:
  - nodeName: k8s-master-01
    orchestrator: k8s
```

同样使用`calicoctl apply -f`应用此配置, 但网络基本没有什么变化.

### 2.3 为 rr 节点与普通 bgp 节点间建立连接

接下来, 我们需要将普通的 bgp 节点与 rr 节点建立连接. 

```yaml
apiVersion: projectcalico.org/v3
kind: BGPPeer
metadata:
  name: bgppeer-global
spec:
  peerIP: 192.168.80.121
  ## 这里与上面 BGPConfiguration 中的`asNumber`字段保持一致.
  asNumber: 61234
```

应用这个配置, 将会发生如下变化:

1. `calicoctl node status`查看到的输出中, `Peer Type`的类型将为`node specific`
2. 用`netstat`查看, km01 `192.168.80.121`节点上将出现与其他节点的连接, 其他节点之间只连接 km01, 不再两两相互连接.

```console
$ calicoctl node status
Calico process is running.

IPv4 BGP status
+----------------+---------------+-------+----------+-------------+
|  PEER ADDRESS  |   PEER TYPE   | STATE |  SINCE   |    INFO     |
+----------------+---------------+-------+----------+-------------+
| 192.168.80.121 | node specific | up    | 09:30:21 | Established |
+----------------+---------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.
```

还没完呢...总觉得哪里怪怪的...

## 3. 后续疑问


### 3.1 路由不完整

首先是, 各主机节点上的路由并不完整, km01 上是完整的, 但是 kw01 和 kw02 上分别只有到对方子网的路由, 没有到 km01 网段的路由. 

km01 的路由如下

```
default via 192.168.80.2 dev ens33 proto static metric 100 
172.16.36.192/26 via 192.168.80.124 dev ens33 proto bird 
172.16.118.64/26 via 192.168.80.125 dev ens33 proto bird 
```

kw01 的路由如下

```
default via 192.168.80.2 dev ens33 proto static metric 100 
blackhole 172.16.36.192/26 proto bird 
172.16.118.64/26 via 192.168.80.125 dev ens33 proto bird
```

kw02 的路由如下

```
default via 192.168.80.2 dev ens33 proto static metric 100 
172.16.36.192/26 via 192.168.80.124 dev ens33 proto bird 
blackhole 172.16.118.64/26 proto bird 
```

km01 上的 Pod 可以 ping 通所有其他的 Pod, 但是 kw01 和 kw02 上的 Pod 就只能 ping 通本机上的 Pod 了, 不过 kw01 和 kw02 在宿主机上还是可以 ping 通所有 Pod的...

### 3.2 peer 表示什么

我在 km01 上用 netstat 查了一下, 有如下输出

```console
$ netstat -anp | grep bird
tcp        0      0 0.0.0.0:179             0.0.0.0:*               LISTEN      31891/bird          
tcp        0      0 192.168.80.121:41558    192.168.80.125:179      ESTABLISHED 31891/bird          
tcp        0      0 192.168.80.121:51614    192.168.80.124:179      ESTABLISHED 31891/bird 
```

你会发现, `:179`端口的才像是服务端, 那就存在了`kw01`和`kw02`两个服务端, 而`km01`作为客户端, 有两个连接...

