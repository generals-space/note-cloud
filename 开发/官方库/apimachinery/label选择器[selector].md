# label选择器[selector]

```go
require (
	k8s.io/apimachinery v0.17.4
	k8s.io/klog v1.0.0
)
```

## label 构造方法

在写 controller 时, 经常遇到需要根据 label 操作指定 pod 的时候, 官方库`apimachinery`提供了相关的方法(不只有`label`, 还有`field`).

```go
package main

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/klog"
)

func main() {
	labelMap := map[string]string{
		"key01": "val01",
		"key02": "val02",
	}
	selector01 := labels.SelectorFromSet(labels.Set(labelMap))
	klog.Infof("%s", selector01.String()) // key01=val01,key02=val02

	selector := labels.NewSelector()
	require01, _ := labels.NewRequirement("key03", selection.DoesNotExist, nil)
	require02, _ := labels.NewRequirement("key04", selection.In, []string{"val04"})
	selector02 := selector.Add(*require01, *require02)
	klog.Infof("%s", selector02.String()) // !key03,key04 in (val04)
}
```

> `selection`包下就一个单文件, 列举了一下可用的操作类型, 也就 7, 8 个的样子.

`labels.SelectorFromSet()`的使用方法虽然简单, 但是貌似只能指定`key=val`这种情况, 对于`!=`, `in`, `not in`这些操作就无能为力了. 

我在相应的`_test.go`文件中也没有找到相关的示例, 只能使用`Requirement`对象.

## 使用方法

```go
	client, _ := clientset.NewForConfig(cfg)

	// 不能简单地用Get("api-server"), 因为Pod名称中还会加一些额外字符串, 比如hostname.
	// 之前用的是 metav1.LabelSelector{}, 然后转换成String(), 但那只是简单的Marshal(),
	// 实际上应该使用 labels.Set{} 结构
	labelSet := labels.Set{
		"component": "kube-apiserver",
	}
	podListOpts := metav1.ListOptions{
		LabelSelector: labelSet.String(),
	}
	podList, err := client.CoreV1().Pods("kube-system").List(podListOpts)
```

------

不过其实可以直接写 label string, 不需要用 label 对象做转换.

```go
	client, _ := clientset.NewForConfig(cfg)
	podList, err := client.CoreV1().Pods("kube-system").List(
		metav1.ListOptions{
		LabelSelector: "key01=val01,!key02",
	})
```
