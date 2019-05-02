# Kubernetes-pod钩子

参考文章

1. [Kubernetes容器上下文环境](https://www.cnblogs.com/zhenyuyaodidiao/p/6558444.html)

2. [Attach Handlers to Container Lifecycle Events](https://kubernetes.io/docs/tasks/configure-pod-container/attach-handler-lifecycle-event/)

3. [multiple command in postStart hook of a container](https://stackoverflow.com/questions/39436845/multiple-command-in-poststart-hook-of-a-container)


## 多个钩子 - 行内脚本的实现

```yml
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