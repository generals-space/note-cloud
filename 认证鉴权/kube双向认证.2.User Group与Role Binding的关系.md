# kuber集群证书认证

参考文章

1. [Users in Kubernetes](https://v1-21.docs.kubernetes.io/docs/reference/access-authn-authz/authentication/)
    - 官方文档

## 1. 引言

**User**

> In this regard, Kubernetes does not have objects which represent normal user accounts. Normal users cannot be added to a cluster through an API call.
>
> k8s自身没有常规意义上的`User`对象, 也不可能通过API去创建/删除.

**Group**

> As of Kubernetes 1.4, client certificates can also indicate a user's group memberships using the certificate's organization fields.
> 
> 从 1.4 版本开始, 客户端证书中就可以表示证书拥有者所属的Group分组了.

实际上, k8s在RBAC认证系统中, 并不存在`User`和`Group`的用户类型.

```console
$ k api-resources | grep rbac
clusterrolebindings    rbac.authorization.k8s.io    false    ClusterRoleBinding
clusterroles           rbac.authorization.k8s.io    false    ClusterRole
rolebindings           rbac.authorization.k8s.io    true     RoleBinding
roles                  rbac.authorization.k8s.io    true     Role
```

但是, 在某些`RoleBinding`和`ClusterRoleBinding`中, 的确能见到`User`和`Group`对象.

```yaml
## kubectl get clusterrolebindings system:kube-scheduler -oyaml
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: system:kube-scheduler
```

```yaml
## kubectl get clusterrolebindings cluster-admin -oyaml
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:masters
```

ta们怎么来的?🤔

## 系统组件与apiserver通信认证

我们知道, 在`/etc/kubernetes/pki/`目录下, 存在多个证书. 但是实际上, 系统组件与apiserver通信并不依赖这些证书, 以`scheduler`为例.

```yaml
## cat /etc/kubernetes/manifests/kube-scheduler.yaml
spec:
  containers:
  - command:
    - kube-scheduler
    - --kubeconfig=/etc/kubernetes/scheduler.conf
```

`scheduler`, `controller-manager`, `kubelet`都通过`--kubeconfig`参数指定yaml文件进行认证的.

```log
$ ll /etc/kubernetes/
-rw------- 1 root root 5453 Sep 24 13:02 admin.conf
-rw------- 1 root root 5489 Sep 24 13:02 controller-manager.conf
-rw------- 1 root root 5497 Sep 24 13:02 kubelet.conf
-rw------- 1 root root 5437 Sep 24 13:02 scheduler.conf
```

查看`kubelet.conf`的内容.

```yaml
## cat /etc/kubernetes/kubelet.conf
apiVersion: v1
clusters:
- cluster:
    ## 这个字段与 pki 目录下的 ca.crt 内容是相同的.
    certificate-authority-data: base64(/etc/kubernetes/pki/ca.crt)
    server: https://k8s-server-lb:8443
  name: kubernetes
kind: Config
users:
- name: system:node:k8s-master-01
  user:
    client-certificate-data: 
      subject= /O=system:nodes/CN=system:node:k8s-master-01
      issuer= /CN=kubernetes
    client-key-data: 
```

对`kubelet.conf`中的`client-certificate-data`字段进行`base64`解码, 然后保存为`kubelet.crt`, 并使用`openssl`命令查看证书信息.

```log
$ openssl x509 -noout -subject -issuer -in kubelet.crt
subject= /O=system:nodes/CN=system:node:k8s-master-01
issuer= /CN=kubernetes
```

- subject: 该证书的所有人信息
- issuer: 为该证书签名的父级证书所有人的信息

双向认证时, 通信双方都可以通过`/etc/kubernetes/pki/ca.crt`对对端证书进行验证, 因为所有组件所用的证书, 都是由该证书签发的. 

认证方式如下

```log
$ openssl verify -CAfile /etc/kubernetes/pki/ca.crt /etc/kubernetes/kubelet.crt
/etc/kubernetes/kubelet.crt: OK
```

------

其余`scheduler`, `controller-manager`, 以及默认kubectl使用的`admin.conf`, 也都是这个套路.

开发者在编写operator代码时, 一般是通过为所在容器指定`ServiceAccount`, 然后为该SA对象绑定权限, 也可以通过`--kubeconfig`加载类似的配置.

## 基于证书的用户层级认证

上面我们通过`openssl`查看到了kubelet组件所使用的证书的信息, 其实`subject`信息中`/O`就表示`Group`, `/CN`则对应`User`.

下面是各系统组件所使用的证书信息列表:

| conponent          | group          | user                           |
| :----------------- | :------------- | :----------------------------- |
| admin.conf         | system:masters | kubernetes-admin               |
| scheduler          |                | system:kube-scheduler          |
| controller-manager |                | system:kube-controller-manager |
| kubelet            | system:nodes   | system:node:k8s-master-01      |

user/group一般在`ClusterRoleBindings`中直接绑定, 可以使用如下命令进行查询.

```bash
kya clusterrolebindings | grep 'system:masters' -B30
```

> `-B`数值可以根据实际情况调整.

上面`admin.conf`与`kubelet`虽然有对应的`User`主体, 但实际上只有`Group`的绑定.

```console
$ kya clusterrolebindings | grep 'kubernetes-admin'
## 无输出
$ kya clusterrolebindings | grep 'system:node:k8s-master-01'
## 无输出
```
