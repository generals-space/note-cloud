# kubectl get --selector标签选择器[label]

参考文章

1. [从零开始入门 K8s| K8s 的应用编排与管理](https://zhuanlan.zhihu.com/p/83681561)
    - 查看, 创建 label, 与使用 label 进行过滤的示例.
2. [k8s对node添加Label](https://blog.csdn.net/wang725/article/details/89786578)
    - 非常清晰, 一文足够
3. [kubernetes：字段选择器（field-selector）标签选择器（labels-selector）和筛选 Kubernetes 资源](https://blog.csdn.net/fly910905/article/details/102572878/)
    - 字段选择器（field-selector）
    - 标签选择器（labels-selector）

- `k get pod -l esname=es-13`
- `k get pod -l esname=es-13,estype=data`
- `k get pod -l esname=es-13,estype!=data`
- `k get pod -l 'esname=es-13,estype in (data,master)'`

支持的运算符有`=`, `!=`, `==`, `in`, `notin`. `k get --help`信息中没有介绍, 我是在错误信息里找到的...

```log
$ k get pod -l 'component not in (kube-apiserver)'
Error from server (BadRequest): Unable to find "/v1, Resource=pods" that match label selector "component not in (kube-apiserver)", field selector "": unable to parse requirement: found 'not', expected: '=', '!=', '==', 'in', notin'
```

还有一种是按照是否存在某个标签来过滤的.

```
k get pod -l 'label01'
k get pod -l '!label01'
```
