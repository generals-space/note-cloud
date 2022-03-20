# runtime.Object通用接口使用

参考文章

1. [crd-ipkeeper/pkg/staticip/new.go](https://github.com/generals-space/crd-ipkeeper/blob/ccdf7e693a4edf309db551e90ab94e5411caf270/pkg/staticip/new.go#L50)

- kuber: 1.16.2
- apimachinery: v0.17.0

`runtime`是指`k8s.io/apimachinery/pkg/runtime`包(注意, 不是`k8s.io/apimachinery/pkg/util/runtime`), ta的`Object`接口声明如下.

```go
type Object interface {
    // ObjectKind 表示资源的 GVK 信息
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() Object
}
```

ta可以表示所有 kube 资源(如`Deployment`, `Daemonset`, `Statefulset`等), 也包含自定义资源(不过这些资源的指针对象才可以).

```go
import(
	apimMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
)
func (c *Controller) caller(){
    sts := &appsv1.StatefulSet{
		ObjectMeta: apimMetav1.ObjectMeta{
			Name: "mysts",
		},
	}
	c.doit(sts)
}
func (c *Controller) doit(obj runtime.Object) {
    klog.Infof("%s", obj.GetObjectKind())   // &TypeMeta{Kind:,APIVersion:,}
}
```

不过这个接口只有上面2个方法, 而且基本没用, 当我们想要获取某个资源对象时(同样不确定ta的类型), 需要将传入的接口对象转换成其他接口类型.

```go
import(
	apimMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
)
func (c *Controller) caller(){
	sts := &appsv1.StatefulSet{
		TypeMeta: apimMetav1.TypeMeta{
            // 这里的 xxx 真的有影响.
            // 不过, 使用 client-go 工具创建 StatefulSet 资源时, 貌似可以不用指定 TypeMeta{} 部分
            Kind: "xxx",    
		},
		ObjectMeta: apimMetav1.ObjectMeta{
			Name: "mysts",
		},
	}
	c.doit(sts)
}
func (c *Controller) doit(obj runtime.Object) {
    klog.Infof("%s", obj.GetObjectKind())                                               // &TypeMeta{Kind:xxx,APIVersion:,}
    klog.Infof("%s", obj.(*appsv1.StatefulSet).GetObjectKind())                         // &TypeMeta{Kind:xxx,APIVersion:,}

    klog.Infof("%s", obj.(apimMetav1.ObjectMetaAccessor).GetObjectMeta().GetName())     // mysts
    klog.Infof("%s", obj.(apimMetav1.Object).GetName())                                 // mysts
}
```

上面的示例只展示了可以通过接口转换获取通用类型的 meta 部分的信息, `GVK`貌似不属于 meta, 我没有找到对应的接口可以获取`GVK`信息, 也没有找到可以获取`Spec`和`Status`部分的信息, 可能后两者根本没有通用的接口类型...

我现在想到的是, 所有 kuber 资源(包括CRD资源), 都拥有`metav1.TypeMeta`和`metav1.ObjectMeta`, 如果一定要获取某个不确定类型资源的`GVK`信息, 可以在主调函数里将这两个成员传入.

对于`Spec`和`Status`, 如果之后仍然无法找到ta们的通用接口的话, 可以尝试将ta们转换成`map[string]interface{}`, 判断其各自的字段.

参考文章1, 是我之前写的小项目, 虽然不是很典型, 但是思路与我上面的不谋而合.

关于如何得到一个`runtime.Object`对象的`GVK`信息, 在我下一篇文章中会讲到.
