# calico-bgp全互联(Node-to-Node)与路由反射(RR).2

参考文章

1. [Kubernetes网络组件之Calico策略实践(BGP、RR、IPIP)](https://blog.51cto.com/14143894/2463392)
    - 第6节: Route Reflector 模式（RR）（路由反射）
2. [Calico配置双RR架构](https://www.cnblogs.com/delacroix429/p/11718491.html)

环境准备

- km01: 192.168.80.121
- kw01: 192.168.80.124
- kw02: 192.168.80.125

## 1. 关闭 node-to-node 连接

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
    route-reflector: true
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
  ## 所有的节点
  nodeSelector: all()
  peerSelector: route-reflector == 'true' 
```

应用这个配置, 将会发生如下变化:

1. `calicoctl node status`查看到的输出中, `Peer Type`的类型将为`node specific`
2. 用`netstat`查看, km01 `192.168.80.121`节点上将出现与其他节点的连接, 其他节点之间只连接 km01, 不再两两相互连接.
3. 此时 Pod 与 Pod 之间, Pod 与集群中其余宿主机节点间便可相互通信.

```console
$ calicoctl node status
Calico process is running.

IPv4 BGP status
+----------------+---------------+-------+----------+-------------+
|  PEER ADDRESS  |   PEER TYPE   | STATE |  SINCE   |    INFO     |
+----------------+---------------+-------+----------+-------------+
| 192.168.80.121 | node specific | up    | 11:11:34 | Established |
+----------------+---------------+-------+----------+-------------+

IPv6 BGP status
No IPv6 peers found.
```

## 

...这个实验仍是没有完成, 还是失败, 部分节点上的 Pod 无法与其他节点上的 Pod 通信. ???

