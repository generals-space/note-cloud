[kube-apiserver 的设计与实现](https://cloud.tencent.com/developer/article/1591764)

[《Kubernetes设计与实现》2.3 Kubernetes API](https://renhongcai.gitbook.io/kubernetes/di-er-zhang-kubernetes-ji-chu/2.3-kubernetes_api)

用于管理Kubernetes资源的API前缀除了api外，还有apis，在Kubernetes早期的设计中，只有api前缀，后来为了引入API组的设计，又添加了apis前缀，简单地说，使用api前缀的API是Kubernetes最核心的一组API，而使用apis前缀的API是Kubernetes发展过程中引入的分组API。

把API分组最大的好处在于用户可以自由地开启和关闭相应的非核心功能。用户可以使用kube-apiserver组件提供的--runtime-config参数来显式地控制某项功能是否开启。

Kubernetes API的分组设计为其提供了无限的扩展能力，借此机制可以轻松地为Kubernetes提供扩展功能，用户不仅可以使用CRD（Custom Resource Definition）功能来提供新的API，还可以通过扩展apiserver来扩展功能。

