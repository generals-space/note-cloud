# kubeadm join-集群中加入新的master与worker节点

参考文章

1. [Discovering what cluster CA to trust](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-join/#discovering-what-cluster-ca-to-trust)
2. [Steps for the first control plane node](https://kubernetes.io/docs/setup/independent/high-availability/#steps-for-the-first-control-plane-node)
    - 其中的note是重点

在 kubeadm 初始化集群成功后会返回 join 命令, 里面有 token、discovery-token-ca-cert-hash等参数, 但ta们是有过期时间的. token 的过期时间是24小时, certificate-key 过期时间是2小时.

```log
$ kubeadm token list
TOKEN                     TTL         EXPIRES                     USAGES                   DESCRIPTION   EXTRA GROUPS
1ugnpn.7qwbt58paanfov8a   <invalid>   2019-07-19T11:17:01+08:00   authentication,signing   <none>        system:bootstrappers:kubeadm:default-node-token
```

如下生成worker节点的join命令

```log
$ kubeadm token create --print-join-command
kubeadm join k8s-master-7-13:8443 --token fw6ywo.1sfp61ddwlg1we27     --discovery-token-ca-cert-hash sha256:52cab6e89be9881e2e423149ecb00e610619ba0fd85f2eccc3137adffa77bb04
```

> `kubeadm token create --ttl 0 --print-join-command`可以创建一个永不过期的token.

要添加master节点, 还要执行如下命令, 得到`certificate-key`

```log
## --experimental-upload-certs已被废弃
## kubeadm init phase upload-certs --experimental-upload-certs
$ kubeadm init phase upload-certs --upload-certs
[upload-certs] Storing the certificates in Secret "kubeadm-certs" in the "kube-system" Namespace
[upload-certs] Using certificate key:
70f399e275cabef0bb2794ea76303da0220574f59e994e755378a359edb5a233
```

使用第一步生成的worker的join命令, 加上上面生成的`certificate-key`, 拼接起来组成master的join命令.

```log
kubeadm join k8s-master-7-13:8443 --token fw6ywo.1sfp61ddwlg1we27     --discovery-token-ca-cert-hash sha256:52cab6e89be9881e2e423149ecb00e610619ba0fd85f2eccc3137adffa77bb04 --control-plane --certificate-key 70f399e275cabef0bb2794ea76303da0220574f59e994e755378a359edb5a233
```

> 注意: 新加的master节点注意拷贝ControlPlane的`/etc/kubernetes`目录下的各个证书.

## 问题处理

### 

```log
$ kubeadm token list
failed to create API client configuration from kubeconfig: invalid configuration: [unable to read client-cert certs.d/kuber.centos7/client.crt for kuber-admin due to open certs.d/kuber.centos7/client.crt: no such file or directory, unable to read client-key certs.d/kuber.centos7/client.key for kuber-admin due to open certs.d/kuber.centos7/client.key: no such file or directory, unable to read certificate-authority certs.d/kuber.centos7/ca.crt for kuber due to open certs.d/kuber.centos7/ca.crt: no such file or directory]
To see the stack trace of this error execute with --v=5 or higher
```

#### 场景描述

单节点集群想添加两个worker节点, 按照上述步骤执行发现出错了. 看这报错是因为我把`kubectl`的配置文件更改了的缘故. 我把base64加密的证书字符串放到单独的文件中了, 配置文件类似如下

```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority: certs.d/kuber.centos7/ca.crt
    server: https://kube-apiserver.generals.space:8443
  name: kuber
contexts:
- context:
    cluster: kuber
    namespace: default
    user: kuber-admin
  name: kube-def
current-context: kube-def
kind: Config
preferences: {}
users:
- name: kuber-admin
  user:
    client-certificate: certs.d/kuber.centos7/client.crt
    client-key: certs.d/kuber.centos7/client.key
```

而且上面的报错是因为在kubectl配置文件中的证书路径写的是相对路径, ta找不到这些证书.

#### 解决办法

要么把证书路径写成绝对路径, 要么在`~/.kube/`目录下执行`kubeadm命令`.

###

主节点上执行如下命令没有输出

```log
[root@k8s-master-01 .kube]# kubeadm token create --print-join-command
W0613 13:01:10.666223    4935 validation.go:28] Cannot validate kube-proxy config - no validator is available
W0613 13:01:10.666324    4935 validation.go:28] Cannot validate kubelet config - no validator is available
```

看了看nginx日志, 没啥问题, `kubectl`的命令也都正常执行, 不过`kubeadm`好像有点问题

```log
[root@k8s-master-01 .kube]# k version
Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.2", GitCommit:"59603c6e503c87169aea6106f57b9f242f64df89", GitTreeState:"clean", BuildDate:"2020-01-18T23:30:10Z", GoVersion:"go1.13.5", Compiler:"gc", Platform:"linux/amd64"}
Unable to connect to the server: EOF
```

```log
[root@k8s-master-01 .kube]# kubeadm token list
failed to list bootstrap tokens: Get https://kube-apiserver.generals.space:8443/api/v1/namespaces/kube-system/secrets?fieldSelector=type%3Dbootstrap.kubernetes.io%2Ftoken: EOF
To see the stack trace of this error execute with --v=5 or higher
```

加上`--v=5`参数也没看出什么.

后来重启了一下主节点, 可以了...
