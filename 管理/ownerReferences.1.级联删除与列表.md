参考文章

1. [为什么kubernetes的ownerReference要设计成列表，存在多个ownerReference的情况吗？](https://joshua.su/8a98cc78605e)
2. [kubernetes/design-proposals-archive](https://github.com/kubernetes/design-proposals-archive/blob/acc25e14ca83dfda4f66d8cb1f1b491f26e78ffe/api-machinery/controller-ref.md)
    - 官方文档

常见的`ownerReferences`配置.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-665bb474c-kkrgb
  namespace: system
  ownerReferences:
  - apiVersion: apps/v1
    kind: ReplicaSet
    name: test-665bb474c
    uid: f2a4a16f-aa5a-4e30-825f-34c9b808159d
    blockOwnerDeletion: true
    controller: true
  resourceVersion: "3767755"
  uid: 8d59503c-1b7f-4ec8-b396-2a20108f0b17
```

## Garbage Collection/垃圾清理/级联删除

该能力实现主要依赖如下4个参数

- apiVersion
- kind
- name
- uid

`blockOwnerDeletion`与`controller`倒不是那么重要.

关于级联删除, 有一些细节需要注意

| Parent     | Child      | 是否可以级联删除 |
| :--------- | :--------- | :--------------- |
| Namespaced | Namespaced | yes              |
| Cluster    | Cluster    | yes              |
| Namespaced | Cluster    | no               |
| Cluster    | Namespaced | yes              |

如果 Parent 是 Namespaced, 而 Child 是 Cluster scope, 则 Parent 资源被删除时, Child 不会有任何变化, 估计是无法找到对应的资源.

Parent在被删除时, 由 k8s (应该是 controller-manager 组件), 向 Child 发起 delete 操作(会添加`deletionTimestamp`字段). 

即使 Child 是 CRD 类型资源, 也不需要 controller 就可以被删除. 不过如果 Child 拥有 finalizer, 还是需要处理一下的.

###

如下是一个简单的 CRD 测试资源

```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: testcrds.kube.generals.space
spec:
  group: kube.generals.space
  names:
    kind: TestCrd
    listKind: TestCrdList
    plural: testcrds
    singular: testcrd
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        ## properties:
        ##   apiVersion:
        ##     type: string
        ##   kind:
        ##     type: string
        ##   metadata:
        ##     type: object
        type: object
    served: true
    storage: true
```

对应的 cr

```yaml
---
apiVersion: kube.generals.space/v1
kind: TestCrd
metadata:
  name: mytestcrd
```

可与普通的 ConfigMap 资源做联合测试.

## 多个ownerReferences/ownerReferences列表

`ownerReferences`是一个列表, 虽然在 k8s 中, 某个子资源同时属于多个父资源并常见, 但也不是不可能.

至少原生资源中没有过, 所以一般都是开发者按照需求自行设置的.

比如一个数据库配置 ConfigMap, 被多个数据库实例 Pod 引用, 那ta就可以属于多个 Pod.

但总不能删除其中一个 Pod, ConfigMap 子资源就跟着被删除了吧.

所以, 一个从属于多个父资源的子资源, 需要在所有父资源被删除后, 才会被级联删除.

```yaml
apiVersion: kube.generals.space/v1
kind: TestCrd
metadata:
  name: mytestcrd
  finalizers:
  - finalizer.kube.generals.space
  ownerReferences:
  - apiVersion: kube.generals.space/v1
    kind: TestCrd
    name: mytestcrd2
    uid: 655022e6-5b97-4852-8d19-50275ab237bf
  - apiVersion: kube.generals.space/v1
    kind: TestCrd
    name: mytestcrd3
    uid: 1cf5cc4d-41db-4ef0-aaa4-0654c30fb6c4
  resourceVersion: "4070170"
  uid: 3622e12a-5e7d-448a-899d-d836e5aa8e17
```

> `finalizers`是想看看最后一个 Parent 被删除后会发生什么.

删除其中一个父资源时, 子资源中对应的`ownerReference`成员会随之删除.

但是删除最后一个父资源时, `ownerReference`并不会被删除, 而是直接发起删除指令, 新增`deletionTimestamp`字段, 直到`finalizers`被处理.

