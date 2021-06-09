# patch.1.type策略[metadata.label json merge strategic]

参考文章

1. [Kubernetes kubectl patch 命令详解](http://docs.kubernetes.org.cn/632.html)
2. [使用kubectl patch更新API对象](http://www.coderdocument.com/docs/kubernetes/v1.14/tasks/run_applications/update_api_objects_in_place_using_kubectl_patch.html)
    - `patch`操作的"策略合并"机制

kuber版本: 1.16.2

`patch`子命令有`--type`选项: 可选值: `json`, `merge`和`strategic`. 默认为`strategic`, 称为"策略合并".

## 1. strategic

按照参考文章2中所说的示例, 更新`ds`资源中的`image`字段, 新生成的Pod资源对象中会出现两个`image`, 即初始部署文件与patch文件进行了"合并". 但是对`tolerations`执行patch操作, 原来的`tolerations`字段会直接被替换掉.

> 但我在实验中单独对Pod与直接对DS资源进行`patch`(都是更改`image`字段), 并没有出现合并(不管是Pod还是DS的`-o`输出中都没有合并).

这种选择性的替换/合并, 被称为"策略合并".

所有类型的资源定义中(`XXXSpec`部分)都声明了该类型哪些字段可以在更新时合并, 以`Pod`为例

```go
type PodSpec struct {
    ...
    Containers []Container `json:"containers" patchStrategy:"merge" patchMergeKey:"name" ...`
```

在对`label`标签进行修改时, 可以使用如下命令.

```
kubectl patch node k8s-master-01 --type strategic -p '{"metadata": {"labels": {"key01": "val01"}}}'
```

## 2. merge

自动合并, 但是这种操作貌似没有办法进行删除操作, 顶多把目标key置空. 同样是以对`label`标签的修改为例

```
kubectl patch node k8s-master-01 --type merge -p '{"metadata": {"labels": {"key01": "val01"}}}'
```

在本文中是以`metadata.labels`字段为例, 无法显示出`merge`和`strategic`在操作行为上的区别, 可以见另一篇文章.

## 3. json (拥有删除能力)

与`merge`和`strategic`相比, `json`方式接收的参数结构会有所不同, 如下

1. 可以指定操作符, 包括`add`, `replace`, `remove`
2. 目标字段路径为目录路径而非json对象, 如`/metadata/labels/key01`

新增

```
kubectl patch node k8s-master-01 --type json -p '[{"op": "add", "path": "/metadata/labels/key01", "value": "val01"}]'
```

更新

```
kubectl patch node k8s-master-01 --type json -p '[{"op": "replace", "path": "/metadata/labels/key01", "value": "val02"}]'
```

删除

```
kubectl patch node k8s-master-01 --type json -p '[{"op": "remove", "path": "/metadata/labels/key01"}]'
```

与`merge`和`strategic`相比, `json`方式拥有了删除字段的能力, 就是使用起来可能不太灵活, `path`路径需要预先确定好, 在程序中只能实现特定功能, 不好通用.
