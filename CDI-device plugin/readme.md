参考文章

资源调度

1. [设备插件](https://kubernetes.io/zh-cn/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/)
2. [【k8s调度】梳理调度相关知识与device plugin](https://blog.csdn.net/qq_24433609/article/details/137020551)

设备注册

1. [基于K8s的SR-IOV网络实践](https://cloud.tencent.com/developer/article/2010030)
    - 对于网络性能有一定要求，但要求不是很高的，例如低配版数据库实例，可以使用macvlan或者ipvlan的网卡.
        - 看来sriov网卡在性能上要比macvlan更强
2. [SR-IOV——网卡直通技术](https://www.cnblogs.com/weiduoduo/p/11068460.html)
3. [SR-IOV研究：一个简单的测试环境](https://www.sensorexpert.com.cn/article/215228.html)
4. [device-plugin 扩展——intel-sriov-device-plugin解读](https://segmentfault.com/a/1190000021061494?sort=votes)

5. [深入剖析 Kubernetes 44 | Kubernetes GPU管理与Device Plugin机制](https://blog.csdn.net/qq_37756660/article/details/134523326)
    - Kubernetes 里对硬件设备的管理，只能处理“设备个数”这唯一一种情况。一旦你的设备是异构的、不能简单地用“数目”去描述具体使用需求的时候，比如，“我的 Pod 想要运行在计算能力最强的那个 GPU 上”，Device Plugin 就完全不能处理了。
6. [深入剖析 Kubernetes 44 | Kubernetes GPU管理与Device Plugin机制](https://time.geekbang.org/column/article/70876)
    - 参考文章5的原文, 评论很有意思

7. [基于KVM的SRIOV直通配置及性能测试](https://cloud.tencent.com/developer/article/2060167)
8. [GPU直通、GPU共享及GPU虚拟化比较](https://blog.csdn.net/hbuxiaofei/article/details/106757002)

9. [盘点来自工业界的 GPU 共享方案](https://www.cvmart.net/community/detail/5370)
    - 来源: https://zhuanlan.zhihu.com/p/398369404
10. [GPUManager虚拟化方案](https://cloud.tencent.com/developer/article/1685122)
11. [腾讯云徐蓓：通过云原生管理Kubernetes GPU资源](https://mp.weixin.qq.com/s/W-Ntu2xdjypFgs5EPMVjkg)
12. [论文笔记《GaiaGPU：Sharing GPUs in Container Clouds》](https://pokerfacesad.github.io/2020/02/07/%E8%AE%BA%E6%96%87%E7%AC%94%E8%AE%B0%E3%80%8AGaiaGPU%20Sharing%20GPUs%20in%20Container%20Clouds%E3%80%8B/)

