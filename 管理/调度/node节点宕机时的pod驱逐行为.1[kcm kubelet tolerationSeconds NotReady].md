# node节点宕机时的pod驱逐行为[kcm kubelet eviction].md

参考文章

1. [谈谈 K8S 的 pod eviction](http://wsfdl.com/kubernetes/2018/05/15/node_eviction.html)
    - 非常透彻, 值得一看
    - K8S `pod eviction` 机制, 某些场景下如节点`NotReady`, 资源不足时, 把 pod 驱逐至其它节点
    - 有两个组件可以发起`pod eviction`: `kube-controller-manager`和`kubelet`, 及这两个场景的具体介绍.
    - `kube-controller-manager`发起的驱逐, 效果需要商榷
2. [kubernetes之node 宕机，pod驱离问题解决](https://www.cnblogs.com/cptao/p/10911959.html)
    - [ ] `kube-controller-manager`的`--pod-eviction-timeout`选项不起作用
        - **没错, 就是不起作用**
    - [x] 部署文件的污点设置`tolerations`, `tolerationSeconds: 10`有效, 当`NotReady`时间过长, 就会重新调度.
3. [k8s 容器的驱逐时间是多少 节点notready api kubelet](https://www.cnblogs.com/gaoyuechen/p/16529774.html)
    - `--node-monitor-period`: 节点控制器(node controller) 检查每个节点的间隔，默认5秒。
    - `--node-monitor-grace-period`: 节点控制器判断节点故障的时间窗口, 默认40秒。即40 秒没有收到节点消息则判断节点为故障。
    - `--pod-eviction-timeout`: 当节点故障时，controller-manager允许pod在此故障节点的保留时间，默认300秒。即当节点故障5分钟后，controller-manager开始在其他可用节点重建pod。
        - `--pod-eviction-timeout`已在1.27版本移除[Removal of --pod-eviction-timeout command line argument](https://kubernetes.io/blog/2023/03/17/upcoming-changes-in-kubernetes-v1-27/#removal-of-pod-eviction-timeout-command-line-argument)
        - kube-controller-manager 的`--pod-eviction-timeout`选项根本不管用. 不如换成 kube-apiserver 的`--default-not-ready-toleration-seconds`与`--default-unreachable-toleration-seconds`选项.

kubelet 发起的驱逐, 往往是资源不足导致.

由 kube-controller-manager 发起的驱逐, 则一般为 kubelet 无响应(无响应的原因有很多种).

kube: 1.17.2, master x 1 + worker * 2

## 默认场景

worker-01, worker-02上各自部署一个 deployment、statefulset 的 Pod.

```log
[root@kube-master-01 /]# kwd pod
NAME                              READY   STATUS    RESTARTS   AGE    IP           NODE          
deploy-test-01-5c44845b5c-hs2lf   1/1     Running   0          120m   10.254.1.3   kube-worker-01
deploy-test-01-5c44845b5c-p2xts   1/1     Running   0          120m   10.254.2.4   kube-worker-02
sts-test-01-0                     1/1     Running   0          163m   10.254.2.2   kube-worker-02
sts-test-01-1                     1/1     Running   0          162m   10.254.1.2   kube-worker-01
```

13:05:11 停止 worker-02 上的 kubelet.

```log
[root@kube-worker-02 data]# systemctl stop kubelet
```

13:05:59 检测到 worker-02 NotReady, 与 kubelet 停止间隔约48s.

```log
[root@kube-master-01 ~]# kwd node -w
kube-worker-02   NotReady   <none>   3h3m   v1.17.2   172.12.0.4    <none>        CentOS Linux 7 (Core)   3.10.0-1062.el7.x86_64   containerd://1.3.9
```

13:11:04 worker-02 上的 deployment、statefulset 的 Pod 同时变为 Terminating, 与 Node NotReady 间隔约6分05秒(305s).

```log
[root@kube-master-01 ~]# kwd pod -w
deploy-test-01-5c44845b5c-p2xts   1/1     Terminating   0          124m   10.254.2.4   kube-worker-02
sts-test-01-0                     1/1     Terminating   0          167m   10.254.2.2   kube-worker-02
deploy-test-01-5c44845b5c-fwq96   0/1     Pending       0          0s     <none>       <none>        
deploy-test-01-5c44845b5c-fwq96   0/1     Pending       0          1s     <none>       kube-master-01
```

其中, deployment Pod **立刻**开始被重建, 并在 master-01 上成功运行, 但是 statefulset Pod 没有, 且原来的 Pod 一直处于 Terminating, 并未被删除, 此时会存在2个 deployment Pod.

------

13:30:04 启动 worker-02 上的 kublet.

> 3秒内`kwd node -w`与`kwd pod -w`就能感知到, 会刷新几条状态, 只不过还是 NotReady 或 Terminating, 没变.

```log
[root@kube-worker-02 data]# systemctl start kubelet
```

13:30:16 worker-02 变为 Ready.

13:30:46 原来处于 Terminating 的 Pod 真正被删除, 同时 statefulset Pod 开始被重建.

## tolerationSeconds

按照参考文章2, 在 deployment/statefulset yaml 中添加 tolerations{}, 可以将从 Node NotReady 到 Pod 变为 Terminating 状态的间隔, 由300s缩短到10s.

deployment Pod 仍然在变为 Terminating 的同时开始重建, 而 statefulset Pod 仍然保持 Terminating.

```yaml
kind: Deployment
spec:
  template:
    spec:
      containers: {}
      tolerations:
      - key: "node.kubernetes.io/not-ready"
        operator: "Exists"
        effect: "NoExecute"
        tolerationSeconds: 10
      - key: "node.kubernetes.io/unreachable"
        operator: "Exists"
        effect: "NoExecute"
        tolerationSeconds: 10
```

其实这两个参数默认由 apiserver 的`--default-not-ready-toleration-seconds`与`--default-unreachable-toleration-seconds`控制, 默认都是300s. 如果想要全局设置`tolerationSeconds`的话, 就可以修改这两个参数了.

> kube-controller-manager 的`--pod-eviction-timeout`选项根本不管用.

## --node-monitor-grace-period 调整 Node NotReady 的判断间隔

在 kube-controller-manager yaml 中增加如下配置, 可以调整 Node NotReady 的判断间隔, 默认为 40s, kcm 重启后即可生效.

```yaml
    - --node-monitor-grace-period=10s
```
