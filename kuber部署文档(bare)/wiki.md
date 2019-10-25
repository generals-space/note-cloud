
1. [Kubernetes扫盲](http://blog.csdn.net/ztsinghua/article/details/52385376)

    - 各个组件讲解的很详细, 虽然有些不够通俗

2. [kubernetes1.6 安装DNS（四）](http://blog.csdn.net/u010278923/article/details/71152796)

    - 这个系列讲得都不错

3. [利用 Kubernetes Service 的 selector 无痛运维在线 pod](https://segmentfault.com/a/1190000007109399#articleHeader0)

    - service的label与selector字段解释

4. [Kubernetes 中的PodIP、ClusterIP 和外部IP](http://blog.csdn.net/liukuan73/article/details/54773579)


5. [Kubernetes如何使用kube-dns实现服务发现](http://www.cnblogs.com/ilinuxer/p/6188804.html)

```
docker pull registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0
docker tag registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0 gcr.io/google_containers/pause-amd64:3.0
```

参考文章

1. [kubelet组件操作API - kubernetes 简介： kubelet 和 pod](http://cizixs.com/2016/10/25/kubernetes-intro-kubelet?utm_source=tuicool&utm_medium=referral)

2. [ovs+kube配置1 - CentOS 7实战Kubernetes部署](http://www.infoq.com/cn/articles/centos7-practical-kubernetes-deployment)

3. [ovs+kube配置1 - kube相关服务systemd脚本1](https://github.com/yangzhares/GetStartingKubernetes)

3. [ovs+kube配置2 - 基于OVS的Kubernetes多节点环境搭建](http://bingotree.cn/?p=828)

4. [kube相关服务systemd脚本2 - 轻松了解Kubernetes认证功能](http://www.tuicool.com/articles/byUnQn7)

5. [kubernetes 集群的安装部署](http://www.cnblogs.com/galengao/p/5780938.html)

6. [kubernetes-create-pods](http://www.liuhaihua.cn/archives/416728.html)

从kubelet直接获取pid/container信息: `curl http://172.32.100.70:10255/pods`

kubelet健康检查: `curl http://127.0.0.1:10248/healthz`

------

从apiserver获取节点信息

```
[root@localhost ~]# kubectl -s http://172.32.100.90:8080 get nodes
NAME            STATUS    AGE       VERSION
172.32.100.70   Ready     4m        v1.8.0-alpha.0.367+0613ae5077b280
172.32.100.80   Ready     4m        v1.8.0-alpha.0.367+0613ae5077b280
```

ovs+docker主机间容器通信必要步骤

1. ip_forward=1 每次重启network服务都会重置sysctl配置, 最好写在文件中.

2. route路由, 写配置文件

3. kbr0网卡配置(不可加HWADDR参数)

4. type=gre/vxlan应该都没关系

5. 关闭防火墙/selinux

6. NetworkManager关闭

ovs网络划分规则, 单个主机拥有其独立子网, 还是可以多主机同属同一子网, 容器ip在什么范围内保持唯一?

20170607
------

flannel, Open vSwitch的存在目的是实现容器之间的跨主机通信, 否则容器只能与其宿主机沟通. 跨主机的容器需要路由与ip支持.

Etcd来存储每台机器的上子网地址

FAQ

1. kubernetes的release只有8M, 好像没法用???