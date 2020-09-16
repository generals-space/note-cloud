# kubectl-label标签操作

参考文章

1. [k8s对node添加Label](https://blog.csdn.net/wang725/article/details/89786578)
    - 非常清晰, 一文足够

添加

```console
$ k label node k8s-master-01 aaa=bbb
node/k8s-master-01 labeled
```

查看

```console
$ kgl node
NAME            STATUS   ROLES    AGE   VERSION   LABELS
k8s-master-01   Ready    master   47h   v1.16.2   aaa=bbb,beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-master-01,kubernetes.io/os=linux,node-role.kubernetes.io/master=
k8s-worker-01   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-01,kubernetes.io/os=linux
k8s-worker-02   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-02,kubernetes.io/os=linux
```

更新

```console
$ k label node k8s-master-01 aaa=ccc
error: 'aaa' already has a value (bbb), and --overwrite is false
```

需要加`--overwrite`参数

```console
$ k label node k8s-master-01 aaa=ccc --overwrite
node/k8s-master-01 labeled
$ kgl node
NAME            STATUS   ROLES    AGE   VERSION   LABELS
k8s-master-01   Ready    master   47h   v1.16.2   aaa=ccc,beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-master-01,kubernetes.io/os=linux,node-role.kubernetes.io/master=
k8s-worker-01   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-01,kubernetes.io/os=linux
k8s-worker-02   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-02,kubernetes.io/os=linux
```

删除

```console
$ k label node k8s-master-01 aaa-
node/k8s-master-01 labeled
$ kgl node
NAME            STATUS   ROLES    AGE   VERSION   LABELS
k8s-master-01   Ready    master   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-master-01,kubernetes.io/os=linux,node-role.kubernetes.io/master=
k8s-worker-01   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-01,kubernetes.io/os=linux
k8s-worker-02   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-02,kubernetes.io/os=linux
```

尝试删除一个不存在的 label

```
$ k label node k8s-master-01 aaa-
label "aaa" not found.
node/k8s-master-01 not labeled
```
