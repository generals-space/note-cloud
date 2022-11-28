# kuber-ConfigMap更新之后[cm]

参考文章

1. [ConfigMap热更新](https://blog.csdn.net/Cui_Cui_666/article/details/105620445)
    - ENV 是在容器启动的时候注入的，启动之后 kubernetes 就不会再改变环境变量的值，且同一个 namespace 中的 pod 的环境变量是不断累加的。
    - 但是Volume不同，kubelet的源码里KubeletManager中是有Volume Manager的，这就说明Kubelet会监控管理每个Pod中的Volume资源，当发现配置的Volume更新后，就会重建Pod，以更新所用的Volume，但是会有一定延迟，大概10秒以内。
2. [Kubernetes - Configmap热更新原理](https://blog.csdn.net/qingyafan/article/details/102848860)
    - configmap具备热更新的能力，但只有通过目录挂载的configmap才具备热更新能力，其余通过**环境变量**，通过**subPath挂载的文件**都不能动态更新。
    - kubelet有一个启动参数`--sync-frequency`，控制同步配置的时间间隔，它的默认值是1min，所以更新configmap的内容后，真正容器中的挂载内容变化可能在0~1min之后。

`ConfigMap`通过`volume`挂载入`Pod`, 然后更新`ConfigMap`中的信息后, `Pod`内部的`ConfigMap`是会同步变动的, 但是由于Pod内的进程没有重启, 所以大部分场景还是需要重启一下Pod才会生效.

但是, 这要求挂载方式为 volume, 且为目录类型, 即不可以为`subPath`形式.

注意内容同步的延迟时间.
