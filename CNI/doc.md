# 网络插件

参考文章

1. [Kubernetes: Flannel networking](https://blog.laputa.io/kubernetes-flannel-networking-6a1cb1f8ec7c)
    - kubernetes网络模型
    - flannel插件实现的overlay网络
    - flannel组件以UDP模式实现的, 容器间跨主机通信时数据包传递的全过程(精彩)
    - 生产环境不推荐使用UDP模式的原因: 性能(数据包在用户空间与内核空间多次拷贝)
2. [Kubernetes利用CNI-bridge插件打通网络](https://blog.csdn.net/qq_36183935/article/details/90735049)
    - 原生CNI插件`bridge`的工作流程
    - 提出"将物理网卡连接到网桥中, 并给bridge设备设置IP"的想法.

参考文章1简洁易懂, ta首先描述了kuber设想中的网络架构: 

1. 容器与容器间(包括不同宿主机间的容器), 容器与宿主机节点间可相互通信, 且不是通过3层NAT的方式, 不能经过地址转换. 
2. 容器A看到的自身的地址, 与其他容器看到的A的地址是同一个. 这其实也是NAT的概念.

kuber只是定义了ta需要的网络模型, 但ta本身并没有实现, 而是把这个任务交给了网络插件. 什么意思呢? 

kuber其实只是一个容器管理系统, 以docker容器为例, docker默认不支持跨主机的通信(其实现在已经支持了, 但不符合ONI标准, kuber采用ONI标准规范其网络插件, 以适用多种container runtime). kuber只是代为控制各节点上容器的运行(还有维护, 监控等功能), ta同时希望不同宿主机节点间的pod可以像上面定义的那样相互通信, 要怎么做呢? 

kuber引入flannel, flannel通过在每个节点(包括Master和Worker)启动一个服务`flanneld`, 并创建TUN设备, 像路由器一样将容器发送出的数据包转发给其他主机, 再由接收方主机转发给本机上的其他容器, 具体的细节可以查看参考文章1.
