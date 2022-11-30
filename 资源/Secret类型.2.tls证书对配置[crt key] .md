# Secret类型.2.tls证书对配置[crt key] 

参考文章

1. [Kubernetes中的Secret配置](https://www.cnblogs.com/Leslieblog/p/10158429.html)
    - 对于`generic`类型的`Secret`资源, base64的编码和解码都是手动完成的...有什么意义???
    - 没有讲`tls`类型的`Secret`的使用方法.
2. [Kubernetes-token（四）](https://www.jianshu.com/p/1c188189678c)

```
k create secret tls etcd-certs --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key
```

这种Secret的type同样有问题, 叫`kubernetes.io/tls`.

```console
$ k get secret
NAME          TYPE                 DATA   AGE
etcd-certs    kubernetes.io/tls    2      4s
```

```yaml
## kya secret etcd-certs
apiVersion: v1
data:
  tls.crt: crt证书内容(base64加密)
  tls.key: key内容(base64加密)
kind: Secret
metadata:
  name: etcd-certs
  namespace: zjjpt-es
type: kubernetes.io/tls
```
