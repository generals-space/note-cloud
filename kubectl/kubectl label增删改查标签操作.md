# kubectl label增删改查标签操作

参考文章

1. [从零开始入门 K8s| K8s 的应用编排与管理](https://zhuanlan.zhihu.com/p/83681561)
    - 查看, 创建 label, 与使用 label 进行过滤的示例.
2. [k8s对node添加Label](https://blog.csdn.net/wang725/article/details/89786578)
    - 非常清晰, 一文足够

## 添加

```console
$ k label node k8s-master-01 aaa=bbb
node/k8s-master-01 labeled
```

## 查看

```console
$ k get node --show-labels
NAME            STATUS   ROLES    AGE   VERSION   LABELS
k8s-master-01   Ready    master   47h   v1.16.2   aaa=bbb,kubernetes.io/hostname=k8s-master-01,node-role.kubernetes.io/master=
```

## 更新

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
k8s-master-01   Ready    master   47h   v1.16.2   aaa=ccc,kubernetes.io/hostname=k8s-master-01,node-role.kubernetes.io/master=
```

## 删除

就用一个减号`-`

```console
$ k label node k8s-master-01 aaa-
node/k8s-master-01 labeled
$ kgl node
NAME            STATUS   ROLES    AGE   VERSION   LABELS
k8s-master-01   Ready    master   47h   v1.16.2   kubernetes.io/hostname=k8s-master-01,node-role.kubernetes.io/master=
```

尝试删除一个不存在的 label 会报错.

```
$ k label node k8s-master-01 aaa-
label "aaa" not found.
node/k8s-master-01 not labeled
```
