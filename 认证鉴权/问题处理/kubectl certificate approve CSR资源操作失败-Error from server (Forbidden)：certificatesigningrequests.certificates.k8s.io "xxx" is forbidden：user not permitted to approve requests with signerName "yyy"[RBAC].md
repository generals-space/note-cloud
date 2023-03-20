# kubectl certificate approve CSR资源操作失败-Error from server (Forbidden)：certificatesigningrequests.certificates.k8s.io "xxx" is forbidden：user not permitted to approve requests with signerName "yyy"

参考文章

1. [CSR permissions error on Kubernetes 1.18 with RBAC enabled](https://github.com/cloudfoundry-incubator/quarks-operator/issues/891)

kube: 1.22.4

```log
$ kubectl certificate approve ${csrName}
Error from server (Forbidden)：certificatesigningrequests.certificates.k8s.io "xxx" is forbidden：user not permitted to approve requests with signerName "yyy"
```

解决办法, 新增并绑定如下权限

```yaml
- apiGroups: ["certificates.k8s.io"]
  resources: ["signers"]
  verbs: ["approve"]
```
