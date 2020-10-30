参考文章

1. [kubelet 源码分析：启动流程](https://cizixs.com/2017/06/06/kubelet-source-code-analysis-part-1/)
2. [kubernetes grace period 失效问题排查](https://developer.aliyun.com/article/609813)
    - 值得收藏
    - `podManager`是`kubelet`在本地缓存`Pod`信息的数据结构, 是比较核心的组件
    - kubelet 获取 Pod 的信息有几个途径：
        1. 通过 podManager 获取本地实时的 Pod 信息
        2. 通过 kube-apiserver 获取 Pod 信息
        3. 静态 Pod 通过配置文件或者 url 获取
        4. 通过本地 Container 列表反向演算 Pod 信息(kubelet 重启的时候就是通过这种方式判断本机上有哪些 Pod 的)
        5. 通过 cgroup 演算 Pod 信息（kubelet 就是通过对比 cgroup 设置的信息和 podManager 的信息判断哪些 Pod 是孤儿 Pod 的，孤儿 Pod 的进程会被 kubelet 立即通过 kill 信号杀进程，非常强暴）

kubelet 重启会导致通过 kubectl exec 进入容器的连接断开.
