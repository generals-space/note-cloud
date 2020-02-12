# kuber集群存储在etcd的数据解析

kuber在etcd中的全部数据都存储在`/registry`目录下.

可以使用如下命令查看各目录下的一级内容而不是递归查看.

```bash
dir=""
dir=/registry
dir=/registry/services
etcdctl get --prefix --keys-only $dir | sed -n "s#$dir\/\([^\/]*\).*#\1#p" | uniq
```

> `sed`中使用`#`为分隔符, 是因为`sed`的pattern中存在变量`$dir`, 且`$dir`中包含`/`字符.

```console
$ dir=/registry/services
$ etcdctl get --prefix --keys-only $dir | sed -n "s#$dir\/\([^\/]*\).*#\1#p" | uniq
endpoints
specs
```

------

查看`/registry`会发现, 该目录下存储的键都是各种类型的资源名称.

```
$ etcdctl get --prefix --keys-only $dir | sed -n "s#$dir\/\([^\/]*\).*#\1#p" | uniq
apiextensions.k8s.io
apiregistration.k8s.io
clusterrolebindings
clusterroles
configmaps
controllerrevisions
crd.projectcalico.org
daemonsets
deployments
etcd.database.coreos.com
leases
masterleases
minions
namespaces
persistentvolumeclaims
persistentvolumes
pods
priorityclasses
ranges
replicasets
rolebindings
roles
samplecontroller.k8s.io
secrets
serviceaccounts
services
storageclasses
```

这些都是集群中已经创建的资源(比如就没有job资源), 但是好像又没那么简单, 上面的列表中没有`endpoints`类型的资源. 是因为ta们在`services`目录下.

```
$ dir=/registry/services
$ etcdctl get --prefix --keys-only $dir | sed -n "s#$dir\/\([^\/]*\).*#\1#p" | uniq
endpoints
specs
```

不知道是出于什么目录, 也许不同类型的资源存储的结构不同吧.

大致上, 不同资源下面就都是按照ns来区分了.

```console
$ dir=/registry/pods
$ etcdctl get --prefix --keys-only $dir | sed -n "s#$dir\/\([^\/]*\).*#\1#p" | uniq
default
etcd
kube-system
local-path-storage
```

即, 上述4个ns中存在pod资源.
