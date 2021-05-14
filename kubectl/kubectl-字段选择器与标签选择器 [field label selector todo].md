# kuber选择器-字段选择器与标签选择器 [field label selector]

参考文章

1. [kubernetes：字段选择器（field-selector）标签选择器（labels-selector）和筛选 Kubernetes 资源](https://blog.csdn.net/fly910905/article/details/102572878/)
    - 字段选择器（field-selector）
    - 标签选择器（labels-selector）

按照Pod标签与所在主机, 查询Pod

```
k get pod --selector middleware=es --field-selector spec.nodeName=
```
