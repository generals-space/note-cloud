# metav1.Time类型与原生time.Time[转换 比较]

参考文章

1. [Convert CreationTimeStamp type to string](https://stackoverflow.com/questions/60069771/convert-creationtimestamp-type-to-string)

kuber: 1.17.3

`Pod.CreationTimeStamp`和`Pod.DeletionTimestamp`两个字段为`metav1.Time`类型.

> `metav1` -> `k8s.io/apimachinery/pkg/apis/meta/v1`

ta的原型很简单

```go
type Time {
    time.Time `protobuf:"-"`
}
```

kuber 对这个结构体又做了一个封装, 添加了很多自定义的方法.

我本来想把`Pod.CreationTimeStamp`使用`time.Time(Pod.CreationTimeStamp)`这种方式做一个强制类型转换, 然后和一个原生的`time.Time`对象(如`initTime`)做个比较的, 但是报错说不能这么用.

没办法, 暂时使用`metav1.Time{initTime}`将原生对象转换一下了.

...`metav1.Time`的`Before()`方法接受的是`metav1.Time`本身, 但是`After()`方法接受的却是`time.Time`对象, 这我是万万妹想到的.
