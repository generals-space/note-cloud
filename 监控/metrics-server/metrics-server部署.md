# 

参考文章

1. [从kubectl top看K8S监控](https://www.jianshu.com/p/64230e3b6e6c)
2. [kubernetes之配置Metrics Server](https://www.cnblogs.com/cptao/p/10912775.html)
    - `--kubelet-preferred-address-types=InternalIP`不加会出现"no such host"错误
3. [kubernetes监控组件Metrics-Server部署-v0.3.1最新版本](https://blog.csdn.net/zyl290760647/article/details/83041991)
    - `--kubelet-insecure-tls`不验证kubelet证书的合法性.
4. [no metrics known for pod](https://github.com/kubernetes-sigs/metrics-server/issues/237)
    - `no metrics known for pod`无法收集 Pod 指标

kuber 版本: 1.16.2 单节点集群, 宿主机配置 4C 8G.

`https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.7/components.yaml`

## FAQ

`metrics-server`在采集数据时会有延迟, 首次成功启动时需要几分钟, 在此期间使用`k top`查询, 可能会报`unable to fetch pod metrics for pod`.

### 无法收集 Pod 信息

```
I0724 03:06:01.186860       1 serving.go:312] Generated self-signed cert (/tmp/apiserver.crt, /tmp/apiserver.key)
I0724 03:06:01.941217       1 secure_serving.go:116] Serving securely on [::]:4443
E0724 03:06:20.070557       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/metrics-server-7cb45bbfd5-l42jm: no metrics known for pod
E0724 03:06:20.070607       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/kube-scheduler-k8s-master-01: no metrics known for pod
E0724 03:06:20.070611       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/coredns-67c766df46-jln7w: no metrics known for pod
E0724 03:06:20.070614       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/kube-controller-manager-k8s-master-01: no metrics known for pod
E0724 03:06:20.070617       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/kube-apiserver-k8s-master-01: no metrics known for pod
E0724 03:06:20.070620       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/coredns-67c766df46-rrw2w: no metrics known for pod
E0724 03:06:20.070623       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/kube-proxy-8x6ch: no metrics known for pod
E0724 03:06:20.070626       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/kube-flannel-ds-amd64-r7dft: no metrics known for pod
E0724 03:06:20.070629       1 reststorage.go:160] unable to fetch pod metrics for pod kube-system/etcd-k8s-master-01: no metrics known for pod
E0724 03:06:31.217028       1 reststorage.go:135] unable to fetch node metrics for node "k8s-master-01": no metrics known for node
E0724 03:07:01.961412       1 manager.go:111] unable to fully collect metrics: unable to fully scrape metrics from source 
```

参考文章4的官方 issue 中, 名为`philipakash`的用户给出了奇妙的方法.

原本官方 readme 中给出的`v0.3.7`的 yaml 文件中启动参数如下

```yaml
        args:
          - --cert-dir=/tmp
          - --secure-port=4443
```

我按照ta的方法, 添加上了`/metrics-server`作为参数, 把`args`修改成了

```yaml
        args:
          - /metrics-server
          - --cert-dir=/tmp
          - --secure-port=4443
          - --kubelet-preferred-address-types=InternalIP
          - --kubelet-insecure-tls
```

竟然可以了...

当然, 把`args`换成`command`也是可以的...

```yaml
        command:
          - /metrics-server
          - --cert-dir=/tmp
          - --secure-port=4443
          - --kubelet-preferred-address-types=InternalIP
          - --kubelet-insecure-tls
```

真不知道该说啥...

### 

```log
$ k logs -f metrics-server-7cb45bbfd5-l42jm
kubelet_summary:k8s-master-01: unable to fetch metrics from Kubelet k8s-master-01 (k8s-master-01): Get https://k8s-master-01:10250/stats/summary?only_cpu_and_memory=true: dial tcp: lookup k8s-master-01 on 10.96.0.10:53: no such host
```

`--kubelet-preferred-address-types=InternalIP`

`coredns`无法根据宿主机节点的hostname得到真实IP, 需要在启动参数中加上这个选项.

### 

```
E0724 03:23:50.470647       1 manager.go:111] unable to fully collect metrics: unable to fully scrape metrics from source kubelet_summary:k8s-master-01: unable to fetch metrics from Kubelet k8s-master-01 (192.168.80.10): Get https://192.168.80.10:10250/stats/summary?only_cpu_and_memory=true: x509: cannot validate certificate for 192.168.80.10 because it doesn't contain any IP SANs
```

`--kubelet-insecure-tls`

kubelet的 10250 是 https 端口, 默认 metric-server 会验证ta的合法性, 但是我在使用 kubeadm 创建集群时有指定 SANs 值, 所以这里会出错.

添加上这个参数就可以了, 类似于`curl`的`-k`选项, 不检测服务端证书.
