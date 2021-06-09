
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
$ k patch pod Pod名称 -p '{"spec":{"containers":[{"name":"centos7","image":"registry.cn-hangzhou.aliyuncs.com/generals-space/centos7-devops"}]}}'
pod/test-ds-sbj8g patched
```

成功了.