# coredns无法启动-certificate signed by unknown  authority

kuber版本: 1.17.2

本文中说的问题, 并非是类似如下这种, 由于 kubectl 中的密钥未正确配置导致与 apiserver 认证失败的问题.

```log
Unable to connect to the server: x509: certificate signed by unknown authority (possibly because of "crypto/rsa: verification error" while trying to verify candidate au
thority certificate "kubernetes")
```

我能确定 setup 集群完毕后, 安装好了 flannel 插件, 理论上此时 coredns 应该可以正常启动了, 但是使用 kubectl 查看时仍是处于`ContainerCreating`状态, describe 了一下, 有如下输出


```log
Events:
  Type     Reason                  Age                 From                    Message
  ----     ------                  ----                ----                    -------
  Normal   Scheduled               119s                default-scheduler       Successfully assigned kube-system/coredns-7f9c544f75-hm7fj to k8s-worker-02
  Warning  FailedCreatePodSandBox  117s                kubelet, k8s-worker-02  Failed to create pod sandbox: rpc error: code = Unknown desc = [failed to set up sandbox
container "1afc8337bbf1b93c5eb240966b2a48aca5b119d22bc535d9d03790d658356e0c" network for pod "coredns-7f9c544f75-hm7fj": networkPlugin cni failed to set up pod "coredns
-7f9c544f75-hm7fj_kube-system" network: error getting ClusterInformation: Get https://[10.96.0.1]:443/apis/crd.projectcalico.org/v1/clusterinformations/default: x509: c
ertificate signed by unknown authority (possibly because of "crypto/rsa: verification error" while trying to verify candidate authority certificate "kubernetes"), faile
d to clean up sandbox container "1afc8337bbf1b93c5eb240966b2a48aca5b119d22bc535d9d03790d658356e0c" network for pod "coredns-7f9c544f75-hm7fj": networkPlugin cni failed
to teardown pod "coredns-7f9c544f75-hm7fj_kube-system" network: error getting ClusterInformation: Get https://[10.96.0.1]:443/apis/crd.projectcalico.org/v1/clusterinfor
mations/default: x509: certificate signed by unknown authority (possibly because of "crypto/rsa: verification error" while trying to verify candidate authority certific
ate "kubernetes")]
  Normal   SandboxChanged          4s (x10 over 117s)  kubelet, k8s-worker-02  Pod sandbox changed, it will be killed and re-created.
```

`/var/log/message`中也有 kubelet 的相关日志.

```log
Jun 20 23:59:19 k8s-master-01 kubelet: E0620 23:59:19.340896  105799 kuberuntime_manager.go:898] Failed to stop sandbox {"docker" "4ab054e3be9f5b3c6b1f440ee8d0e7e1370b2
28b198573a5b69b071e90707b65"}
Jun 20 23:59:19 k8s-master-01 kubelet: E0620 23:59:19.340936  105799 kuberuntime_manager.go:676] killPodWithSyncResult failed: failed to "KillPodSandbox" for "90811698-
2454-44cf-9d92-2b4bece37207" with KillPodSandboxError: "rpc error: code = Unknown desc = networkPlugin cni failed to teardown pod \"coredns-7f9c544f75-8thxc_kube-system
\" network: error getting ClusterInformation: Get https://[10.96.0.1]:443/apis/crd.projectcalico.org/v1/clusterinformations/default: x509: certificate signed by unknown
 authority (possibly because of \"crypto/rsa: verification error\" while trying to verify candidate authority certificate \"kubernetes\")"
Jun 20 23:59:19 k8s-master-01 kubelet: E0620 23:59:19.340957  105799 pod_workers.go:191] Error syncing pod 90811698-2454-44cf-9d92-2b4bece37207 ("coredns-7f9c544f75-8th
xc_kube-system(90811698-2454-44cf-9d92-2b4bece37207)"), skipping: failed to "KillPodSandbox" for "90811698-2454-44cf-9d92-2b4bece37207" with KillPodSandboxError: "rpc e
rror: code = Unknown desc = networkPlugin cni failed to teardown pod \"coredns-7f9c544f75-8thxc_kube-system\" network: error getting ClusterInformation: Get https://[10
.96.0.1]:443/apis/crd.projectcalico.org/v1/clusterinformations/default: x509: certificate signed by unknown authority (possibly because of \"crypto/rsa: verification er
ror\" while trying to verify candidate authority certificate \"kubernetes\")"
```

最开始没注意日志中打印的 url 包含了 projectcalico, 后来才发现是 coredns 所在节点的 `/etc/cni/net.d`目录下还残留着上次部署 calico 插件生成的配置文件, 将所有节点上的 calico 配置文件删除后再重新 setup 一遍就可以了...
