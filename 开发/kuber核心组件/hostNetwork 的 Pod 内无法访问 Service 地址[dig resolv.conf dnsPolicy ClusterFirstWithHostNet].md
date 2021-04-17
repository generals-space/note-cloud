# hostNetwork 的 Pod 内无法访问 Service 地址[dig resolv.conf]

参考文章

1. [k8s hostNetwork涉及到dns问题](https://linuxeye.com/470.html)
2. [dnsPolicy in hostNetwork not working as expected](https://github.com/kubernetes/kubernetes/issues/87852)
3. [kuber 官方文档 Pod's DNS Policy](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy)

环境

- kuber v1.17.3, 单 master 节点
- 宿主机(master节点): 172.16.91.10
- 宿主机 dns : 172.16.91.2 (vmware虚拟机的网关)

## 问题描述

在写扩展调度器时发现, 通过`hostNetwork`部署的`kube-scheduler`无法通过 Service 名称访问到`kueb-scheduler-extender`. 做了一些实验, 发现这是正常现象.

`hostNetwork`模式直接使用宿主机网络, 而就算是在 kuber 宿主机上, 也没办法 ping 通 `Service`(不过可以 ping 通 `ServiceIP`).

如下是直接使用`hostNetwork`部署的 Pod 中, `/etc/resolv.conf`文件的内容, 与该 Pod 所在的宿主机完全相同.

```
nameserver 172.16.91.2
```

> `172.16.91.2`是宿主机中配置的 dns 服务器地址

按照参考文章1所说, 在使用`hostNetwork`模式时, 同时指定`dnsPolicy`字段为`ClusterFirstWithHostNet`, 就可以了. 如下

```yaml
  hostNetwork: true
  dnsPolicy: ClusterFirstWithHostNet
```

这样 Pod 内部的`/etc/resolv.conf`的内容就变成了如下

```
nameserver 10.96.0.10
search kube-system.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
```

> `10.96.0.10`为 coredns 组件的地址

这样在 Pod 内就可以访问 Service 了.

## 原理解析

kuber 内部的 Service 是通过 coredns 组件解析成 ServiceIP 的, 宿主机环境用的 dns 服务没有办法解析. 

在宿主机环境执行如下命令

```console
$ dig kube-dns

;; QUESTION SECTION:
;kube-dns.			IN	A
```

没有结果, 但我们可以指定以 coredns 服务的地址进行解析.

```console
$ kwd pod | grep coredns
coredns-67c766df46-2xsqh                  1/1     Running   1          5d22h   10.254.0.11    k8s-master-01   <none>           <none>
coredns-67c766df46-xr7wb                  1/1     Running   1          5d22h   10.254.0.10    k8s-master-01   <none>           <none>
$ dig @10.254.0.10 kube-dns

;; QUESTION SECTION:
;kube-dns.			IN	A

$ dig @10.254.0.10 kube-dns.kube-system.svc

;; QUESTION SECTION:
;kube-dns.kube-system.svc.	IN	A

$ dig @10.254.0.10 kube-dns.kube-system.svc.cluster.local

;; QUESTION SECTION:
;kube-dns.kube-system.svc.cluster.local.	IN A

;; ANSWER SECTION:
kube-dns.kube-system.svc.cluster.local.	30 IN A	10.96.0.10
```

注意, 虽然指定了 coredns 的地址, 但是在 kuber 环境外(宿主机环境)进行解析, 需要为 Service 名称添加上全部的域名前缀, 默认为`cluster.local`.

------

如果你了解`resolv.conf`的文件配置的话, 就会理解, 为什么在 Pod 内部访问(相同`namespace`的) Service 只需要指定名称就行, 而不是像上面需要指定全域名.

该文件中的`search`字段, 就是为待查询的域名自动加上后缀再进行查询的. 比如在容器内`ping kube-dns`, dns 服务会先尝试将域名补全为`kube-dns.kube-system.svc.cluster.local`再去查询.

但是, 不同`namespace`下的 Pod 中, `resolv.conf`的`search`字段是不同的, 只会带有自己命令空间的名称, 所以只有在访问相同命名空间的 Service 可以写短名称, 如果想访问`ns01`下的`service01`, 那么补全成为`service01.kube-system.svc.cluster.local`也是找不到的.
