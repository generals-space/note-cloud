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

