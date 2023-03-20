# ClusterRole.rules.apiGroups.resourceNames限制单个资源对象

参考文章

1. [Proper use of Role.rules.resourceNames for creating pods with limited access to resources](https://stackoverflow.com/questions/65202615/proper-use-of-role-rules-resourcenames-for-creating-pods-with-limited-access-to)
2. [Using RBAC Authorization - Referring to resources](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-resources)
    - 官方文档

某次在RBAC的 Role/ClusterRole 定义中, 发现有一个`resourceNames`, 不知道是啥涵义.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: configmap-updater
rules:
- apiGroups: [""]
  #
  # at the HTTP level, the name of the resource for accessing ConfigMap
  # objects is "configmaps"
  resources: ["configmaps"]
  resourceNames: ["my-configmap"]
  verbs: ["update", "get"]
```

按照参考文章1, 2所说, `resourceNames`可以将资源限制到其指定的对象上. 如上面的角色配置, 就只允许其绑定的用户有能力 get/update 名为 "my-configmap"的 ConfigMap 对象, 其他的 ConfigMap 则无权限操作.

`resourceNames`只能针对已经存在的对象, 就算在`verbs`中配置了`create`操作, 也没办法创建.

如果是`list`, `watch`, 则在请求的时候, 需要添加`--field-selector=metadata.name`进行筛选, 如`kubectl get configmaps --field-selector=metadata.name=my-configmap`.
