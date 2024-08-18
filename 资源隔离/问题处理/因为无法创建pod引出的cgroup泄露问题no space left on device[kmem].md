# 因为无法创建pod引出的cgroup泄露问题

参考文章

1. [Deployment doesn't create new replica set nor give any error](https://github.com/kubernetes/kubernetes/issues/36117)
2. [Cgroup泄漏--潜藏在你的集群中](https://tencentcloudcontainerteam.github.io/2018/12/29/cgroup-leaking/)
    - 复现条件中做的实验值得借鉴
3. [kubernetes频繁刷no space left on device错误](https://www.orchome.com/1511)
4. [Kernel not freeing memory cgroup causing no space left on device](https://github.com/moby/moby/issues/29638)
5. [Cgroup leaking, no space left on /sys/fs/cgroup](https://github.com/kubernetes/kubernetes/issues/70324)
6. [cgroup limit reached - no space left on device](https://stackoverflow.com/questions/45278379/cgroup-limit-reached-no-space-left-on-device)
7. [K8S 问题排查： cgroup 内存泄露问题 - kmem](https://www.cnblogs.com/leffss/p/15019898.html)

## 场景描述

某个开发团队在用kube集群, 把pod删除后突然不再新建, 给我整懵了. 很突然, 之前还好好的.

最初以为是node上的label被误删, 导致无法找到合适的节点去调度ta, 我手动添加了一下, 但还是没有新的pod创建.

于是我新开了一个命名空间testns, 在其中创建了一个pod资源, 仍然无法创建. 所以这其实并不是偶然事件, 一定是哪里出问题了.

## 排查思路

但是目前最重要的是把开发想要创建的pod启动起来, 先不管其他pod的问题, 想着既然之前还可以重建, 现在不行, 那么需要查一下对应的deploy资源. 我describe了一下, 看到如下部分的结果, 觉得很可疑.

```log
$ kubectl describe deploy 目标deploy资源
OldReplicaSets:		fika-io-749979362 (3/3 replicas created)
NewReplicaSet:		<none>
No events.
```

因为deployment底层是通过rs来管理pod数量的, 但是旧的rs失效后没有新的rs创建就有问题了. 

通过google, 我找到了参考文章1, 有回答说可能是controller manager出了问题, 建议查一下其日志.

于是我又查了一个kube-system空间下的情况, 发现controller和scheduler都处于CrashBackoff的状态, 且已经重启了几百次. 我describe了一下这两个对象, 发现有如下明显的报错:

```log
oci runtime error: process_linux.go:258: applying cgroup configuration for process caused "mkdir /sys/fs/cgroup/memory/docker/406cfca0c0a597091854c256a3bb2f09261ecbf86e98805414752150b11eb13a: no space left on device"
```

我几乎条件反射式地查了下硬盘空间, 50G才用了25G, 不是这个问题. 

于是又是一番google.

最终腾讯云团队的一篇文章给出了答案, 建议读一下.

由于是内核的问题(内核版本3.10), 而且没有解决办法, 只能重启. 要彻底解决此问题, 只能升级内核, 否则仍然会出现这种问题.
