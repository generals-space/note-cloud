# cgroup目录下kubepods.slice与Pod UID

参考文章

1. [解析容器化技术中的资源管理](http://www.dockone.io/article/9004)
2. [解析容器化技术中的资源管理](https://mp.weixin.qq.com/s/jT6m05vy601paKNi-Wh-6A)
    - 参考文章1的原文地址
3. [Kubernetes生产实践系列之三十：Kubernetes基础技术之集群计算资源管理](https://blog.csdn.net/cloudvtech/article/details/107634724)

`cgroup`目录的各子系统下(`/sys/fs/cgroup/{cpu,memory}/kubepods`), 都存在`kubepods`目录, 存放着被调度到当前主机的所有Pod的资源配置.

![](https://gitee.com/generals-space/gitimg/raw/master/5a6b7a80eeb7db228c49b4873ca67e36.png)

如果是`Guaranteed`级别的 Pod, 则会直接出现该目录下, 就像上面的`kubepods-podfa6c2e7a_35c6_42dd_a4ab_03641cfa4508.slice`. 其中文件名里的`kubepods-pod`后面的部分, 是该Pod的UID值.

![](https://gitee.com/generals-space/gitimg/raw/master/61680d6fc2744c2f67eccccfe0ae4304.png)

而`Burstable`和`BestEffort`级别的 Pod 则会放到对应的字目录下. 

在各自`kubepods-pod${PodUID}.slice`目录下, 则存在Pod中各个`container`的目录及资源配置, 如下

![](https://gitee.com/generals-space/gitimg/raw/master/c1d7027e0082a9f9c18bceac490806f3.png)

`docker-xxx.scope`中的`xxx`当然就是容器的id, 有2个, 其中一个是`pause`.
