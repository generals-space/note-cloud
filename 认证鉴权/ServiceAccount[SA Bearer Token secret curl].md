# ServiceAccount[SA Bearer Token secret curl]

参考文章

1. [官方文档 Managing Service Accounts](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/)
    - 介绍了用户账户(User)与服务账户(SA)的区别和联系(可以查看对应的中文文档, 可以理解更快):
        1. User是针对个人而言的, SA是针对运行在pod中的进程而言的.
        2. User是全局性的, 其名称在集群各ns中都是全局唯一的, 未来的用户资源不会做ns隔离, SA是ns隔离的.
        3. 实际场景中, User可能会与企业数据库保持同步(如ldap机制), 其使用对象为个人, 创建时需要特殊权限, 且可能涉及到复杂的业务流程; 而SA则是为了具体任务而存在的, ta的创建更轻量, 权限分配也更细致.
    - "Token Controller"部分介绍了SA与Token的关系: 每个SA可能有多个对应的Token

在kuber中, 每创建一个ns, 都会在该ns下默认创建一个名为`default`的SA对象.

```log
$ k get sa -n default
NAME      SECRETS   AGE
default   1         18d
```

每一个`SA`对象, 都拥有一个对应的`secret`对象.

```yaml
$ k get sa | grep default
default   1         16h
$ kya sa default
apiVersion: v1
kind: ServiceAccount
metadata:
  name: default
  namespace: kube-system
secrets:
- name: default-token-2xvwc
```

```yaml
$ k get secret | grep coredns
coredns-token-pz4dj    kubernetes.io/service-account-token    3    11h
$ kya secret coredns-token-pz4dj
apiVersion: v1
data:
  ca.crt: "base64字符串"
  namespace: base64("kube-system")
  token: "base64字符串"
kind: Secret
metadata:
  name: coredns-token-pz4dj
  namespace: kube-system
type: kubernetes.io/service-account-token
```

我们可以通过如下方式得到 token 原本的内容.

```
Token=$(echo -n "token的base64字符串内容" | base64 -d)
```

然后使用 curl , 在请求头中附上 token 内容, 就可以拥有该 SA 对象的相应权限了.

```
curl --header "Authorization: Bearer ${TOKEN}" -X GET -k https://localhost:6443/
```

## Pod 与 SA

另外, 每一个pod都**必须**设置SA, 如果部署文件中没有写明, 则默认使用ta所在ns的名为`default`的SA. 

```yaml
spec:
  serviceAccount: coredns
  serviceAccountName: coredns
```

pod在启动时都会挂载包含该`SA`相关信息的`volume`, 位置在pod内部的`/var/run/secrets/kubernetes.io/serviceaccount`.

```log
$ pwd
/var/run/secrets/kubernetes.io/serviceaccount
$ ll
total 0
lrwxrwxrwx 1 root root 13 Dec 10 16:33 ca.crt -> ..data/ca.crt
lrwxrwxrwx 1 root root 16 Dec 10 16:33 namespace -> ..data/namespace
lrwxrwxrwx 1 root root 12 Dec 10 16:33 token -> ..data/token
```

目录中的文件为对应secret对象中的数据字段, 且其中的内容已经经过base64解密, 无需额外操作.

这个`volume`中的信息将提供给pod中运行的程序以访问kuber API的权限.
