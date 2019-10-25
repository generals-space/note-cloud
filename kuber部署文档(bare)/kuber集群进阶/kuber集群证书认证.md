# kuber集群证书认证

参考文章

1. [Kubernetes安装之证书验证](https://www.kubernetes.org.cn/1861.html)

2. [和我一步步部署 kubernetes 集群](https://github.com/opsnull/follow-me-install-kubernetes-cluster)

3. [Kubernetes集群安全配置案例](http://www.cnblogs.com/breg/p/5923604.html)

kubernetes 系统的各组件需要使用 TLS 证书对通信进行加密. apiserver与etcd之间, apiserver与kubectl/controller-manager之间等. 本文介绍如何使用openssl工具为kuber集群生成自签名证书.

先来说一下证书认证的工作流程. 在为普通网站颁发自签名https证书时, 只需要生成`.key`私钥文件与`.crt`证书文件就行, 但kubernetes使用的认证方式称为**双向认证**.

前者只需要客户端信任服务器的https证书即可, 而后者不只要apiserver信任与之建立连接的客户端的证书, 还需要连接的客户端信任apiserver本身的证书. 后者的认证流程被称为**双向认证**, 毕竟前者只需要保护客户端免于威胁就行了, 而后者, 双方都很重要.


其实双向认证原理与普通https站的证书相似. 

传统的认证方式是, 第三方认证机构CA事先将自己生成的根证书(这个证书不再需要被信任人签名)发送给各浏览器厂商, 然后这些第三方机构再使用自己的根证书为前来申请https证书的站长签名, 于是访问者在使用这些浏览器访问这些网站的时候就能拿到对方的证书, 浏览器再与自己存储的根证书库对其进行校验, 如果校验失败则会提示用户该网站的证书不合法(给其签名的机构不被浏览器信任).

双向认证的实现方法是, 首先生成一个根证书, 用这个根证书分别为认证双方颁发证书和私钥, 然后双方都使用这个根证书验证连接对端的证书...

## 1. 证书生成

首先生成根证书, 其根证书就是普通的https证书, 不过不会直接当作普通证书来用, 而是用来给其他证书作签名的.

```
$ openssl genrsa -out ca.key 2048
## -subj选项填你自己指定的域名就可以了, 随便指, 之后如何使用与它无关
$ openssl req -x509 -new -nodes -key ca.key -subj '/CN=sky-mobi.com' -days 5000 -out ca.crt
```

然后我们得到`ca.key`与`ca.crt`

我们首先尝试为apiserver一对证书和私钥.

```
$ openssl genrsa -out apiserver.key 2048
## 这里的-subj选项将用于客户端连接, 而且不能是IP, 必须是域名, 如果两者不匹配则认证失败. 所以客户端需要添加hosts域名映射.
$ openssl req -new -key apiserver.key -subj "/CN=apiserver.sky-mobi.com" -out apiserver.csr
$ openssl x509 -req -in apiserver.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 5000 -out apiserver.crt
```

然后是kubectl, 步骤完全一样.

```
$ openssl genrsa -out kubectl.key 2048
## 这里的-subj选项没什么用, 好像可以随便写?
$ openssl req -new -key kubectl.key -subj "/CN=kubectl所在的IP或域名" -out kubectl.csr
$ openssl x509 -req -in kubectl.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 5000 -out kubectl.crt
```

------

配置apiserver系统服务, 保证有如下参数

```
[Service]
ExecStart=/usr/local/kubernetes/bin/kube-apiserver  \
    --bind-address=0.0.0.0 \
    --secure-port=6443 \
    --client-ca-file=/usr/local/kubernetes/etc/ca.crt \                 ## 根证书路径 
    --tls-cert-file=/usr/local/kubernetes/etc/apiserver.crt \           ## apiserver证书路径 
    --tls-private-key-file=/usr/local/kubernetes/etc/apiserver.key \    ## apiserver私钥路径
    ...
```

重启apiserver服务.

然后配置kubectl.

```
$ kubectl config set-cluster sky-test --server=https://apiserver.sky-mobi.com:6443 --certificate-authority=/usr/local/kubernetes/etc/ca.crt --embed-certs=true
$ kubectl config set-credentials admin --certificate-authority=/usr/local/kubernetes/etc/ca.crt --client-key=/root/.kube/kubectl.key --client-certificate=/root/.kube/kubectl.crt
$ kubectl config set-context sky-test --cluster=sky-test --user=admin
$ kubectl config use-context sky-test
```

记得在kubectl所在主机上添加apiserver的hosts映射.

kubectl客户端的证书和私钥中, `-subj`选项的CN值真是可以随便指定的, 也因此客户端的证书与私钥可以拷贝到其他机器上使用.

...怎么感觉这种双向认证这么水呢???