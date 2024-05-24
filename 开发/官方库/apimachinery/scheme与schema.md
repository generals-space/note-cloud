```go
k8s.io/client-go/kubernetes/scheme

scheme.Codecs.UniversalDecoder()
```

```
k8s.io/apimachinery/pkg/runtime/schema

schema.GroupVersion{}
schema.GroupResource
schema.ParseGroupVersion(()
```

另外

```go
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/runtime" // 这个包中也有一个 Scheme 成员
```

`runtime.Scheme`与`scheme.Scheme`可以相互赋值.

