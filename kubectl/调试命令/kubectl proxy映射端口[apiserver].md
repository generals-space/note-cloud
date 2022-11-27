# kubectl proxy映射端口

参考文章

1. [kuber官方文档 Accessing Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/)

虽然同是用于端口映射, `proxy`与`port-forward`的区别在于, 后者是用于映射端口到某个pod的, 前者是专用于映射端口到apiserver的.

由于对于访问apiserver需要双向认证, 普通的http客户端工具无法附加ssl证书, 所以才提供了`proxy`这个命令, 通过kubectl进行映射, 实际上就是一个跳过了认证步骤的微服务网关.

```console
$ curl -k https://k8s-server-lb:6443/api/
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {

  },
  "status": "Failure",
  "message": "forbidden: User \"system:anonymous\" cannot get path \"/api/\"",
  "reason": "Forbidden",
  "details": {

  },
  "code": 403
}
```

使用proxy将apiserver的端口映射到本地的8080端口

```
$ k proxy --port=8080
Starting to serve on 127.0.0.1:8080
## 阻塞
```

```console
$ curl http://localhost:8080/api/
{
  "kind": "APIVersions",
  "versions": [
    "v1"
  ],
  "serverAddressByClientCIDRs": [
    {
      "clientCIDR": "0.0.0.0/0",
      "serverAddress": "192.168.0.101:6443"
    }
  ]
}
```

> 注意: 映射出来的端口已经不再是https类型, 所以访问协议应为http.

> 当然, 也有通过curl工具访问apiserver的方法, 参考文章1中也有给出, 不过不在本文的讨论范围, 这里就不细说了.

