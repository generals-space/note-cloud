# kuber-secret对象

参考文章

1. [Kubernetes中的Secret配置](https://www.cnblogs.com/Leslieblog/p/10158429.html)

`secret`类型目前有3种: 

1. docker-registry(仓库配置)
2. generic(通用类型, 任意文件/目录/键值对) type为`Opaque`
3. tls(密钥配置) type为`kubernetes.io/tls`, 同时包括crt与key.

## generic

```

```

## tls

```
k create secret tls etcd-certs --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key
```
