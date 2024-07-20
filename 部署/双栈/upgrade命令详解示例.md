# upgrade命令详解示例

参考文章

1. [kubernetes集群版本升级攻略](https://blog.51cto.com/newfly/2440901)
    - 超全超详细
2. [kubeadm部署k8s-1.9](https://www.cnblogs.com/lixuebin/p/10813970.html)
    - 命令行中`--feature-gates`选项及config配置文件中`featureGates`字段的使用方法.
3. [Failed to parse subnet when creating IPv4/IPv6 dual stack](https://github.com/kubernetes/kubeadm/issues/1828)
4. [fialed to test IPv6DualStack feature of release version 1.16.0](https://github.com/kubernetes/kubernetes/issues/83006)
    - `serviceSubnet: Invalid value: "10.96.0.0/12,2019:30::/24": couldn't parse subnet`
    - 开启双栈配置后, `networking`字段下, `podSubnet`的解析是正确的, 但是无法解析`serviceSubnet`字段.
    - the verison 1.16 of the kubeadm does not support to pass a comma separated list of to --serviceSubnet with "IPv4, IPv6"
5. [Kubernetes v1.17 参考指南 - kubeadm upgrade](https://www.bookstack.cn/read/kubernetes-1.17-reference/0321732120b50209.md)

## `kubeadm upgrade plan`

> 检查可升级到哪些版本，并验证您当前的集群是否可升级。 要跳过互联网检查，请传递可选的 [version] 参数 --参考文章5

1. 只检测目标版本是否可升级, 不检查其他的(比如`kubeadm-config.yaml`中定义的双栈参数)...
2. 跨版本升级(`1.16.2`->`1.18.2`)貌似可以, 没有报错, 只有Warning.
    - `WARNING: No recommended etcd for requested Kubernetes version (v1.18.2)`).
3. ta甚至不检查目标版本是否存在(`1.17.2000`)
4. 如果目标版本比当前版本小, 则会告诉你`Awesome, you're up-to-date! Enjoy!`.

...错了, 理论上来说要先升级`kubeadm`工具, 再执行`upgrade plan`. 所以当告诉你`Awesome, you're up-to-date! Enjoy!`的时候, 是可以升级的时候.

```log
## 两者的输出几乎一样, upgrade plan 并不检查具体参数.
## kubeadm upgrade plan v1.17.2 
$ kubeadm upgrade plan --config ./kubeadm-config.v1.17.2.yaml
[upgrade/config] Making sure the configuration is correct:
[preflight] Running pre-flight checks.
[upgrade] Making sure the cluster is healthy:
[upgrade] Fetching available versions to upgrade to
[upgrade/versions] Cluster version: v1.16.2
[upgrade/versions] kubeadm version: v1.16.2

Components that must be upgraded manually after you have upgraded the control plane with 'kubeadm upgrade apply':
COMPONENT   CURRENT       AVAILABLE
Kubelet     3 x v1.16.2   v1.17.2

Upgrade to the latest version in the v1.16 series:

COMPONENT            CURRENT   AVAILABLE
API Server           v1.16.2   v1.17.2
Controller Manager   v1.16.2   v1.17.2
Scheduler            v1.16.2   v1.17.2
Kube Proxy           v1.16.2   v1.17.2
CoreDNS              1.6.2     1.6.2
Etcd                 3.3.15    3.3.15-0

You can now apply the upgrade by executing the following command:

        kubeadm upgrade apply v1.17.2

Note: Before you can perform this upgrade, you have to update kubeadm to v1.17.2.
```

## `kubeadm upgrade diff`

> 显示哪些差异将被应用于现有的静态 pod 资源清单。 --参考文章5

**静态Pod(static pod)**是kuber中的专有概念, 一般是指在`/etc/kubernetes/manifests`目录下的`yaml`部署文件, 包括`apiserver`, `controller-manager`, `scheduler`和`etcd`.

不过`diff`貌似不检查`etcd`配置. 另外, 反而会检测`/etc/kubernetes/admin.conf`文件...

`diff`需要我们使用`--config`选项传入作为对比的配置文件, 不支持在命令行传入. 这里使用`kubeadm-config.v1.17.2.yaml`为例, 开启了双栈.

```log
$ kubeadm upgrade diff --config ./kubeadm-config.v1.17.2.yaml 
--- /etc/kubernetes/manifests/kube-apiserver.yaml
+++ new manifest
@@ -21,6 +21,7 @@
     - --etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt
     - --etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key
     - --etcd-servers=https://127.0.0.1:2379
+    - --feature-gates=IPv6DualStack=true
     - --insecure-port=0
     - --kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt
     - --kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key
@@ -37,7 +38,7 @@
     - --service-cluster-ip-range=10.96.0.0/12
     - --tls-cert-file=/etc/kubernetes/pki/apiserver.crt
     - --tls-private-key-file=/etc/kubernetes/pki/apiserver.key
-    image: registry.cn-hangzhou.aliyuncs.com/google_containers/kube-apiserver:v1.16.2
+    image: registry.cn-hangzhou.aliyuncs.com/google_containers/kube-apiserver:v1.17.2
     imagePullPolicy: IfNotPresent
     livenessProbe:
       failureThreshold: 8
--- /etc/kubernetes/manifests/kube-controller-manager.yaml
+++ new manifest
@@ -16,19 +16,20 @@
     - --authorization-kubeconfig=/etc/kubernetes/controller-manager.conf
     - --bind-address=127.0.0.1
     - --client-ca-file=/etc/kubernetes/pki/ca.crt
-    - --cluster-cidr=10.254.0.0/16
+    - --cluster-cidr=10.254.0.0/16,2019:20::/24
     - --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt
     - --cluster-signing-key-file=/etc/kubernetes/pki/ca.key
     - --controllers=*,bootstrapsigner,tokencleaner
+    - --feature-gates=IPv6DualStack=true
     - --kubeconfig=/etc/kubernetes/controller-manager.conf
     - --leader-elect=true
     - --node-cidr-mask-size=24
     - --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt
     - --root-ca-file=/etc/kubernetes/pki/ca.crt
     - --service-account-private-key-file=/etc/kubernetes/pki/sa.key
-    - --service-cluster-ip-range=10.96.0.0/12
+    - --service-cluster-ip-range=10.96.0.0/12,2019:30::/24
     - --use-service-account-credentials=true
-    image: registry.cn-hangzhou.aliyuncs.com/google_containers/kube-controller-manager:v1.16.2
+    image: registry.cn-hangzhou.aliyuncs.com/google_containers/kube-controller-manager:v1.17.2
     imagePullPolicy: IfNotPresent
     livenessProbe:
       failureThreshold: 8
--- /etc/kubernetes/manifests/kube-scheduler.yaml
+++ new manifest
@@ -14,9 +14,10 @@
     - --authentication-kubeconfig=/etc/kubernetes/scheduler.conf
     - --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
     - --bind-address=127.0.0.1
+    - --feature-gates=IPv6DualStack=true
     - --kubeconfig=/etc/kubernetes/scheduler.conf
     - --leader-elect=true
-    image: registry.cn-hangzhou.aliyuncs.com/google_containers/kube-scheduler:v1.16.2
+    image: registry.cn-hangzhou.aliyuncs.com/google_containers/kube-scheduler:v1.17.2
     imagePullPolicy: IfNotPresent
     livenessProbe:
       failureThreshold: 8
```

可以看到, 主要就是`images`版本, 各组件启动参数中添加了`--feature-gates`选项, 另外就是`--cluster-cidr`和`--service-cluster-ip-range`选项的变化.

## kubeadm upgrade apply

> 将 Kubernetes 集群升级到指定版本 --参考文章5

```
kubeadm upgrade apply v1.17.2 --feature-gates IPv6DualStack=true
kubeadm upgrade apply --config ./kubeadm-config.v1.17.2.yaml
```

## kubeadm upgrade node

> 升级集群中某个节点的命令 --参考文章5

这条命令需要在 worker 节点上执行, 所以 worker 节点上也要安装并升级 `kubeadm`.

当 control plane 与 非 control plane 的 master 节点升级完成后, 便可以一台一台升级 worker 节点了. 无需再执行 `kubeadm upgrade plan`.

```log
$ kubeadm upgrade node
[upgrade] Reading configuration from the cluster...
[upgrade] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -oyaml'
[upgrade] Skipping phase. Not a control plane node.
[kubelet-start] Downloading configuration for the kubelet from the "kubelet-config-1.17" ConfigMap in the kube-system namespace
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[upgrade] The configuration for this node was successfully updated!
[upgrade] Now you should go ahead and upgrade the kubelet package using your package manager.
```

ta会自动读取名为`kubeadm-config`的 cm 资源(就是`kubeadm-config.v1.17.2.yaml`的内容), 然后生成`/var/lib/kubelet/config.yaml`配置文件, 最后只要重启一下 kubelet 即可(当然, 用户需要事先升级`kubelet`组件(与`kubeadm`一同升级即可)).
