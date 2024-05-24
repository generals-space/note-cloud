# patch.1.示例

1. [Kubernetes kubectl patch 命令详解](http://docs.kubernetes.org.cn/632.html)
2. [使用kubectl patch更新API对象](http://www.coderdocument.com/docs/kubernetes/v1.14/tasks/run_applications/update_api_objects_in_place_using_kubectl_patch.html)
    - `patch`操作的"策略合并"机制

kuber版本: 1.16.2

在对`metadata.label`字段的`patch`操作时, 没有体现出`strategic`与`merge`的区别, 这里以`containers`部分为例.

在我目前的理解中, 这两个在对map类型的类型进行`patch`时没有区别, 都是直接合并替换. 只有在对包含数组成员的部分时, 才有区别.

以如下ds资源为例

> 最好不要用pod实验, 因为pod资源只能修改`image`, `tolerations`等有限个字段, 可能会对`patch`操作的理解出现偏差.

```yaml
spec:
  template:
    spec:
      containers:
      - name: centos7
        image: centos7:v1
```

其中`containers`部分为数组类型, 如果以如下语句进行`patch`

```yaml
spec:
  template:
    spec:
      containers:
      - name: centos8
        image: centos8:v1
```

结果会是什么?

`strategic`策略合并的结果会是

```yaml
spec:
  template:
    spec:
      containers:
      - name: centos7
        image: centos7:v1
      - name: centos8
        image: centos8:v1
```

`merge`的结果会是

```yaml
spec:
  template:
    spec:
      containers:
      - name: centos8
        image: centos8:v1
```

以下给出对test-ds资源的`patch`操作示例.

## 2. `strategic`实践

### 合并 env

下面的示例是为`centos7`这个容器指定`env`参数(可与原`env`进行合并, 而非单纯替换)

```log
$ k patch ds test-ds -p '{"spec":{"template":{"spec":{"containers":[{"env":[{"name":"updateTime","value":"123456"}]}]}}}}'
Error from server: map: map[env:[map[name:updateTime value:123456]]] does not contain declared merge key: name
```

这是因为策略合并在合并`containers`部分时, 需要指定`contaners[n].name`字段, 以确定合合并的是`containers`数组中的哪个容器, 如下

```log
$ k patch ds test-ds -p '{"spec":{"template":{"spec":{"containers":[{"name":"centos7","env":[{"name":"updateTime","value":"123456"}]}]}}}}'
daemonset.apps/test-ds patched
```

其他类型的资源, 或是除了`containers`部分的策略合并时, 应该也有类似的"锚点"机制, 不过我还没有找到明确的文档, 这里就不深入讨论了.

### 合并 containers

如果指定的`contaners[n].name`值不等于已有的字段值, 那么策略合并机制会认为我们想要在`containers`中再追加一个容器.

```
$ k patch ds test-ds -p '{"spec":{"template":{"spec":{"containers":[{"name":"centos8","image":"registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7"}]}}}}'
daemonset.apps/test-ds patched
```

## 3. `merge`实践

与`strategic`相比, `merge`的操作就非常单纯了, 对于数组中的成员, 全部都是直接替换.

```log
$ k patch ds test-ds --type merge -p '{"spec":{"template":{"spec":{"containers":[{"env":[{"name":"updateTime","value":"123456"}]}]}}}}'
The DaemonSet "test-ds" is invalid:
* spec.template.spec.containers[0].name: Required value
* spec.template.spec.containers[0].image: Required value
```

因为是直接替换, 所以`containers`成员中必须存在`name`, `image`.

```
k patch ds test-ds --type merge -p '{"spec":{"template":{"spec":{"containers":[{"name": "centos7", "image":"registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7", "env":[{"name":"updateTime","value":"123456"}]}]}}}}'
```

注意, 这一句执行了, 就不能达到合并`env`的目的了, 整个`containers`都被替换掉了...
