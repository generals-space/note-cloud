# control plane升级失败-service-cluster-ip-range is too large

参考文章

1. [kubeadm init --service-cidr= validations are still broken for IPv6 subnets](https://github.com/kubernetes/kubeadm/issues/2132)
2. [Kubeadm DualStack Support for List of Service IPs](https://github.com/kubernetes/kubernetes/pull/82473/files)
  - 单元测试代码中给出了合法的 IPv6 网段示例, 掩码位为 112.

```log
$ kubeadm upgrade apply --config ./kubeadm-config.v1.17.2.yaml 
W0605 11:40:18.853549   50129 validation.go:28] Cannot validate kube-proxy config - no validator is available
W0605 11:40:18.853598   50129 validation.go:28] Cannot validate kubelet config - no validator is available
[upgrade/config] Making sure the configuration is correct:
W0605 11:40:18.874401   50129 common.go:94] WARNING: Usage of the --config flag for reconfiguring the cluster during upgrade is not recommended!
W0605 11:40:18.875702   50129 validation.go:28] Cannot validate kube-proxy config - no validator is available
W0605 11:40:18.875717   50129 validation.go:28] Cannot validate kubelet config - no validator is available
[preflight] Running pre-flight checks.
...省略
[upgrade/staticpods] Preparing for "kube-apiserver" upgrade
...省略
[upgrade/staticpods] This might take a minute or longer depending on the component/version gap (timeout 5m0s)
Static pod: kube-apiserver-k8s-master-01 hash: 8810a8c8d7f86ec2518c4cae06f33257
[upgrade/apply] FATAL: couldn't upgrade control plane. kubeadm has tried to recover everything into the earlier state. Errors faced: timed out waiting for the condition
To see the stack trace of this error execute with --v=5 or higher
```

开了日志

```log
$ kubeadm upgrade apply --config ./kubeadm-config.v1.17.2.yaml -v 5
[upgrade/staticpods] This might take a minute or longer depending on the component/version gap (timeout 5m0s)
Static pod: kube-apiserver-k8s-master-01 hash: 8810a8c8d7f86ec2518c4cae06f33257
I0605 11:49:15.814611   71811 request.go:848] Got a Retry-After 1s response for attempt 1 to https://kube-apiserver.generals.space:8443/api/v1/namespaces/kube-system/pods/kube-apiserver-k8s-master-01?timeout=10s
I0605 11:49:25.312877   71811 request.go:848] Got a Retry-After 1s response for attempt 1 to https://kube-apiserver.generals.space:8443/api/v1/namespaces/kube-system/pods/kube-apiserver-k8s-master-01?timeout=10s
...省略 n 次
timed out waiting for the condition
couldn't upgrade control plane. kubeadm has tried to recover everything into the earlier state. Errors faced
...省略 error
```

好在升级时用的新版本 apiserver 的 docker 容器还在, 于是查看 ta 的日志.

```log
$ d logs k8s_kube-apiserver_kube-apiserver-k8s-master-01_kube-system_b46fbe29509e9b523b237db7c4fe7c6f_5
Flag --insecure-port has been deprecated, This flag will be removed in a future version.
I0605 05:41:43.755378       1 server.go:596] external host was not specified, using 192.168.80.121
Error: specified --secondary-service-cluster-ip-range is too large
Usage:
  kube-apiserver [flags]
...省略
```

我找了老半天, 都没找到 apiserver `--secondary-service-cluster-ip-range` 这个参数在哪, 后来在[源码](https://github.com/kubernetes/kubernetes/blob/v1.17.2/cmd/kube-apiserver/app/options/validation.go#L56)里找到了, 只有在开启`IPv6DualStack`特性才有效.

相关问题中没有找到解决方法, 只在参考文章2的 pull request 的单元测试中, 发现了合法的 IPv6 网段...掩码位112(之前从24, 48, 64, 一直试到 96, 没到112这么大)

```yaml
controllerManager:
  extraArgs:
    cluster-cidr: 10.254.0.0/16,2019:20::/96
    ## service-cluster-ip-range: 10.96.0.0/12,fec0:30::/48
    service-cluster-ip-range: 10.96.0.0/12,fec0:30::/112
networking:
  dnsDomain: cluster.local
  ## 不知道是不是 kubeadm 还是 apiserver 的限制, serviceSubnet 的掩码位小于 96 都会出错:
  ## specified --secondary-service-cluster-ip-range is too large.
  ## 目前来说, 112 是比较合适的.
  podSubnet: 10.254.0.0/16,2019:20::/96
  ## serviceSubnet: 10.96.0.0/12,fec0:30::/48
  serviceSubnet: 10.96.0.0/12,fec0:30::/112
```
