# runtime.Object通用接口使用.2.GVK信息

参考文章

- kube: 1.16.2
- apimachinery: v0.17.0

上一篇写了我没找到相关的接口获取一个`runtime.Object`接口对象的`GVK`信息, 需要在主调函数里手动传入对应类型的`Kind`参数, 但是调用一个通用的函数还要手动传`Kind`信息, 太low了.

但其实还有一种方法, 应该是更适合, 更优雅的方式, 就是`switch..type()`.

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
    var name, kind string
	switch obj.(type) {
	case *appsv1.DaemonSet:
		name = obj.(*appsv1.DaemonSet).Name
		kind = obj.(*appsv1.DaemonSet).Kind
		break
	case *appsv1.StatefulSet:
		name = obj.(*appsv1.StatefulSet).Name
		kind = obj.(*appsv1.StatefulSet).Kind
		break
	}

    klog.Infof("%s", name)  // mysts
    klog.Infof("%s", kind)  // 空
}
```

使用`switch..type()`可以在一定范围内(开发者一定能知道`doit()`方法需要处理哪几种资源对象吧)确认`obj`的类型信息, 这对不同类型资源的处理将会是极大的帮助.

然后是`GVK`信息, 不管是手动创建`StatefulSet`对象, 还是使用`client-go`通过接口获取到的`StatefulSet`对象, 直接获取`Kind`信息必然都会得到空字符串.

但是我们通过`switch..type()`已经能确定obj的具体类型, 而ta的GVK信息是需要自行组装的.

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
	var name string
	var gvk apimSchema.GroupVersionKind
	switch obj.(type) {
	case *appsv1.DaemonSet:
		name = obj.(*appsv1.DaemonSet).Name
		gvk = appsv1.SchemeGroupVersion.WithKind("DaemonSet")
		break
	case *appsv1.StatefulSet:
		name = obj.(*appsv1.StatefulSet).Name
		gvk = appsv1.SchemeGroupVersion.WithKind("StatefulSet")
	case *corev1.Pod:
		name = obj.(*corev1.Pod).Name
		gvk = corev1.SchemeGroupVersion.WithKind("Pod")
		break
	}

    klog.Infof("%s", name)      // mysts
    klog.Infof("%s", gvk)       // apps/v1, Kind=StatefulSet
    klog.Infof("%s", gvk.Kind)  // StatefulSet
}
```

每种资源都是拥有自己所属的`Group`和`Version`的, 如`apps/v1`, 而在一个指定的`GroupVersion`下添加资源, 如`StatefulSet`或是`DaemonSet`, 则是需要在定义资源的`types.go`中进行注册的. [api/apps/v1/register.go](https://github.com/kubernetes/api/blob/v0.17.0/apps/v1/register.go)中, 对`apps/v1`下的各种资源进行了注册.

```go
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Deployment{},
		&DeploymentList{},
		&StatefulSet{},
		&StatefulSetList{},
		&DaemonSet{},
		&DaemonSetList{},
		&ReplicaSet{},
		&ReplicaSetList{},
		&ControllerRevision{},
		&ControllerRevisionList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
```

所以根据`GroupVersion`信息去拼接`GVK`信息是比较正规的做法.
