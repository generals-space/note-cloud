参考文章

1. [深入浅出Kubernetes网络：跨节点网络通信之Flannel](https://cloud.tencent.com/developer/article/1450296)
    - 详细讲解了`flannel`的3种网络模型: `UDP`, `VxLAN`, `host-gw`, 以及ta们的工作机制和缺陷. 配图, 十分详细.
    - `host-gw`模式`flannel`的唯一作用就是负责主机上路由表的动态更新, 不过缺陷是需要宿主机处于同一子网, 否则路由无法直达. 
2. [Kubernetes学习之路（二十一）之网络模型和网络策略](https://www.cnblogs.com/linuxk/p/10517055.html)
3. [Kubernetes网络插件对比分析（Flannel、Calico、Weave）](https://network.51cto.com/art/201907/598970.htm)

