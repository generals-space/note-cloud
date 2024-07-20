# kube演进历程

参考文章

1. [使用CoreDNS实现Kubernetes基于DNS的服务发现](https://www.kubernetes.org.cn/4694.html)
    - dns组件的演进: skydns(1.3之前) -> kubedns(1.3) -> coredns(1.11) 可能是因为 CoreDNS 支持 IPv6?
2. [从kubectl top看K8S监控](https://www.jianshu.com/p/64230e3b6e6c)
    1. k8s 1.6开始, kubernetes将cAdvisor开始集成在kubelet中, 不需要单独配置
    2. k8s 1.7开始, Kubelet metrics API 不再包含 cadvisor metrics, 而是提供了一个独立的 API 接口来做汇总
    3. k8s 1.12开始, cadvisor 监听的端口在k8s中被删除, 所有监控数据统一由Kubelet的API提供
3. [kubernetes 权威指南第5版 5.3]()
    1. v1.2版本引入 Scheduler Extender, 支持外部扩展;
    2. v1.5版本为调度器的优先级算法引入 Map/Reduce 的计算模式;
    3. v1.15版本实现了基于 Scheduling Framework 的方式, 开始支持组件化开发;
    4. v1.18版本将所有策略(Predicates与Priorities)全部组件化, 将默认的调度流程切换为 Scheduling Framework;
    5. v1.19版本将抢占过程组件化, 同时支持 Multi Scheduling Profile;

