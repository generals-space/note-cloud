# calico

参考文章

1. [Calico网络模型](https://www.cnblogs.com/menkeyi/p/11364977.html)
    - calico网络模型, 讲解很详细
    - 容器跨主机通信的转发流程
    - felix: 路由配置组件; bird: 路由广播组件(BGP Speaker)
    - ipip网络模型的目标.
2. [Calico网络方案](https://www.cnblogs.com/netonline/p/9720279.html)
    - 文章末尾给出的calico与flannel各模型网络的对比很值得一看.
3. [calico网络策略](https://yq.aliyun.com/articles/674020)
    - `calico`不同于`flannel`, 不需要为每个node分配子网段, 所以只需要考虑pod的数量;
    - 宿主机上默认初始会创建一个掩码位为26的子网网段, 当pod超过这个值后, 可以从其他可用子网中取值, 并不是固定的.
4. [docker 容器网络方案：calico 网络模型](https://cizixs.com/2017/10/19/docker-calico-network/)

ipip网络模型是为了解决集群中各节点不在同一网段的问题. 一般情况下集群中的节点都在同一局域网, 但也有可能为了做冗余和灾备, 节点位于不同地点的机房, 此时应该是可以通过外网IP实现集群通信的. 但是BGP模型无法在这种类型的网络中通过三层路由完成, 所以ipip出现了.

按照参考文章2末尾所说, calico的ipip与flannel的vxlan/udp一样是overlay的解决方案. 另外也提到, calico的BGP与flannel的host-gw, 都要求集群各节点在相同子网网段中, 也在参考文章1中找到了答案.
