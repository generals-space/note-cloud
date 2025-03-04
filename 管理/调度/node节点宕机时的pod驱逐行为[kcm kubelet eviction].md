# node节点宕机时的pod驱逐行为[kcm kubelet eviction].md

参考文章

1. [谈谈 K8S 的 pod eviction](http://wsfdl.com/kubernetes/2018/05/15/node_eviction.html)
    - 非常透彻, 值得一看
    - K8S `pod eviction` 机制, 某些场景下如节点`NotReady`, 资源不足时, 把 pod 驱逐至其它节点
    - 有两个组件可以发起`pod eviction`: `kube-controller-manager`和`kubelet`, 及这两个场景的具体介绍.
    - `kube-controller-manager`发起的驱逐, 效果需要商榷
2. [kubernetes之node 宕机，pod驱离问题解决](https://www.cnblogs.com/cptao/p/10911959.html)
    - [ ] `kube-controller-manager`的`--pod-eviction-timeout`选项不起作用
    - [x] 部署文件的污点设置`tolerations`, `tolerationSeconds: 10`有效, 当`NotReady`时间过长, 就会重新调度.
3. [k8s 容器的驱逐时间是多少 节点notready api kubelet](https://www.cnblogs.com/gaoyuechen/p/16529774.html)
    - `--node-monitor-period`: 节点控制器(node controller) 检查每个节点的间隔，默认5秒。
    - `--node-monitor-grace-period`: 节点控制器判断节点故障的时间窗口, 默认40秒。即40 秒没有收到节点消息则判断节点为故障。
    - `--pod-eviction-timeout`: 当节点故障时，controller-manager允许pod在此故障节点的保留时间，默认300秒。即当节点故障5分钟后，controller-manager开始在其他可用节点重建pod。
        - `--pod-eviction-timeout`已在1.27版本移除[Removal of --pod-eviction-timeout command line argument](https://kubernetes.io/blog/2023/03/17/upcoming-changes-in-kubernetes-v1-27/#removal-of-pod-eviction-timeout-command-line-argument)

kubelet 发起的驱逐，往往是资源不足导致.

由 kube-controller-manager 发起的驱逐, 则一般为 kubelet 无响应(无响应的原因有很多种)
