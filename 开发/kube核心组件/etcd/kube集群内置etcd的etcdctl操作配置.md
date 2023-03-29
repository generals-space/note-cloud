# kuber集群内置etcd的etcdctl操作配置

参考文章

1. [etcdctl环境变量设置](https://www.cnblogs.com/lizhaoxian/p/11498268.html)
2. [官方仓库 etcdctl readme](https://github.com/etcd-io/etcd/tree/master/etcdctl)

**环境**

- version: 3.3
- API: 3
- kube: 1.16.2(1主2从)
    - 1主 `192.168.0.101`
    - 2从 `192.168.1.104/105`

etcd服务的Pod在启动时将宿主机的`/etc/kubernetes/pki/etcd`目录挂载到pod内部.

进入到pod内部, 执行如下命令即可.

```
ETCDCTL_API=3 etcdctl --endpoints=https://127.0.0.1:2379 --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key get / --prefix
```

写这么一大长串很麻烦, 可以使用环境变量, 用以简化.

```bash
export ETCDCTL_API=3
export ETCDCTL_ENDPOINTS=https://127.0.0.1:2379
export ETCDCTL_CACERT=/etc/kubernetes/pki/etcd/ca.crt
export ETCDCTL_CERT=/etc/kubernetes/pki/etcd/server.crt
export ETCDCTL_KEY=/etc/kubernetes/pki/etcd/server.key
etcdctl get / --prefix
```

另外, 由于kubeadm创建的etcd服务使用的是`hostNetwork`, 可以直接在集群外访问, 拷贝所需的证书和密钥即可.

```
ETCDCTL_API=3 etcdctl --endpoints=https://192.168.0.101:2379 --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key get / --prefix
```
