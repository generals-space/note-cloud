# kuber-CRD.1.概念理解

参考文章

1. [kubernetes系列之十四：Kubernetes CRD(CustomResourceDefinition)概览](https://blog.csdn.net/cloudvtech/article/details/80277960)
    - 可以说是官网 sameple-controller 的运行示例.
2. [kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)
    - kuber官方工程
    - 实际上是[kubernetes/kubernetes](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/sample-controller)的子工程.
3. [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

> Custom resources are extensions of the Kubernetes API --来自参考文章3

我看了一下, CRD本身是kuber的一种资源类型, 与Pod, Service等平级. 

但某些时候也许这些资源太过通用, 不具有专业性, 比如如果某客户希望将机器学习部署在kuber集群中, 如何实现?

常规操作是构建镜像, 创建ConfigMap, 部署deployment及Service, 这些操作实际进行中将出现很多问题(我尝试把harbor从compose部署到kuber集群花了好多天).

CRD允许用户自定义资源类型, 如下, 使用`kubectl apply -f 文件名`即可创建一种名为`Foo`的资源类型, 同样与Pod, Service平级.

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: foos.samplecontroller.k8s.io
spec:
  group: samplecontroller.k8s.io
  version: v1alpha1
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
```

此处创建的Foo资源无namespace的限制, 全局通用, 毕竟Pod, Service也不是只能特定namespace才能使用的嘛.

然后, 使用如下配置创建Foo资源实例

```yaml
apiVersion: samplecontroller.k8s.io/v1alpha1
kind: Foo
metadata:
  name: example-foo
spec:
  deploymentName: example-foo
  replicas: 1
```

这里ta只定义了两个属性: `deploymentName`和`replicas`. 参考文章1和2中还有一个`crd-validation.yaml`用于声明对自定义资源属性字段的校验规则, 很容易理解, 这里就不写了.

这两个属性其实是与Foo资源的行为有关的(废话, 哪个资源的属性不是跟ta的行为有关?).

我们使用`kubectl apply -f 文件名`来创建Foo资源实例. 但是对于Pod, Service这种资源, kuber的controller manager是可以理解的, 并且按照不同资源规定的行为去处理. 

比如创建了一个Pod资源, controller manager需要保证pod已经正确在目标节点上运行, Pod实例的port, volume, network等属性也需要与之关系. 说到底, kuber代码中是有Pod这类结构体的存在的, yaml文件中可以定义的属性其实就是结构体的字段...不难理解吧?

但是Foo资源呢? kuber虽然允许你通过yaml文件自定义资源, 但是资源的行为仍然需要用代码来完成.

> 这样, 用户自定义资源可以看作是kuber的插件实现, 这下理解文章开头引用参考文章3的那句话了吧.

继续, 自定义资源作为kuber的一种扩展, 需要注册. 按照参考文章2中readme所说, 需要先build, 然后再执行.

```console
$ go build -o sample-controller .
./sample-controller -kubeconfig=$HOME/.kube/config
I1109 17:25:47.265085  110201 controller.go:227] Successfully synced 'default/example-foo'
I1109 17:25:47.265217  110201 event.go:281] Event(v1.ObjectReference{Kind:"Foo", Namespace:"default", Name:"example-foo", UID:"8fef5f20-e97e-45be-9ba8-f6c60e2a9b66", APIVersion:"samplecontroller.k8s.io/v1alpha1", ResourceVersion:"2259372", FieldPath:""}): type: 'Normal' reason: 'Synced' Foo synced successfully
```

这是一个前台进程, 不会退出 (但是资源部署完成后结束这个进程并不会影响已经存在的资源实例). 

ta会创建一个deployment对象, 其名称正是`deploymentName`字段所指定的名称, 并且该deployment的replicas值也正是上面Foo的replicas属性所表示的值, 于是也会创建1个pod对象.
