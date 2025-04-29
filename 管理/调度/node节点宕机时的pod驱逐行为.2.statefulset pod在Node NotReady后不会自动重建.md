# node节点宕机时的pod驱逐行为.2.statefulset pod在Node NotReady后不会自动重建

参考文章

1. [Pods stuck in Terminating state when worker node is down (never redeployed on healthy nodes), how to fix this?](https://stackoverflow.com/questions/68979835/pods-stuck-in-terminating-state-when-worker-node-is-down-never-redeployed-on-he)
2. [Statefulset should be able to evicted if the worker node goes down #74947](https://github.com/kubernetes/kubernetes/issues/74947)
3. [Add new Taint Effect: ForceEviction #719](https://github.com/kubernetes/enhancements/pull/719)
4. [跟我学 K8S--代码: Kubernetes StatefulSet 代码分析与Unknown 状态处理](https://segmentfault.com/a/1190000019488735)

官方设计如此, 不想让 statefulset pod 自动迁移, 吵了好多年.
