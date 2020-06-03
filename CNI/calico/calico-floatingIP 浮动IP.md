# calico-floatingIP 浮动IP

参考文章

1. [calico官方文档 - Add a floating IP to a pod](https://docs.projectcalico.org/networking/add-floating-ip)
    - floatingIP 相比于 service 资源的异同
    - floatingIP 功能的开启与使用方法

floatingIP, 类似于`keepalived`维护的VIP(`virtual IP`虚IP), 初始状态指向某一固定后端, 当这一后端因故障无法提供服务时, VIP将指向其他的后端.

在 kuber 概念中, 其实 service 类型资源已经基本实现了这样的功能, 不过与 service 相比, 除了 TCP/UDP/SCTP 协议(service 是由iptables/ipvs 实现的), floatingIP 还支持其他类型的协议(虽然这样的需求非常少)

## 1. 配置方法

首先需要更新 configmap 资源, 找到`cni_network_config`块, 在`plugins`字段下添加如下内容.

```json
    "feature_control": {
        "floating_ips": true
    }
```

注意 json 格式中的逗号, 另外, 更新完成 configmap 后还需要重启.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: centos7-pod
  labels:
    app: centos7
  annotations:
    ## 注解值只能是字符串, 不可以是数组.
    ## cni.projectcalico.org/floatingIPs: ["10.0.0.1"]   ## 错误
    ## cni.projectcalico.org/floatingIPs: "[\"10.0.0.1\"]"  ## 正确
    cni.projectcalico.org/floatingIPs: '["10.0.0.1"]'    ## 正确
spec:
  containers:
  - name: centos7
    image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos7
    command: ['tail', '-f', '/etc/os-release']
```

`floatingIPs`中写的IP必须是合法的, 未被占用的, PodCIDR 中的 IP. 

## 2. 使用方法

`floatingIPs`应该是多个 Pod 定义同一个 IP 值, 如`10.0.0.1`, 在这些 Pod 创建完成后, calico 会将`10.0.0.1`指向这此 Pod 中的其中一个(一般是第1个创建的). 当`floatingIPs` 后端 Pod 故障退出时, calico 会将其再指向另外一个 Pod, 实现切换.

不过这些 Pod 本身的 IP 并非是`10.0.0.1`, 而是各自拥有本身的 IP. 而这个`floatingIPs`地址在集群中根本没地方找, 所以运维人员需要找个地方记下来...

另外, 我试了下在`Deployment`和`DaemonSet`中指定`floatingIPs`, 本来以为可以代理ta们派生出来的 Pod 的, 结果不是...看来目前`floatingIPs`只支持 Pod 类型的资源.

以下是`DeamonSet`的示例(失败)

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: test-ds
  labels:
    app: test-ds
  annotations:
    ## 注解值只能是字符串, 不可以是数组.
    ## cni.projectcalico.org/floatingIPs: ["10.0.0.1"]   ## 错误
    ## cni.projectcalico.org/floatingIPs: "[\"10.0.0.1\"]"  ## 正确
    cni.projectcalico.org/floatingIPs: '["172.16.80.200"]'    ## 正确
  namespace: default
spec:
  selector:
    matchLabels:
      app: test-ds
  template:
    metadata:
      labels:
        app: test-ds
    spec:
      containers:
      - name: centos7
        image: registry.cn-hangzhou.aliyuncs.com/generals-space/centos7-devops
        imagePullPolicy: IfNotPresent
        command: ["tail", "-f", "/etc/os-release"]
      ## 允许在master节点部署
      tolerations:
      ## 这一段表示pod需要容忍拥有master角色的node们, 
      ## 且这个污点的效果是 NoSchedule.
      ## 因为key, operator, value, effect可以唯一确定一个污点对象(Taint).
      - key: node-role.kubernetes.io/master
        operator: Exists
        effect: NoSchedule
```

