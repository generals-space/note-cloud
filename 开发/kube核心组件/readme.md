参考文章

[为什么Kubernetes不使用Libnetwork？](http://dockone.io/article/973)
    - 写的不错, 论据很充分
[kubernetes pod为什么需要pause容器](https://blog.csdn.net/lywzgzl/article/details/100625594)
    - 不错
[跟我学 K8S--代码: 调试 Kubernetes](https://segmentfault.com/a/1190000015807702)
    - 脚本官方脚本`local-up-cluster.sh`启动集群
    - `dlv`进行调试
[k8s 源码阅读系统文章](https://chenxy.blog.csdn.net/category_6361923.html)

[kubernetes源码剖析读后感（三）](https://blog.csdn.net/weixin_45413603/article/details/107170116)
[sigma敏捷版系列文章：kubernetes grace period 失效问题排查](https://developer.aliyun.com/article/609813)
[kube-apiserver 的设计与实现](https://www.jianshu.com/p/7100880a8858)

[kube-scheduler源码分析（3）-抢占调度分析](https://blog.csdn.net/kyle18826138721/article/details/123306546)

- `etcd`: 保存了整个集群的状态;
- `apiserver`: 提供了资源操作的唯一入口，并提供认证、授权、访问控制、API 注册和发现等机制;
- `controller-manager`: 负责维护集群的状态，比如故障检测、自动扩展、滚动更新等;
- `scheduler`: 负责资源的调度，按照预定的调度策略将 Pod 调度到相应的机器上;
- `kubelet`: 负责维护容器的生命周期，同时也负责 Volume(CVI)和网络(CNI)的管理;
- `container runtime`: 负责镜像管理以及 Pod 和容器的真正运行(CRI);
- `kube-proxy`: 负责为 Service 提供 cluster 内部的服务发现和负载均衡;
