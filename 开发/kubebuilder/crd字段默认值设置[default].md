
```go
type Reference struct {
	Name string `json:"name"`
}

type ResourceSpec struct {
    // 简单字段默认值
    // +kubebuilder:default=Delete
	DeletionPolicy string `json:"deletionPolicy,omitempty"`

    // 结构体默认值
	// +kubebuilder:default={"name": "default"}
	ProviderConfigReference *Reference `json:"providerConfigRef,omitempty"`
}

```

这样通过 make 生成的 crd yaml 如下.

```yaml
              deletionPolicy:
                default: Delete
                description: ''
                enum:
                - Orphan
                - Delete
                type: string

              providerConfigRef:
                default:
                  name: default
                description: ''
                properties:
```
