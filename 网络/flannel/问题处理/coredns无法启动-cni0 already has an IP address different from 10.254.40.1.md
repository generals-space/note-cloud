# coredns无法启动-cni0 already has an IP address different from 10.254.40.1

## 问题描述

coredns有一个Pod一直处于`ContainerCreating`状态.

```console
$ kwd pod
NAME                                    READY   STATUS              RESTARTS   AGE    IP                NODE            NOMINATED NODE   READINESS GATES
coredns-67c766df46-2dmd9                0/1     Running             2339       25d    10.254.0.5        k8s-master-01   <none>           <none>
coredns-67c766df46-lr4sf                0/1     ContainerCreating   0          4h5m   <none>            k8s-worker-02   <none>           <none>
```

`describe`一下有如下输出

```
Events:
  Type     Reason                  Age                 From                    Message
  ----     ------                  ----                ----                    -------
  Normal   Scheduled               <unknown>           default-scheduler       Successfully assigned kube-system/coredns-67c766df46-lr4sf to k8s-worker-02
  Warning  FailedCreatePodSandBox  36s                 kubelet, k8s-worker-02  Failed create pod sandbox: rpc error: code = Unknown desc = failed to set up sandbox container "30c7d82500c8f5a3e830be46a265fbc3f4deb85cc1bc618328b9225ecf06be23" network for pod "coredns-67c766df46-lr4sf": networkPlugin cni failed to set up pod "coredns-67c766df46-lr4sf_kube-system" network: failed to set bridge addr: "cni0" already has an IP address different from 10.254.40.1/24
```

重启Pod是无法解决的, 本来以为这下只能重启主机或reset集群了, 不过从事容器云这么长这么时间, 这点问题还是应该能解决才行.

## 解决方法

event信息中提到, `cni0`接口上已经有一个IP段了, 所以ipam阶段失败.

登陆到问题pod所在的`k8s-worker-02`主机, 发现的确如此.

```
9: cni0: <BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether ae:98:6b:21:aa:ab brd ff:ff:ff:ff:ff:ff
    inet 10.254.34.1/24 scope global cni0
       valid_lft forever preferred_lft forever
    inet6 fe80::ac98:6bff:fe21:aaab/64 scope link
       valid_lft forever preferred_lft forever
```

将该接口上的`10.254.34.1/24`删除.

```
ip a del 10.254.34.1/24 dev cni0
```

再次查看时, 该接口立刻就又拥有了一个IP段, 正是`10.254.40.1/24`.

```
9: cni0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1450 qdisc noqueue state UP group default qlen 1000
    link/ether ae:98:6b:21:aa:ab brd ff:ff:ff:ff:ff:ff
    inet 10.254.40.1/24 scope global cni0
       valid_lft forever preferred_lft forever
    inet6 fe80::ac98:6bff:fe21:aaab/64 scope link
       valid_lft forever preferred_lft forever
```

此时coredns的Pod也正常启动了.
