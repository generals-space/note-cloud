# kuber-Pod钩子

参考文章

1. [Kubernetes容器上下文环境](https://www.cnblogs.com/zhenyuyaodidiao/p/6558444.html)
2. [Attach Handlers to Container Lifecycle Events](https://kubernetes.io/docs/tasks/configure-pod-container/attach-handler-lifecycle-event/)
3. [multiple command in postStart hook of a container](https://stackoverflow.com/questions/39436845/multiple-command-in-poststart-hook-of-a-container)
4. [sigma敏捷版系列文章：kubernetes grace period 失效问题排查](https://developer.aliyun.com/article/609813)
    - `PostStart`执行的时机是在容器启动以后，但是并不是等容器启动完成再执行(可以说`PostStart`和`Entrypoint`是并行执行的)

## 多个钩子 - 行内脚本的实现

```yaml
  lifecycle:
    postStart:
      exec:
        command:
          - "sh"
          - "-c"
          - >
            if [ -s /var/www/mybb/inc/config.php ]; then
            rm -rf /var/www/mybb/install;
            fi;
            if [ ! -f /var/www/mybb/index.php ]; then
            cp -rp /originroot/var/www/mybb/. /var/www/mybb/;
            fi
```

kube 在创建容器后立即发送`postStart`事件. 但是, **不能保证`postStart`处理程序在容器的`entrypoint`调用之前被调用**. 相对于容器的代码, `postStart`处理程序以**异步方式**运行, 但 kube 对容器的管理会阻塞直到`postStart`处理程序完成. 容器的状态直到`postStart`处理程序完成后才会设置为`RUNNING`. 

kube 在容器终止之前立即发送`preStop`事件. kube 对容器的管理一直阻塞直到`preStop`处理程序完成, 除非 Pod 的宽限期过期. 
