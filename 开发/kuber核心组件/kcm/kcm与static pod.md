# kcm与static pod

参考文章

1. [官方文档  Create static Pods](https://kubernetes.io/docs/tasks/configure-pod-container/static-pod/)

写在前面

kubernetes版本: v1.16.2(1主2从)

在研究`kube-controller-manager`源码的时候, 希望通过`go run main.go`执行测试程序. 不过由于分布式资源锁的存在, 只有成功获得锁的实例才能继续执行, 所以需要将主节点上运行着的`controller-manager`的pod删除. 

但删除后总是会重启, 而且pod输出的yaml文件没有Reference块, 也没有发现该ns下对应的deployment, daemonset, statefulset, rs, rc等相关的对象.

然后查看apiserver, scheduler的日志, 也没有发现可疑的输出.

最后快绝望的时候, 查看了master节点上的kubelet的日志.

```
Jan  5 13:48:00 k8s-master-01 kubelet[888]: I0105 13:48:00.531576     888 kubelet.go:1647] Trying to delete pod kube-controller-manager-k8s-master-01_kube-system 7ef7f259-6c68-4986-8b69-a89887e19b19
Jan  5 13:48:00 k8s-master-01 kubelet[888]: W0105 13:48:00.760125     888 kubelet.go:1651] Deleted mirror pod "kube-controller-manager-k8s-master-01_kube-system(7ef7f259-6c68-4986-8b69-a89887e19b19)" because it is outdated
```

于是再次搜索`mirror pod`, 找到了参考文章1.

------

按照参考文章1中所说, 查看kubelet的配置文件, `/var/lib/kubelet/config.yaml`(master与worker节点都有此文件)

```yaml
staticPodPath: /etc/kubernetes/manifests
```

在master节点的该目录下, 存在如下文件

```console
$ pwd
/etc/kubernetes/manifests
$ ll
total 16
-rw------- 1 root root 1810 Dec 26 23:28 etcd.yaml
-rw------- 1 root root 2647 Dec 26 23:28 kube-apiserver.yaml
-rw------- 1 root root 2573 Dec 26 23:28 kube-controller-manager.yaml
-rw------- 1 root root 1160 Dec 26 23:28 kube-scheduler.yaml
```

worker节点中, 该目录为空.

------

`kubelet`服务在启动时会自动创建`staticPodPath`目录中声明的static pod, 然后监控这些pod的运行状态.

另外, `kubelet`也会定期扫描`staticPodPath`目录, 动态地创建或删除该目录下新增或不见的yaml部署文件.

1. `Static Pods`由某个节点上的kubelet服务管理, 而不是由`control plane`. 
2. `Static Pods`会绑定到某个指定的节点上.
3. 虽然不由`control plane`管理, 但是可以通过`kubectl`向apiserver查询到pod记录. 这是因为kubelet在自行管理`Static Pods`时, 还会向apiserver注册`Mirror Pods`. 
4. 通过`kubectl`删除`mirror pod`对象, 会引起相应的kubelet服务kill并重新创建pod实例, apiserver中也会对应更新.
