参考文章

1. [使用启动引导令牌（Bootstrap Tokens）认证](https://kubernetes.io/zh-cn/docs/reference/access-authn-authz/bootstrap-tokens/)
    - bootstrap token 必须符合正则表达式 `[a-z0-9]{6}\.[a-z0-9]{16}`
2. [Kubernetes - kubelet bootstrap 流程](https://lingxiankong.github.io/2018-09-18-kubelet-bootstrap-process.html)
    - 向集群新增节点时, kubelet 使用低权限的 bootstrap token 跟 api server 建立连接后, 要能够自动向 api server 申请自己的证书, 并且 api server 要能够自动审批证书.
3. [kubeadm 工作原理 -kubeadm init 原理分析 -kubeadm join 原理分析](https://xie.infoq.cn/article/040212f39e0bcdac15f40f9b6)

## 介绍

bootstrap token 是在 kubeadm init 进行集群初始化时, 保存在 kube-system/bootstrap-token-${xxxxxx} 的一个 Secret 对象.

```yaml
## kubectl -n kube-system get secret bootstrap-token-fl6bwd -oyaml
apiVersion: v1
data:
  ## system:bootstrappers:kubeadm:default-node-token
  auth-extra-groups: c3lzdGVtOmJvb3RzdHJhcHBlcnM6a3ViZWFkbTpkZWZhdWx0LW5vZGUtdG9rZW4=
  expiration: MjAyMi0xMi0xMVQxNjo1Mjo1NSswODowMA==
  token-id: Zmw2Yndk
  token-secret: M3ViZTF1d2N5azNhNW82Yg==
  usage-bootstrap-authentication: dHJ1ZQ==  ## true
  usage-bootstrap-signing: dHJ1ZQ==         ## true
kind: Secret
metadata:
  name: bootstrap-token-fl6bwd ## 最后6个字符是 token-id 的值
  namespace: kube-system
type: bootstrap.kubernetes.io/token
```

1. type为`bootstrap.kubernetes.io/token`
2. name中, `bootstrap-token-fl6bwd`, 后6位字符是 token-id 的值

kubeadm init 完成后, 会打印 join 子命令, 类似如下

```
kubeadm join k8s-master-7-13:8443 --token fw6ywo.1sfp61ddwlg1we27 \ 
--discovery-token-ca-cert-hash sha256:52cab6e89be9881e2e423149ecb00e610619ba0fd85f2eccc3137adffa77bb04 \
--control-plane
```

这里的`--token`就是 bootstrap token 的内容.

## 有什么用?

bootstrap token 也是 token, 使用方法跟常规 bearer token 一样, 如下

```json
// curl -k -H 'Content-Type: application/json' -H 'Authorization: Bearer fl6bwd.3ube1uwcyk3a5o6b' 'https://127.0.0.1:6443/api/v1/namespaces/kube-system/pods'
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {
    
  },
  "status": "Failure",
  "message": "pods is forbidden: User \"system:bootstrap:fl6bwd\" cannot list resource \"pods\" in API group \"\" in the namespace \"kube-system\"",
  "reason": "Forbidden",
  "details": {
    "kind": "pods"
  },
  "code": 403
}
```

认证通过了, 不过没权限, 毕竟 bootstrap token 不是用来发这种请求的. 按照参考文章3中所说, bootstrap token 的权限是非常低的, 只能让客户端请求通过"认证"阶段, 只保留"authentication", "signing"权限.
