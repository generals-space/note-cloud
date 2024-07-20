参考文章

1. [Kubeadm init with failed node error](https://discuss.kubernetes.io/t/kubeadm-init-with-failed-node-error/20042)

## 问题描述

- OS: ubuntu 20.04
- hostname: ubuntu
- IP: 192.168.128.135

执行`kubeadm init`显示超时失败, 此时kubelet已经启动, 但是日志中显示

```log
$ kubeadm init --kubernetes-version=v1.25.1 --image-repository=registry.aliyuncs.com/google_containers --ignore-preflight-errors=all
Sep 20 00:14:40 ubuntu kubelet[36562]: W0920 00:14:40.631878   36562 reflector.go:424] vendor/k8s.io/client-go/informers/factory.go:134: failed to list *v1.Service: Get "https://192.168.128.135:6443/api/v1/services?limit=500&resourceVersion=0": dial tcp 192.168.128.135:6443: connect: connection refused
Sep 20 00:14:40 ubuntu kubelet[36562]: E0920 00:14:40.633037   36562 reflector.go:140] vendor/k8s.io/client-go/informers/factory.go:134: failed to watch *v1.Service: failed to list *v1Service: Get "https://192.168.128.135:6443/api/v1/services?limit=500&resourceVersion=0": dial tcp 192.168.128.135:6443: connect: connection refused
Sep 20 00:14:40 ubuntu kubelet[36562]: E0920 00:14:40.663165   36562 kubelet.go:2448] "Error getting node" err="node \"ubuntu\" not found"
Sep 20 00:14:40 ubuntu kubelet[36562]: E0920 00:14:40.663165   36562 kubelet.go:2448] "Error getting node" err="node \"ubuntu\" not found"
Sep 20 00:14:40 ubuntu kubelet[36562]: E0920 00:14:40.663165   36562 kubelet.go:2448] "Error getting node" err="node \"ubuntu\" not found"
Sep 20 00:14:40 ubuntu kubelet[36562]: E0920 00:14:40.663165   36562 kubelet.go:2448] "Error getting node" err="node \"ubuntu\" not found"
...省略
```

## 

最初以为是无法解析 ubuntu 主机名, 于是在hosts配置了`192.168.128.135 ubuntu`的解析, reset后重新安装, 但是没用.

后来发现是由`192.168.128.135:6443`连接失败引起的, `docker ps -a`根本没容器生成, 理论上这个时候 apiserver, scheduler 等组件应该以 static pod 形式启动才对.

翻阅官方文档后才知道 kubernetes 1.25.0 后正式移除了 docker 接口, 转而使用实现了 ORI 接口的 runtime.

由于 containerd 还是 docker 自身相关的接口, 所以这次一步到位, 直接使用 crio.

> 好吧, 其实是因为 docker+cri-dockerd 也无效.

> 在 centos 7下安装 crio 后无法启动(xfs文件系统的问题), 不过 ubuntu 上启动 crio 倒是挺正常.

其实最开始换成 crio 接口还是无效, 本来kubelet要创建容器一定是要通过接口让 CRI 运行时完成的, 但是static容器没启动 kubelet 也没报错, 一点信息也看不出来.

后来想修改一下 crio 本身的日志级别, 查看了下ta的配置文件`/etc/crio/crio.conf`

```conf
# The image used to instantiate infra containers.
# This option supports live configuration reload.
pause_image = "k8s.gcr.io/pause:3.2"
```

将`pause_image`字段换成`registry.aliyuncs.com/google_containers/pause:3.8`, 清理后重新执行`kubeadm init`, 竟然可以了...

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress:
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///var/run/crio/crio.sock
  imagePullPolicy: IfNotPresent
  ## name: ubuntu
  ## taints: null
  ## kubeletExtraArgs:
  ##   v: "5"
---
apiServer:
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta3
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns: {}
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.aliyuncs.com/google_containers
kind: ClusterConfiguration
kubernetesVersion: 1.25.1
networking:
  dnsDomain: cluster.local
  serviceSubnet: 10.96.0.0/12
scheduler: {}
---
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
cgroupDriver: systemd
```
