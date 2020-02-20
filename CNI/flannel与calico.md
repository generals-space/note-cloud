参考文章

1. [深入浅出Kubernetes网络：跨节点网络通信之Flannel](https://cloud.tencent.com/developer/article/1450296)
    - 详细讲解了`flannel`的3种网络模型: `UDP`, `VxLAN`, `host-gw`, 以及ta们的工作机制和缺陷. 
    - 配图, 十分详细, 值得阅读.
    - `UDP`模型: 创建`tun`设备, 数据包经典`tun`设备封包再通过`eth0`发出.
    - `VxLAN`: 
    - `host-gw`模型: `flannel`的唯一作用就是负责主机上路由表的动态更新, 不过缺陷是需要宿主机处于同一子网, 否则路由无法直达. 
2. [Kubernetes学习之路（二十一）之网络模型和网络策略](https://www.cnblogs.com/linuxk/p/10517055.html)
3. [Kubernetes网络插件对比分析（Flannel、Calico、Weave）](https://network.51cto.com/art/201907/598970.htm)
4. [Kubernetes: Flannel networking](https://blog.laputa.io/kubernetes-flannel-networking-6a1cb1f8ec7c)
    - kubernetes网络模型
    - flannel插件实现的overlay网络
    - flannel组件以UDP模式实现的, 容器间跨主机通信时数据包传递的全过程(精彩)
    - 生产环境不推荐使用UDP模式的原因: 性能(数据包在用户空间与内核空间多次拷贝)
3. [Kubernetes网络方案的三大类别和六个场景](https://sq.163yun.com/blog/article/223878660638527488)
    - 总纲级别, 值得一读
    - Calico提供了两种网络模型: 
        1. BGP (`Underlay L3`)
        2. ipip (`Overlay L3`)

