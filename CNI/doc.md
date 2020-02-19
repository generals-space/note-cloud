# 网络插件

参考文章

1. [Kubernetes: Flannel networking](https://blog.laputa.io/kubernetes-flannel-networking-6a1cb1f8ec7c)
    - kubernetes网络模型
    - flannel插件实现的overlay网络
    - flannel组件以UDP模式实现的, 容器间跨主机通信时数据包传递的全过程(精彩)
    - 生产环境不推荐使用UDP模式的原因: 性能
2. [Kubernetes网络方案的三大类别和六个场景](https://sq.163yun.com/blog/article/223878660638527488)
    - 总纲级别, 值得一读
3. [官方文档 Network Plugins](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/)

4. [Kubernetes CNI网络插件](https://www.cnblogs.com/rexcheny/p/10960233.html)

5. [Kubernetes利用CNI-bridge插件打通网络](https://blog.csdn.net/qq_36183935/article/details/90735049)
6. [Flannel是如何工作的](https://cloud.tencent.com/developer/article/1096997)
    - vxlan, hostgw, udp才是有真正使用场景的网络模型, 其他都是实验性的, 不建议上生产.
    - [containernetworking/plugins/plugins/meta/flannel/README.md](https://github.com/containernetworking/plugins/blob/master/plugins/meta/flannel/README.md)工程才是真正的cni插件.


参考文章1简洁易懂, ta首先描述了kuber设想中的网络架构: 

1. 容器与容器间(包括不同宿主机间的容器), 容器与宿主机节点间可相互通信, 且不是通过3层NAT的方式, 不能经过地址转换. 
2. 容器A看到的自身的地址, 与其他容器看到的A的地址是同一个. 这其实也是NAT的概念.

kuber只是定义了ta需要的网络模型, 但ta本身并没有实现, 而是把这个任务交给了网络插件. 什么意思呢? 

kuber其实只是一个容器管理系统, 以docker容器为例, docker默认不支持跨主机的通信(其实现在已经支持了, 但不符合ONI标准, kuber采用ONI标准规范其网络插件, 以适用多种container runtime). kuber只是代为控制各节点上容器的运行(还有维护, 监控等功能), ta同时希望不同宿主机节点间的pod可以像上面定义的那样相互通信, 要怎么做呢? 

kuber引入flannel, flannel通过在每个节点(包括Master和Worker)启动一个服务`flanneld`, 并创建TUN设备, 像路由器一样将容器发送出的数据包转发给其他主机, 再由接收方主机转发给本机上的其他容器, 具体的细节可以查看参考文章1.

------

这就是kuber中网络插件的角色, 不同的网络插件拥有不同的特性, 表现在不同的方面. 

比如参考文章1中介绍的flannel的UDP模式, 架构简单, 又容易理解, 但是由于数据包在用户空间与内核空间多次拷贝, 因此性能较差. 

另外仔细想想, 如果虚拟环境需要实现多租户的功能, 类似vlan的网络隔离, 该怎么办呢? 又或者, pod的IP未知且不固定, 如果希望通过dhcp进行IP分配(分配后固定, 且可以分配指定地址), 又该怎么办呢?

这就是不同网络插件自由发挥的地方了, ta们定义的虚拟网络架构各不相同, 各擅胜场, 需要运维人员仔细比较后选型. 参考文章2中有一张图片提到了多种网络插件及ta们各自的特性, 可以了解一下.

------

尝试部署了kuber集群后发现若无网络插件, 会影响域名解析, 跨宿主机的容器间通信, 以及pod创建时的IP分配(pod可能永远也无法启动)
