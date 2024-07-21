# kubectl proxy映射端口

参考文章

1. [Accessing Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/)
    - 官方文档

虽然同是用于端口映射, `proxy`与`port-forward`的区别在于, 后者是用于映射端口到某个pod的, 前者是专用于映射端口到apiserver的.

由于对于访问apiserver需要双向认证, 普通的http客户端工具无法附加ssl证书, 所以才提供了`proxy`这个命令, 通过kubectl进行映射, 实际上就是一个跳过了认证步骤的微服务网关.

```json
// curl -k https://kube-apiserver.generals.space:6443/api/
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

```log
$ k proxy --port=8080
Starting to serve on 127.0.0.1:8080
## 阻塞
```

```json
// curl http://localhost:8080/api/
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

------

上述代理只能在本机上访问, 如果要将该端口开放到公网(仅供测试), 可以使用如下命令.

```
kubectl proxy --port=6444 --address='0.0.0.0' --api-prefix=/ --accept-hosts='^*$' &
```

kubectl proxy命令是在执行者所在终端搭建了一个 http server, ta的后端就是 kubectl 通过 /root/.kube/config 与 apiserver 鉴权后建立的连接, 因此 proxy 所暴露出的权限最大就是该 kubeconfig 中配置的权限.

## 扩展

常规 operator , 不管是运行在 Pod 里还是 Pod 外, 都使用 ssl 证书与 apiserver 进行双向认证.

有时我们希望调试 operator 与 apiserver 之间的通信数据, 但是通过 tcpdump 抓包, 再使用 wireshark 进行分析时, 界面上只能显示协议类型为`TLS`, 无法展示为 http, 也就无法查看 uri, query 等信息.

此时就可以通过 kubectl proxy 将 apiserver 代理成 http 服务, 然后在 operator 使用的 kubeconfig 中, 将`clusters[].cluster.server`的地址, 由`https`改成`http`, 同时修改端口为 proxy 的端口. 

然后再启动 operator, 此时双方就使用 http 协议进行通信了.
