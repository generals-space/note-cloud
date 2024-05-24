参考文章

1. [Certificate Signing Requests](https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/#normal-user)
    - 官方文档
2. [证书签名请求](https://kubernetes.io/zh-cn/docs/reference/access-authn-authz/certificate-signing-requests/#normal-user)
    - 官方文档, 中文版

在kube集群的`RBAC`框架下, 可以通过创建`ServiceAccount`, 实现赋权及鉴权. 认证凭据为token(Bearer Token)+ca.crt(SA的Secret中有保存), 对应到curl底层命令, 为

```bash
curl -k -H 'Content-Type: application/json' -H 'Authorization: Bearer xxxxxx' 'https://127.0.0.1:6443/api/v1/namespaces/default/pods'
```

权限的主体还有其他选择: `User`及`Group`, 基于ca.crt, user.key与user.crt的证书链机制, 对应到curl命令是

```bash
curl -k -H 'Content-Type:application/json' --cacert /etc/kubernetes/pki/ca.crt --cert user.crt --key user.key 'https://127.0.0.1:6443/api/v1/namespaces/default/pods'
```

从目前来看, SA的认证机制都是kube的应用层在用, 而ca.crt的认证机制都是kube原生组件在用???

> 不过 Bearer Token 底层仍然是通过 x509 证书来实现的, 具体逻辑要看 [apiserver]().

创建SA, 然后通过RBAC赋权的方式我们已经掌握了, 本文讲的是如何创建一个合法用户, 然后通过RBAC为该用户赋权.

## 

创建新用户要求先准备好一个key, 以及csr文件, 见参考文章1, 2.

```log
$ openssl genrsa -out myuser.key 2048
$ openssl req -new -key myuser.key -out myuser.csr
Country Name (2 letter code) [AU]:
State or Province Name (full name) [Some-State]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]: ## 这里应该是组名
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:60099@internal.users ## 这里是 CN, 最重要.
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
```

> 设置 CSR 的 CN 和 O 属性很重要。CN 是用户名, O 是该用户归属的组。

然后创建 csr 资源对象, 让 kube 集群使用 ca.crt 完成签发.

```yaml
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: myuser
spec:
  ## request 字段是 CSR 文件内容的 base64 编码值. 要得到该值, 可以执行命令
  ## cat myuser.csr | base64 | tr -d "\n"
  request: 
  signerName: kubernetes.io/kube-apiserver-client
  expirationSeconds: 86400 # one day
  ## usages 字段必须是 'client auth'
  usages:
  - client auth
```

```log
$ k apply -f csr.yaml 
certificatesigningrequest.certificates.k8s.io/myuser created
$ k get csr
NAME     AGE   SIGNERNAME                            REQUESTOR          REQUESTEDDURATION   CONDITION
myuser   3s    kubernetes.io/kube-apiserver-client   kubernetes-admin   24h                 Pending
```

常规证书的申请-签发流程, 其实可以直接使用`openssl`通过`ca.crt`自行完成, 不过这是基于操作者自身是管理员的情况的. 

这里的csr签发流程中, 添加了一个"审批"的流程, 普通用户通过创建`csr`资源提交给管理员进行审批, 上面`Pending`表示该csr还未完成审批.

------

接下来, 管理员可以使用如下命令进行批准.

```
kubectl certificate approve myuser
```

批准完成后, 新用户的证书可以在该csr资源的`.status.certificate`字段中获取, 该csr资源的状态也会发生变化.

```log
$ k get csr
NAME     AGE     SIGNERNAME                            REQUESTOR          REQUESTEDDURATION   CONDITION
myuser   7m58s   kubernetes.io/kube-apiserver-client   kubernetes-admin   24h                 Approved,Issued
```

然后就是通过RBAC进行赋权了.

------

其实kube的签发, 底层也是 openssl 那一套, 我们也可以通过类似如下的命令去验证证书链

```log
$ openssl verify -CAfile /etc/kubernetes/pki/ca.crt ./myuser.crt 
./myuser.crt: OK
```

## 关于 signerName: kubernetes.io/kube-apiserver-client

见参考文章2.

kube-controller-manager 内建了一个证书批准者, 其`signerName`为`kubernetes.io/kube-apiserver-client-kubelet`, 该批准者将 CSR 上用于节点凭据的各种权限委托给权威认证机构.

kube-controller-manager 将 SubjectAccessReview 资源发送（POST）到 apiserver, 以便检验批准证书的授权.

