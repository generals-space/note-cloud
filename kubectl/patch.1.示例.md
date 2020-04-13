# patch.1.示例

1. [Kubernetes kubectl patch 命令详解](http://docs.kubernetes.org.cn/632.html)
2. [使用kubectl patch更新API对象](http://www.coderdocument.com/docs/kubernetes/v1.14/tasks/run_applications/update_api_objects_in_place_using_kubectl_patch.html)
    - `patch`操作的"策略合并"机制

kuber版本: 1.16.2

## 1. patch 试验

用`test-ds.yaml`部署了ds资源, 每个节点上生成一个Pod. 尝试使用`patch`命令修改, 添加一个环境变量`updateTime`, 但是却失败了.

```
$ k patch pod Pod名称 -p '{"spec":{"containers":[{"env":[{"name":"updateTime","value":"123456"}]}]}}'
Error from server: map: map[env:map[updateTime:123456]] does not contain declared merge key: name
```

最初以为是`test-ds`的yaml文件中没有为template添加该环境变量, 于是尝试在yaml文件中添加并更新, 但没有用.

后来我对比了一下`kubectl patch -h`中的示例, 发现好像还要在`containers`数组中添加一个`name`字段(由于`containers`是一个数组, 在对其中一个`container`做出修改时, 需要添加`name`字段定位目标).

```console
$ k patch pod Pod名称 -p '{"spec":{"containers":[{"name":"centos7","env":[{"name":"updateTime","value":"123456"}]}]}}'
The Pod "test-ds-sbj8g" is invalid: spec: Forbidden: pod updates may not change fields other than `spec.containers[*].image`, `spec.initContainers[*].image`, `spec.activeDeadlineSeconds` or `spec.tolerations` (only additions to existing tolerations)
```

还是失败了, 因为Pod资源只能修改`image`, `tolerations`等有限个字段.

```console
$ k patch pod Pod名称 -p '{"spec":{"containers":[{"name":"centos7","image":"registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:devops"}]}}'
pod/test-ds-sbj8g patched
```

成功了.

其他资源的patch字段倒是没有字段限制.

## 2. patch 后果

上面更新了ds派生的其中一个Pod的`image`, 但是Pod并没有重启, IP也未变. 但是对应的container被重建了, 使用`docker ps`可以看到原来的容器已经被删除了(pause容器没有).

Pod未重启这个事情让我觉得挺疑惑的, 于是我又尝试如下两种方法, 更改ds部署文件中container的`image`和`env`, 结果所有Pod都重启了, 两次都是.

```
k patch ds test-ds -p '{"spec": {"template": {"spec": {"containers": [{"name": "centos7", "image": "registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:devops"}]}}}}'

k patch ds test-ds -p '{"spec": {"template": {"spec": {"containers": [{"name": "centos7", "env": [{"name": "newItem", "value": "hehe"}]}]}}}}'
```

## 3. patch 从文件读入

patch子命令没有直接使用`-f`等选项指定读取哪个文件, 但是可以在命令行配合cat使用.

```yaml
{
   "spec": {
        "template": {
            "spec": {
                "containers": [
                {
                    "name": centos7,
                    "image": registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:devops
                }
                ]
            }
        }
   }
}
```

```
k patch ds test-ds --patch "$(cat patch-file.yaml)"
```

> 虽然`patch`不支持从文件读入更改内容, 但是`kubectl`还有一个子命令`replace`可以实现, 看名字就知道, ta的替换类型是全部替换...

