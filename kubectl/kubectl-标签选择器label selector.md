# kubectl-标签选择器label selector

参考文章

1. [从零开始入门 K8s| K8s 的应用编排与管理](https://zhuanlan.zhihu.com/p/83681561)
    - 查看, 创建 label, 与使用 label 进行过滤的示例.
2. [k8s对node添加Label](https://blog.csdn.net/wang725/article/details/89786578)
    - 非常清晰, 一文足够

- `k get pod -l esname=es-hjl-13`
- `k get pod -l esname=es-hjl-13,estype=data`
- `k get pod -l esname=es-hjl-13,estype!=data`
- `k get pod -l 'esname=es-hjl-13,estype in (data,master)'`

支持的运算符有`=`, `!=`, `==`, `in`, `notin`. `k get --help`信息中没有介绍, 我是在错误信息里找到的...

```console
$ k get pod -l 'component not in (kube-apiserver)'
Error from server (BadRequest): Unable to find "/v1, Resource=pods" that match label selector "component not in (kube-apiserver)", field selector "": unable to parse requirement: found 'not', expected: '=', '!=', '==', 'in', notin'
```

还有一种是按照是否存在某个标签来过滤的.

```
k get pod -l 'label01'
k get pod -l '!label01'
```

## 示例

### 添加

```console
$ k label node k8s-master-01 aaa=bbb
node/k8s-master-01 labeled
```

### 查看

```console
$ kgl node
NAME            STATUS   ROLES    AGE   VERSION   LABELS
k8s-master-01   Ready    master   47h   v1.16.2   aaa=bbb,beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-master-01,kubernetes.io/os=linux,node-role.kubernetes.io/master=
k8s-worker-01   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-01,kubernetes.io/os=linux
k8s-worker-02   Ready    <none>   47h   v1.16.2   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,kubernetes.io/arch=amd64,kubernetes.io/hostname=k8s-worker-02,kubernetes.io/os=linux
```

### 更新

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

### 删除

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
