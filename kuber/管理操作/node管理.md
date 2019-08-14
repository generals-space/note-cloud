# node管理

参考文章

1. [向Kubernetes集群添加／删除Node](https://blog.51cto.com/wutengfei/2113791)
2. [使用 kubectl drain 从集群中移除节点](https://www.cnblogs.com/weifeng1463/p/10359581.html)
3. [Nodes Rancher官方文档](https://rancher.com/docs/rancher/v2.x/en/cluster-admin/nodes/#aggressive-and-safe-draining-options-for-rancher-prior-to-v2-2-x)
    - drain驱散失败的情况分析及解决方案
4. [kubectl drain](http://kubernetes.kansea.com/docs/user-guide/kubectl/kubectl_drain/)
    - drain的几个强制性选项及其作用

kubernetes版本: 1.15.10

## drain驱散

- kubectl drain node名
- kubectl cordon node名
- kubectl uncordon node名

cordon/uncordon只修改node状态(`SchedulingDisabled`不可调度), 不涉及pod操作.

而drain首先将node设置为`SchedulingDisabled`, 之后迁移pod(先在其他node上启动新pod, 然后移除本地pod).

## 删除节点

一般执行drain操作后就可以使用`kubectl delete node node名`删除node节点了.

`delete node`会将目标node节点上的kuber组件全部停止并删除, 不过貌似iptables和ipvsadm没有清空???, 还是手动执行一下`kubeadm reset`吧.

`delete node`可以删除master节点, 不过目标节点的kuber组件并未停止, 同样需要手动`kubeadm reset`.

## FAQ

### drain驱散操作失败分析及解决

```
$ kubectl drain k8s-worker-7-17
node/k8s-worker-7-17 cordoned
error: unable to drain node "k8s-worker-7-17", aborting command...

There are pending nodes to be drained:
 k8s-worker-7-17
cannot delete Pods not managed by ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet (use --force to override): default/fortest
cannot delete DaemonSet-managed Pods (use --ignore-daemonsets to ignore): kube-system/kube-proxy-np5gf, kube-system/terway-vlan-171-z25mh
```

按照参考文章3, drain操作失败有3种可能

1. there are pods not managed by a ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet
    - 有一些pod没有被上述对象管理, 没有可能再被调度到其他node上. 比如**单纯的pod资源(没有使用deployment方式部署), 删除就直接删了, kuber无法在其他node重新创建**. kuber希望管理员自行处理这类pod的删除逻辑, 加上`--force`选项可强制删除此类pod.
2. there are DaemonSet-managed pods
    - drain默认不会对被`DaemonSet`对象管理的pod做任何操作, 使用`--ignore-daemonsets`忽略此类pod.
3. there are pods using emptyDir
    - 如果有pod使用`emptyDir`存储本地数据, `emptyDir`中的数据会随着pod的移除而删除. 与第一种情况一样, kuber希望管理员能明确指定, 使用`--delete-local-data`强制删除(可能管理员需要自行备份数据).

注意: 即使drain失败, 也只是某些pod没有被删除, 但是该节点还是被修改为不可调度的状态.

```
$ kubectl get node
k8s-worker-7-17   Ready,SchedulingDisabled   <none>   2d18h   v1.15.0
```
