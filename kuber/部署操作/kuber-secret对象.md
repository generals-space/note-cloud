# kuber-secret对象

参考文章

1. [Kubernetes中的Secret配置](https://www.cnblogs.com/Leslieblog/p/10158429.html)
    - 对于`generic`类型的`Secret`资源, base64的编码和解码都是手动完成的...有什么意义???
    - 没有讲`tls`类型的`Secret`的使用方法.

`secret`类型目前有3种: 

1. docker-registry(仓库配置)
2. generic(通用类型, 任意文件/目录/键值对) type为`Opaque`
3. tls(密钥配置) type为`kubernetes.io/tls`, 同时包括`crt`与`key`.

## docker-registry

不多说, 私有仓库的认证信息, 使用场景固定.

## generic

这个类型的`Secret`资源与`ConfigMap`几乎完全一致, 使用方法也大致相同. 

唯一区别应该是, `kubectl describe`一个`Secret`资源时, 并不会打印其中的内容吧? 不过感觉没什么用, `kubectl get`仍然可以通过`-o yaml`得到其中存储的信息...

```
kubectl create secret generic db-user-pass --from-file=./username.txt --from-file=./password.txt
```

## tls

```
k create secret tls etcd-certs --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key
```
