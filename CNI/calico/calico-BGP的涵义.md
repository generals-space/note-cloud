# calico-BGP的涵义

1. [什么是BGP服务器？BGP服务器有什么特点？](https://zhuanlan.zhihu.com/p/93153990)
    - BGP(一般都是多线)机房的实现方式及优势
2. [能够使服务器实现“多线路互联”和自动切换“最佳访问路由”的BGP路由协议！](https://www.yisu.com/news/id_414.html)
    - 对自治系统(AS)有一个清晰的解释, 可以参考
    - 介绍了AS间用于连接的BGP节点的位置, 建立连接的每个AS中必须存在一个BGP节点, 通常由"路由器"来执行BGP.
3. [华为交换机OSPF和BGP知识](https://blog.51cto.com/11009796/2126410)
    - OSPF属于IGP协议(内部网关协议), 一般运行在AS自治系统内部;
    - BGP属于EGP(外部网关协议), 一般是由ISP服务提供商运用在各个AS之间;
    - 多图易懂
4. [Calico网络模型](https://www.cnblogs.com/menkeyi/p/11364977.html)
    - 大规模部署架构图
5. [Calico官方 Why BGP?](https://www.projectcalico.org/why-bgp/)
    - working end points 应该是指网络终端, 比如常规的服务器, PC, 手机等网络设备.
    - AS自治系统: Autonomous System
    - OSPF与IS-IS是IGP的两种具体实现, 用于AS内部的路由设备发现网络中的其他设备
    - 中文译文 [为什么Calico网络选择BGP？](https://blog.51cto.com/weidawei/2152319)

自治系统(AS)可以理解为独立机房, 机房内部实现自治(即有权自主决定采用何种路由协议).

`BGP`本质上是一种3层协议, 机房在接入多线ISP后, 还需要专用设备将ta们集中连接并处理, 这种设备必须具备3层路由功能, 但应该也不是直接用路由器做的, 应该大部分采用的是3层交换机.

对应到calico网络, calico把每个宿主机当作网关路由器, 宿主机与其上的容器就形成了一个自治系统. 但还不够, flannel也有类似的`host-gw`模型. 最重要的, BGP提供的是多AS间的互联方式, 至少需要部署一个实现了BGP协议的设备. calico中由`bird`进程实现了`BGP Router Reflector`的功能(bird本身就是一个常用的路由软件, 支持多种路由协议), 具体的架构可以见参考文章4的"大规模部署架构图", 应该就比较容易理解了.

但我目前还不太明白BGP执行者的实际原理, 待探讨.
