参考文章

1. [Kubernetes controllers for CRDs的示例及Rancher RKE的使用](https://www.cnblogs.com/hindsight/p/9036362.html)

2. [Prometheus Operator的工作原理](https://yunlzheng.gitbook.io/prometheus-book/part-iii-prometheus-shi-zhan/operator/what-is-prometheus-operator#prometheus-operator-de-gong-zuo-yuan-li)
    - 描述了声明式API的涵义.
3. [kubernetes/apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver)
    - 用于管理`CRD`资源的包(通过`client-go`没有针对CRD的实现)
4. [通过自定义资源扩展Kubernetes](https://blog.gmem.cc/crd)

- `CRD`: Custom Resource Definition 自定义资源对象
- `CR`: Custom Resource CRD的实例
- `CRD Controller`: 执行 CRD 实际逻辑的代码部分

1. 没有对应的 controller, 也可以创建 crd 对象, 并且创建 cr 实例, 只是没有对应的 controller 去处理罢了.
2. 只有 controller, 没有 crd 也是可以的, 有两种情况:
    1. 只处理 k8s 内置的资源, 不需要创建 crd golang 对象(有 KVG, Spec 和 Status 等字段), 这种可以见 informer 示例, 可以说这是最简的 controller 了.
    2. 需要创建 crd golang 对象 (KVG, Spec, Status)

只创建 CRD golang 结构体, 能否在 controller 创建? (这种肯定无法通过 kubectl 获取 crd 对象列表的).
