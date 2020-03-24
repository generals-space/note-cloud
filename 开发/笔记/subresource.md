参考文章

1. [kubernetes 审计日志功能](https://www.jianshu.com/p/8117bc2fb966)

2. [Kubernetes API 概念](https://k8smeetup.github.io/docs/reference/api-concepts/)

kubernetes 每个资源对象都有 subresource,通过调用 master 的 api 可以获取 kubernetes 中所有的 resource 以及对应的 subresource,比如 pod 有 logs、exec 等 subresource。 --参考文章1

一些资源类型将具有一个或多个子资源，在资源下方表示为子路径：`GET /apis/GROUP/VERSION/RESOURCETYPE/NAME/SUBRESOURCE` --参考文章2

