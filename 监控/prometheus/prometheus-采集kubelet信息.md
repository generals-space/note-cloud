# prometheus-采集kubelet信息

参考文章

1. [从Kubelet获取节点运行状态](https://yunlzheng.gitbook.io/prometheus-book/part-iii-prometheus-shi-zhan/readmd/use-prometheus-monitor-kubernetes#cong-kubelet-huo-qu-jie-dian-yun-hang-zhuang-tai)
    - 两种方式从node采集信息: 1. 直接访问kubelet暴露在宿主机上的端口; 2. 通过Kubernetes的apiserver提供的代理API访问各个节点中kubelet的metrics服务

各节点的kubelet组件中除了包含自身的监控指标信息以外，kubelet组件还内置了对cAdvisor的支持。cAdvisor能够获取当前节点上运行的所有容器的资源使用情况.
