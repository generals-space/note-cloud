# 根据field进行过滤[fieldSelector]

参考文章

1. [Field Selector spec.nodeName does not work](https://github.com/kubernetes-sigs/kubebuilder/issues/547)
2. [client.MatchingFields should be a fields.Selector](https://github.com/kubernetes-sigs/controller-runtime/issues/576)
3. [List custom resources from caching client with custom fieldSelector](https://stackoverflow.com/questions/57083221/list-custom-resources-from-caching-client-with-custom-fieldselector)
    - 采纳答案值得阅读

希望查询某一主机上的所有pod, 在使用`List()`方法与`ListOptions{}`参数时, 执行却出错了.

```go
import(
    "sigs.k8s.io/controller-runtime/pkg/client"
    "k8s.io/apimachinery/pkg/fields"
)
podList := &corev1.PodList{}
err := r.Client.List(ctx, podList, &client.ListOptions{
    FieldSelector: fields.SelectorFromSet(fields.Set{
        "spec.nodeName": "xxx",
    })
})
```

在运行到`List()`的时候, 出现了如下错误

```log
Index with name field:spec.nodeName does not exist
```

在使用`code-generator`生成的工程中, 上述操作是没有问题的. 

但是在`kubebuilder`创建的工程中, `r.Client`被分为了2部分, 其中`Get/List`为读操作, `Create/Update/Delete`为写操作, 分别由`r.Client`中的两个成员`Reader`和`Writer`完成.

`r.Client`在查询本地内容时, `Reader`借助了本地缓存(`Indexer`), 并不直接请求kube集群. 如果使用 fieldSelector, 那么需要事先创建该 field 的索引.

```go
import(
    "k8s.io/apimachinery/pkg/runtime"
)
cache := mgr.GetCache()
indexFunc := func(obj runtime.Object) []string {
    return []string{obj.(*corev1.Pod).Spec.NodeName}
}
if err := cache.IndexField(&corev1.Pod{}, "spec.nodeName", indexFunc); err != nil {
    panic(err)
}
```
