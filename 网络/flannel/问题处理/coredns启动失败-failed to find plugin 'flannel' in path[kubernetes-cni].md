# 

参考文章

1. [failed to find plugin “flannel” in path [/opt/cni/bin]，k8sNotReady解决方案](https://blog.csdn.net/qq_29385297/article/details/127682552)
    - 手动从[containernetworking/plugins](https://github.com/containernetworking/plugins/releases/tag/v0.8.6)仓库下载release包, 拷贝`flannel`文件到`/opt/cni/bin/`目录下.
    - 在1.0.0版本后CNI Plugins库中就没有`flannel`
2. [使用kubeadm搭建kubernetes单机master，亲测无异常](https://developer.aliyun.com/article/827206)
3. [kubernetes入门：使用kubeadm搭建单机master，亲测无异常，建议收藏](https://bbs.huaweicloud.com/blogs/306548)
    - 与1是同一篇文章

flannel: 0.11.0

```log
$ kwd pod
NAME                                    READY   STATUS              RESTARTS   AGE   IP               NODE            NOMINATED NODE   READINESS GATES
coredns-7f9c544f75-9qn66                0/1     ContainerCreating   0          21m   <none>           k8s-master-01   <none>           <none>
coredns-7f9c544f75-r7kdk                0/1     ContainerCreating   0          21m   <none>           k8s-master-01   <none>           <none>
kube-flannel-ds-amd64-6tksh             1/1     Running             0          93s   172.16.156.172   k8s-master-01   <none>           <none>
$ kde pod coredns-7f9c544f75-9qn66 | tail
                 node-role.kubernetes.io/master:NoSchedule
                 node.kubernetes.io/not-ready:NoExecute for 300s
                 node.kubernetes.io/unreachable:NoExecute for 300s
Events:
  Type     Reason                  Age                 From                    Message
  ----     ------                  ----                ----                    -------
  Warning  FailedScheduling        84s (x16 over 22m)  default-scheduler       0/1 nodes are available: 1 node(s) had taints that the pod didn't tolerate.
  Normal   Scheduled               82s                 default-scheduler       Successfully assigned kube-system/coredns-7f9c544f75-9qn66 to k8s-master-01
  Warning  FailedCreatePodSandBox  82s                 kubelet, k8s-master-01  Failed to create pod sandbox: rpc error: code = Unknown desc = failed to setup network for sandbox "508e458d484ee31b8411bbb2834044e84152896e55ff8b6ec0c1ba0abd897e8f": plugin type="flannel" failed (add): failed to find plugin "flannel" in path [/opt/cni/bin]
  Normal   SandboxChanged          5s (x7 over 81s)    kubelet, k8s-master-01  Pod sandbox changed, it will be killed and re-created.
```

我看了下`/opt/cni/bin`目录, 的确没有`flannel`文件. 

```log
$ pwd
/opt/cni/bin
$ ls
bandwidth  bridge  dhcp  dummy  firewall  host-device  host-local  ipvlan  loopback  macvlan  portmap  ptp  sbr  static  tuning  vlan  vrf
```

本来以为这东西是`flannel`镜像自己带的, 结果yaml文件里`initContainers`只拷贝了`/etc/cni/net.d/`下的配置文件, 镜像里也没有可执行程序.

按照参考文章1中所说, 要手动从[containernetworking/plugins](https://github.com/containernetworking/plugins/releases/tag/v0.8.6)仓库下载release包, 拷贝`flannel`文件到`/opt/cni/bin/`目录下.

而按照参考文章2, `flannel`本来应该包含在`kubernetes-cni`包中的, 这个包会随着`kubelet`一起安装.

```log
$ rpm -qa | grep cni
kubernetes-cni-1.2.0-0.x86_64
```

这里装的`kubernetes-cni`是1.2.0版本的, 而参考文章1中所说, 自1.0.0及之后, 该包就不再包含`flannel`了.

于是还是要手动下载, 解决.
