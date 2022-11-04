# runtime.Object通用接口使用.3.List

参考文章

- kube: 1.16.2
- apimachinery: v0.17.0

apimachinery@v0.17.2/pkg/apis/meta/v1/types.go -> List{}

apimachinery@v0.17.2/pkg/apis/meta/v1/types.go -> ObjectMeta{}

apimachinery@v0.17.2/pkg/apis/meta/v1/types.go -> TypeMeta{}


apimachinery@v0.17.2/pkg/api/meta/help.go -> ExtractList(obj runtime.Object)

将 runtime.Object 转换成 []List 列表对象.
