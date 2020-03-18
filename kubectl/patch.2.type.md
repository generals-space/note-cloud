参考文章

1. [使用kubectl patch更新API对象](http://www.coderdocument.com/docs/kubernetes/v1.14/tasks/run_applications/update_api_objects_in_place_using_kubectl_patch.html)
    - `patch`操作的"策略合并"机制

> `patch`子命令有`--type`选项: 可选值: `json`, `merge`和`strategic`, 默认为`strategic`, 称为"策略合并".

按照参考文章1中所说的示例, 更新部署文件中的`image`字段, 新生成的Pod资源对象中会出现两个`image`, 即初始部署文件与patch文件进行了"合并". 但是对`tolerations`执行patch操作, 原来的`tolerations`字段会直接被替换掉.

这种选择性的替换/合并, 被称为"策略合并".

但我在实验中单独对Pod与直接对DS资源进行patch(都是更改`image`字段), 并没有出现合并(不管是Pod还是DS的`-o`输出中都没有合并).

另外两个patch选项, `merge`是直接合并, 而`json`也并不用于指定一种格式(因为patch内容本来就应该是json字符串), 而是指替换操作, 直接替换成patch的内容.

所有类型的资源定义中(`XXXSpec`部分)都声明了该类型哪些字段可以在更新时合并, 以Pod为例

```go
type PodSpec struct {
    ...
    Containers []Container `json:"containers" patchStrategy:"merge" patchMergeKey:"name" ...`
```

------

试了下`merge`完全没效果.

```
$ k patch ds test-ds --type merge -p '{"spec": {"template": {"spec": {"containers": [{"name": "centos7", "image":"generals/centos7"}]}}}}'
daemonset.apps/test-ds patched
```

以为ds生成的Pod中会出现两个container呢.