# kubectl-字段选择器field selector

参考文章

1. [使用Kubernetes对象6---字段选择器（Field Selectors）](https://blog.csdn.net/u014637098/article/details/88843633)

字段选择器是指可以使用目标资源的 yaml 配置中的某个字段作为过滤条件进行查询, 而不单纯限于 label .

不过字段选择器仍然是有限制的, 不同类型的资源支持的字段选择器不同, 使用不支持的字段选择器会出错. 比如我想筛选`status.hostIP`为`192.168.10.101`的所有Pod.

```
$ k get pod --field-selector status.hostIP=192.168.10.101
Error from server (BadRequest): Unable to find "/v1, Resource=pods" that match label selector "", field selector "status.hostIP=192.168.10.101": field label not supported: status.hostIP
```

这说明不能使用`hostIP`作为过滤条件...

> 所有类型资源都支持`metadata.name`和`metadata.namespace`字段.

