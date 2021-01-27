# kubectl logs -p 查询pod中container重启前的日志

参考文章

1. [kubernetes查看重启pod的日志](https://blog.csdn.net/u012517061/article/details/108155670)

有时 Pod 内的 container 因为 CPU/内存 过高, 导致了系统 OOM 或是健康检查多次无响应后被 kill 掉.

而重启的 container 在使用`kubectl logs -f Pod名称`时是没有办法看到重启前的日志的, 这就导致无法排查原因.

此时可以使用`kubectl logs -p Pod名称`, `-p/--previous`, 只要Pod本身还存在, 就可以查看 container 之前的日志信息.

> 疑问: container 多次重启, `-p`是可以查看重启之前一次的日志, 还是可以查看从 Pod 启动开始的所有日志???

