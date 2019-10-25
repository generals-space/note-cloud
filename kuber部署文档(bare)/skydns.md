1. [官方文档](https://github.com/skynetservices/skydns)

    - 可以参考其中的安装方法

2. [Kubernetes DNS服务的安装与配置](http://blog.csdn.net/bluishglc/article/details/52438917)

    - 可以参考其中的架构图

    - 提到了dns服务本地安装与pod安装, 不必单独装skydns了

安装方法在参考文章1中写的已经足够简洁明了, 这里不再强调.

> 关于Skydns和Kube2sky是在本地安装还是以Pod的方式安装到Kube集群里，笔者在网上看到两种方式都有，但是笔者对本地安装的方式持怀疑态度，主要是涉及到虚拟网络和物理网络的联通性问题，具体地说就是Skydns Server的IP应该kube集群虚拟网络中的某个IP地址，也就是说这个IP需要在kube-apiserver启动参数–service-cluster-ip-range指定的IP地址范围内。而Skydns如果是本地化安装，是无法绑定DNS Server的IP为一个虚拟网络的IP（就是参数-addr的值）。笔者看到网上有文章是这样做的，但是笔者没有尝试成功，不排除需要某种网络位置才能打通两个网络之间的通信。总的来说，还是倾向使用镜像方式安装。

安装完成后, 

```
[root@localhost log]# kubectl -s http://172.32.100.90:8080 get pods --namespace=kube-system
NAME                        READY     STATUS             RESTARTS   AGE
kube-dns-3362242485-x595c   1/3       CrashLoopBackOff   24         42m
kubernetes-dashboard-1790416121-qtrpf   0/1       CrashLoopBackOff   7         15m

```

获取所有命名空间的服务状态`kubectl -s http://172.32.100.90:8080 get svc --all-namespaces`