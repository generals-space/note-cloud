# kuber选择器-标签选择器label selector

参考文章

1. [从零开始入门 K8s| K8s 的应用编排与管理](https://zhuanlan.zhihu.com/p/83681561)
    - 查看, 创建 label, 与使用 label 进行过滤的示例.

- `k get pod -l esname=es-hjl-13`
- `k get pod -l esname=es-hjl-13,estype=data`
- `k get pod -l esname=es-hjl-13,estype!=data`
- `k get pod -l 'esname=es-hjl-13,estype in (data,master)'`

支持的运算符有`=`, `!=`, `==`, `in`, `notin`. `k get --help`信息中没有介绍, 我是在错误信息里找到的...

```console
$ k get pod -l 'component not in (kube-apiserver)'
Error from server (BadRequest): Unable to find "/v1, Resource=pods" that match label selector "component not in (kube-apiserver)", field selector "": unable to parse requirement: found 'not', expected: '=', '!=', '==', 'in', notin'
```
