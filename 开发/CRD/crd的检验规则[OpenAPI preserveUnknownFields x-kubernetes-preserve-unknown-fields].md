# crd的检验规则[OpenAPI v1]

参考文章

1. [kubernetes 自定义资源（CRD）的校验](https://cloud.tencent.com/developer/article/1557507)

创建一个名为LogstashCluster的crd对象, 由于还在开发阶段, 有可能需要在`spec`块中增删字段, 所以我们希望尽量不在 crd 资源中把字段写死. 否则, 在创建 cr 实例时, kube 集群会把未在 crd 对象中声明的字段移除掉.

在 v1.13.2 集群中, 可以使用`preserveUnknownFields`, 同时只定义spec/status类型为`object`, 不声明该块下的`properties`信息.

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: logstashclusters.crd.generals.space
spec:
  group: crd.generals.space
  names:
    kind: LogstashCluster
    listKind: LogstashClusterList
    plural: logstashclusters
    shortNames:
    - logstash
    singular: logstashcluster
  preserveUnknownFields: true
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: LogstashCluster is the Schema for the logstashclusters API
        properties:
          apiVersion:
            description: ""
            type: string
          kind:
            description: ""
            type: string
          metadata:
            type: object
          spec:
            description: ""
            type: object
          status:
            description: ""
            type: object
        type: object
    served: true
    storage: true
```

在 v1.23.0 集群中, 可以使用`x-kubernetes-preserve-unknown-fields`, 同时只定义spec/status类型为`object`, 不声明该块下的`properties`信息.

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: logstashclusters.crd.generals.space
spec:
  group: crd.generals.space
  names:
    kind: LogstashCluster
    listKind: LogstashClusterList
    plural: logstashclusters
    shortNames:
    - logstash
    singular: logstashcluster
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: LogstashCluster is the Schema for the logstashclusters API
        properties:
          apiVersion:
            description: ""
            type: string
          kind:
            description: ""
            type: string
          metadata:
            type: object
          spec:
            description: ""
            type: object
            x-kubernetes-preserve-unknown-fields: true
          status:
            description: ""
            type: object
            x-kubernetes-preserve-unknown-fields: true
        type: object
    served: true
    storage: true
```
