# crd 与 controller

## 1. CustomResourceDefinition

首先, 我们可以通过声明`CustomResourceDefinition`类型的资源创建一种新的资源. 以`PodGroup`为例, 我们创建了一种名为`PodGroup`的资源.

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: podgroups.testgroup.k8s.io
spec:
  group: testgroup.k8s.io
  version: v1
  names:
    kind: PodGroup
    plural: podgroups
  scope: Namespaced
```

部署这个yaml之后我们可以通过`kubectl get PodGroup`获取这个类型的资源(当然, 此时是没有的).

```console
$ k apply -f crd.yaml
customresourcedefinition.apiextensions.k8s.io/podgroups.testgroup.k8s.io created
$ k get PodGroup
No resources found in kube-system namespace.
```

现在可以通过创建类型为`PodGroup`类型的资源了.

```yaml
apiVersion: testgroup.k8s.io/v1
kind: PodGroup
metadata:
  name: podgroup01
spec:
  field: field01
```

其中`apiVersion`为`CustomResourceDefinition`部署文件中声明的, `PodGroup`资源所属的`group`与`version`字段的组合. 并且`spec`下的字段可以随便写, 因为此时还没有创建对应的`Controller`, 而字段的验证工作是要由`Controller`对象完成的.

```console
$ kap pg.yaml
podgroup.testgroup.k8s.io/podgroup01 created
$ k get podgroup
NAME         AGE
podgroup01   9s
```

## 2. PodGroup{} 对象 与 Controller{} 对象

`PodGroup`声明了该类型资源的`spec`下可配置的所有合法字段, `Controller`在开发工作中是一个通用概念, 用于处理一段具体逻辑.

通过运行一个`Controller`, 监听来自 apiserver 的`PodGroup`资源的 CURD 操作, 对新增的`PodGroup`做字段验证, 然后对其监听的Pod做出反应.

但也有很多场合, `Controller`是不需要配合`PodGroup`(这一类资源结构对象)的, ta们直接监听已经存在的`deploy`, `daemonset`等资源, 做一些操作, 比如记录, 回调通知, 作为钩子等, 都可以.

