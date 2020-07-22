# pod-qos服务质量[qosClass request limit burstable]

参考文章

1. [kubernetes 之QoS服务质量管理](https://www.cnblogs.com/tylerzhou/p/11043280.html)
2. [kubenetes之配置pod的QoS](https://www.cnblogs.com/tylerzhou/p/11043282.html)
3. [Configure Quality of Service for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/)
4. [Configure Out of Resource Handling](https://kubernetes.io/docs/tasks/administer-cluster/out-of-resource/)

## 1. 引言

`QoS`: Quality of Service, 服务质量.

在kubernetes中，每个POD都有个QoS标记，通过这个Qos标记来对POD进行服务质量管理。它取决于用户对服务质量的预期，也就是期望的服务质量。

简单来说就是, 当你对一个 Pod 的性能表现的期望比较确定, 希望为ta分配指定大小的资源(也许在这个资源配比下项目组经过了严密的压力测试, 性能达标), 那么`kuber`会为你的服务提供更加可靠的资源保障. 而很多时候, 一些不那么重要的, 测试性的 Pod 没有那么多规则, 一般都不会指定ta的`resources.{requests,limits}`, 那么`kuber`可能会在资源紧张的场景下挤压(甚至说掠夺)这些Pod的资源.

当你指定的资源分配越笃定, 越规范, kuber就越会保障这样的 Pod 更好地运行.

按照资源指定的严密性, kuber中的`QoS`分为3个级别: `Guaranteed` > `Burstable` > `BestEffort`.

kuber中只有Pod拥有`QoS`属性, 且体现在两个指标: CPU和内存.

## 2.
