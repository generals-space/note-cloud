# kubeadm创建基于containerd的kubernetes集群[cri docker crictl ctr]

参考文章

1. [kubernetes那些事 —— 使用containerd作为容器运行时](https://zhuanlan.zhihu.com/p/334742204)
    - 灵感来源
2. [kubeadm安装k8s高可用集群](https://blog.csdn.net/wuxingge/article/details/119462915)
    - 灵感来源
    - containerd 版本与 docker 版本的 kubernetes 集群在部署上的区别
3. [Container Runtimes](https://kubernetes.io/docs/setup/production-environment/container-runtimes/)
    - 官方文档
    - kubernetes 针对不同底层 runtime 的部署方法
4. [kubeadm init 报错 ”unknown service runtime.v1alpha2.RuntimeService”](https://blog.csdn.net/weixin_40668374/article/details/124849090)
5. [Kubeadm unknown service runtime.v1alpha2.RuntimeService](https://github.com/containerd/containerd/issues/4581)
    - 官方issue
6. [containerd拉取私库镜像失败(kubelet)](https://blog.csdn.net/u010566813/article/details/125990298)
7. [搭建 k8s 环境](https://www.cnblogs.com/-ori/p/16968520.html)
8. [【K8S】ctr和crictl的区别](https://blog.csdn.net/u010157986/article/details/126118897)

## 前言

kube 在 1.24.0 开始, 移除了 dockershim, 不再直接与 dockerd 服务进行通信. 转而通过 CRI 接口协议, 直接接入 runtime. 支持 CRI 协议的 runtime 有.

- Docker Engine
    - 仍然与 dockerd 进行通信, 但是需要增加额外组件`cri-dockerd`, 可以说, 该工程表示的就是本次 kubernetes 移除的部分.
- containerd
    - docker 底层的容器实现, 调用关系大概是 docker -> dockerd -> contianerd -> runc
- CRI-O
    - redhat 维护的仓库, 其底层也是 runc

## 准备

kubernetes: v1.17.2

虽然 kube 是在 1.24 才正式移除 dockershim 的, 但其实很早就开始通过 CRI 协议兼容其他 runtime 了.

在 1.17 版本中, 可以通过`criSocket`指定 containerd 的 sock 地址. containerd 随 docker 安装, 也随 docker 服务启动而启动.

```yaml
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
nodeRegistration:
  ## 默认通过 dockershim 与 dockerd 通信, 现在直接指定 containerd 服务, 绕过 dockerd.
  ## criSocket: /var/run/dockershim.sock
  criSocket: unix:///var/run/containerd/containerd.sock
```

直接使用 kubeadm init 即可.

安装完后, criSocket 信息可以在 node yaml 的 annotation 中看到.

```yaml
## kubectl get node k8s-master-01 -oyaml
apiVersion: v1
kind: Node
metadata:
  annotations:
    kubeadm.alpha.kubernetes.io/cri-socket: unix:///var/run/containerd/containerd.sock
  labels:
    kubernetes.io/hostname: k8s-master-01
    node-role.kubernetes.io/master: ""
  name: k8s-master-01
```

## ctr 与 crictl

安装完成后, 无法通过 docker ps 看到容器了, docker images 也看不到镜像信息了.

ctr是containerd自带的CLI命令行工具, 可以单独下载镜像, 启动容器, 但是ta也没办法看到 kubernetes 运行的容器和镜像.

不过, 还存在一个crictl命令, ta是k8s中CRI（容器运行时接口）的客户端, k8s使用该客户端和containerd进行交互;

crictl, 顾名思义, 只要是实现了 CRI 接口的 runtime, 都可以通过该命令行管理.

```console
$ crictl ps
WARN[0000] runtime connect using default endpoints: [unix:///var/run/dockershim.sock unix:///run/containerd/containerd.sock unix:///run/crio/crio.sock unix:///var/run/cri-dockerd.sock]. As the default settings are now deprecated, you should set the endpoint instead.
ERRO[0000] unable to determine runtime API version: rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial unix /var/run/dockershim.sock: connect: no such file or directory"
WARN[0000] image connect using default endpoints: [unix:///var/run/dockershim.sock unix:///run/containerd/containerd.sock unix:///run/crio/crio.sock unix:///var/run/cri-dockerd.sock]. As the default settings are now deprecated, you should set the endpoint instead.
ERRO[0000] unable to determine image API version: rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial unix /var/run/dockershim.sock: connect: no such file or directory"
CONTAINER        IMAGE            CREATED           STATE      NAME                      ATTEMPT    POD ID           POD
8d5cb5cd82e4e    cba2a99699bdf    20 minutes ago    Running    kube-proxy                0          062cc1b19a0ea    kube-proxy-xxcqd
d76602bea5efa    f52d4c527ef2f    20 minutes ago    Running    kube-scheduler            0          11a2b0058e122    kube-scheduler-k8s-master-01
73b711ee173c5    da5fd66c4068c    20 minutes ago    Running    kube-controller-manager   0          fd271192295d9    kube-controller-manager-k8s-master-01
491e4e090f718    41ef50a5f06a7    20 minutes ago    Running    kube-apiserver            0          99359a9953fc4    kube-apiserver-k8s-master-01
f0a1eacd820ff    303ce5db0e90d    20 minutes ago    Running    etcd                      0          c64df243c9ad2    etcd-k8s-master-01
```

从上述输出的"WARN"信息可以看到, crictl 可以获取 kubernetes 所支持的所有 CRI runtime 接口.

由于我们是直接使用的 containerd 服务, 可以让 crictl 指定其连接的 sock 地址.

```console
$ crictl --runtime-endpoint unix:///var/run/containerd/containerd.sock ps -a
CONTAINER        IMAGE            CREATED           STATE      NAME                      ATTEMPT    POD ID           POD
8d5cb5cd82e4e    cba2a99699bdf    22 minutes ago    Running    kube-proxy                0          062cc1b19a0ea    kube-proxy-xxcqd
d76602bea5efa    f52d4c527ef2f    23 minutes ago    Running    kube-scheduler            0          11a2b0058e122    kube-scheduler-k8s-master-01
73b711ee173c5    da5fd66c4068c    23 minutes ago    Running    kube-controller-manager   0          fd271192295d9    kube-controller-manager-k8s-master-01
491e4e090f718    41ef50a5f06a7    23 minutes ago    Running    kube-apiserver            0          99359a9953fc4    kube-apiserver-k8s-master-01
f0a1eacd820ff    303ce5db0e90d    23 minutes ago    Running    etcd                      0          c64df243c9ad2    etcd-k8s-master-01
```

用如下命令长久设置 endpoint

```
crictl config runtime-endpoint /run/containerd/containerd.sock
```

## F&Q

### unknown service runtime.v1alpha2.RuntimeService

```log
[root@k8s-master-01 ~]# kubeadm init --config=./kubeadm-config.yaml
W1215 15:00:51.835929    2220 validation.go:28] Cannot validate kube-proxy config - no validator is available
W1215 15:00:51.835967    2220 validation.go:28] Cannot validate kubelet config - no validator is available
[init] Using Kubernetes version: v1.17.2
[preflight] Running pre-flight checks
error execution phase preflight: [preflight] Some fatal errors occurred:
	[ERROR CRI]: container runtime is not running: output: E1215 15:00:51.862184    2230 remote_runtime.go:948] "Status from runtime service failed" err="rpc error: code = Unimplemented desc = unknown service runtime.v1alpha2.RuntimeService"
time="2022-12-15T15:00:51+08:00" level=fatal msg="getting status of runtime: rpc error: code = Unimplemented desc = unknown service runtime.v1alpha2.RuntimeService"
, error: exit status 1
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
To see the stack trace of this error execute with --v=5 or higher
```

参考文章4中说, 删除`/etc/containerd/config.toml`, 即`containerd`服务的配置文件, 然后重启`containerd`服务即可.

按照参考文章5中所说, 其实是因为通过`yum install docker-ce`, 或是`yum install containerd.io`安装的`containerd`服务, 其`/etc/containerd/config.toml`中默认包含如下配置.

```ini
disabled_plugins = ["cri"]
```

而 kubernetes 就是通过 cri 接口协议与`containerd`进行通信的, 所以把这行注释掉, 再重启`containerd`即可.

### 无法下载 pause 镜像

#### 问题描述

kubeadm init 超时, 仍然是 kubelet 没能正常启动.

```log
[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s
[kubelet-check] Initial timeout of 40s passed.

Unfortunately, an error has occurred:
	timed out waiting for the condition

This error is likely caused by:
	- The kubelet is not running
	- The kubelet is unhealthy due to a misconfiguration of the node in some way (required cgroups disabled)

If you are on a systemd-powered system, you can try to troubleshoot the error with the following commands:
	- 'systemctl status kubelet'
	- 'journalctl -xeu kubelet'

Additionally, a control plane component may have crashed or exited when started by the container runtime.
To troubleshoot, list all containers using your preferred container runtimes CLI, e.g. docker.
Here is one example how you may list all Kubernetes containers running in docker:
	- 'docker ps -a | grep kube | grep -v pause'
	Once you have found the failing container, you can inspect its logs with:
	- 'docker logs CONTAINERID'
error execution phase wait-control-plane: couldn't initialize a Kubernetes cluster
To see the stack trace of this error execute with --v=5 or higher
```

不过`ps -ef | grep kubelet`, 其实`kubelet`进程已经启动了. 

查看`/var/log/message`日志, 发现 kubelet 倒没什么报错, 只是连接不上 apiserver, 因为 apiserver 根本没有启动, `/etc/kubernetes/manifests`下的 static pod 都没有启动.

不禁让人怀疑, 是不是 kubelet -> containerd 过程中是不是出了问题, 于是过滤出`containerd`的日志, 如下

```log
Dec 15 15:12:50 k8s-master-01 containerd: time="2022-12-15T15:12:50.168472282+08:00" Start cri plugin with config {...省略}

Dec 15 15:12:50 k8s-master-01 containerd: time="2022-12-15T15:12:50.168472282+08:00" level=info msg="Connect containerd service"
Dec 15 15:12:50 k8s-master-01 containerd: time="2022-12-15T15:12:50.168524630+08:00" level=info msg="Get image filesystem path \"/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs\""
Dec 15 15:12:50 k8s-master-01 containerd: time="2022-12-15T15:12:50.169026672+08:00" level=error msg="failed to load cni during init, please check CRI plugin status before setting up network for pods" error="cni config load failed: no network config found in /etc/cni/net.d: cni plugin not initialized: failed to load cni config"

Dec 15 15:14:08 k8s-master-01 containerd: time="2022-12-15T15:14:08.548569068+08:00" level=info msg="trying next host" error="failed to do request: Head \"https://asia-east1-docker.pkg.dev/v2/k8s-artifacts-prod/images/pause/manifests/3.6\": dial tcp 64.233.189.82:443: i/o timeout" host=registry.k8s.io
Dec 15 15:14:08 k8s-master-01 containerd: time="2022-12-15T15:14:08.549918017+08:00" level=error msg="RunPodSandbox for &PodSandboxMetadata{Name:etcd-k8s-master-01,Uid:2540ac6d2ab7eda00fbffa2863fbb465,Namespace:kube-system,Attempt:0,} failed, error" error="failed to get sandbox image \"registry.k8s.io/pause:3.6\": failed to pull image \"registry.k8s.io/pause:3.6\": failed to pull and unpack image \"registry.k8s.io/pause:3.6\": failed to resolve reference \"registry.k8s.io/pause:3.6\": failed to do request: Head \"https://asia-east1-docker.pkg.dev/v2/k8s-artifacts-prod/images/pause/manifests/3.6\": dial tcp 64.233.189.82:443: i/o timeout"
Dec 15 15:14:08 k8s-master-01 containerd: time="2022-12-15T15:14:08.838780081+08:00" level=info msg="trying next host" error="failed to do request: Head \"https://asia-east1-docker.pkg.dev/v2/k8s-artifacts-prod/images/pause/manifests/3.6\": dial tcp 64.233.189.82:443: i/o timeout" host=registry.k8s.io
```

看来是 pause 镜像无法下载.

kubeadm-config.yaml 中通过`imageRepository`定义了镜像源, 但无法对 pause 生效(kubelet有一个`--pod-infra-container-image`选项专门额外配置)

按照参考文章6, 7, 可以在 containerd 的配置文件(/etc/containerd/config.toml)中定义 pause 的镜像源

```yaml
## disabled_plugins = ["cri"]
[plugins]
  [plugins."io.containerd.grpc.v1.cri"]
    sandbox_image = "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6"
```

...但是, 重启 containerd 后, 重新 kubeadm init, 还是下载不下来, 而且好像根本就没生效. 重试了4, 5次, 我都绝望了.

#### 解决方法

后来, 按照参考文章6中提到的, 使用如下命令打印了一份全量的配置

```
containerd config default > /etc/containerd/config.toml
```

在这个全量的配置中, 修改`sandbox_image`的值, 再重启 containerd 服务, 然后 kubeadm init 就可以了.
