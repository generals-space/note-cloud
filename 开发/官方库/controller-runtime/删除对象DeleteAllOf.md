参考文章

1. [Kubebuilder v1 版本 VS v2 版本](https://cloudnative.to/kubebuilder/migration/v1vsv2.html)
2. [controller-runtime/pkg/client/example_test.go](https://github.com/kubernetes-sigs/controller-runtime/blob/00af7b6464ff24864470171799a2755032fa4499/pkg/client/example_test.go#L234)

在不预先查出一个对象列表时, 直接删除



```go
// controller-runtime:pkg/client/example_test.go

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)
// This example shows how to use the client with typed and unstrucurted objects to delete collections of objects.
func ExampleClient_deleteAllOf() {
	// Using a typed object.
	// c is a created client.
	_ = c.DeleteAllOf(context.Background(), &corev1.Pod{}, client.InNamespace("foo"), client.MatchingLabels{"app": "foo"})

	// Using an unstructured Object
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps",
		Kind:    "Deployment",
		Version: "v1",
	})
	_ = c.DeleteAllOf(context.Background(), u, client.InNamespace("foo"), client.MatchingLabels{"app": "foo"})
}
```
