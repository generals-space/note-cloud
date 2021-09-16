# kubernetes版本与CAP_配置的格式[memlock CAP_SYS_RESOURCE]

参考文章

1. [Capabilities in security context need to be specified differently for docker vs rkt](https://github.com/kubernetes/kubernetes/issues/33104)

## 问题描述

之前已经有了 es.v7.5.1 arm 平台镜像, 后来要求将该版本迁移到 x86 平台.

es v7 中`elasticsearch.yml`新增了一条`bootstrap.memory_lock: true`, 需要同时修改`limits.conf`.

```
* soft nofile 65536
* hard nofile 65536
* soft nproc 65535
* hard nproc 65535
* soft memlock unlimited
* hard memlock unlimited
```

但是该镜像在部署时, 一直无法启动, 显示`CrashLoopBackOff`.

## 排查思路

查看容器日志, 发现启动脚本中的`runuser -u elasticsearch elasticsearch`命令有如下输出

```
runuser: cannot open session: Permission denied
```

将`command`修改为`tail -f /etc/profile`, 容器启动后进入容器终端, 发现不只`runuser`, 连`su`命令都没法用了.

```
[root@esc-210915-143151-data-1 elasticsearch]# su -l elasticsearch
Last login: Wed Sep 15 15:50:44 CST 2021 on pts/0
su: cannot open session: Permission denied
```

后查明是因为修改`memlock`需要容器有`CAP_SYS_RESOURCE`能力, 算是个低级错误.

但是在部署文件中添加如下配置后.

```yaml
        securityContext:
          capabilities:
            add:
            - CAP_SYS_RESOURCE
          privileged: false
          procMount: Default
```

容器启动仍然失败

```
esc-210915-143151-data-0       1/2     RunContainerError   15         57m     192.168.31.244   ly-xjf-r021110-gyt   <none>           <none>
esc-210915-143151-data-1       1/2     RunContainerError   16         57m     192.168.31.135   ly-xjf-r020905-gyt   <none>           <none>
esc-210915-143151-data-2       1/2     CrashLoopBackOff    16         57m     192.168.31.137   ly-xjf-r020803-gyt   <none>           <none>
esc-210915-143151-exporter-0   1/1     Running             0          4h57m   192.168.31.228   ly-xjf-r021110-gyt   <none>           <none>
esc-210915-143151-kibana-0     0/1     CrashLoopBackOff    62         4h57m   192.168.31.107   ly-xjf-r020901-gyt   <none>           <none>
esc-210915-143151-master-0     1/2     CrashLoopBackOff    40         3h2m    192.168.31.25    ly-xjf-r021110-gyt   <none>           <none>
esc-210915-143151-master-1     1/2     CrashLoopBackOff    40         3h2m    192.168.31.83    ly-xjf-r020803-gyt   <none>           <none>
esc-210915-143151-master-2     1/2     CrashLoopBackOff    40         3h2m    192.168.31.38    ly-xjf-r020905-gyt   <none>           <none>
```

这次的报错变成了`RunContainerError`, 而不在是单纯的`CrashLoopBackOff`, 查看日志时有如下输出.

```
[monitor@ly-xjf-r020807-gyt v1.crd]$ k logs -f esc-210915-143151-data-0 -c es-cluster
failed to open log file "/var/log/pods/b8012435-15fb-11ec-b6be-000af79b1e70/es-cluster/9.log": open /var/log/pods/b8012435-15fb-11ec-b6be-000af79b1e70/es-cluster/9.log: no such file or directory
```

到这些容器所在的宿主机上, 用`docker ps -a`查看, 发现有容器启动, 但立刻就结束了, `docker log`也没有任何输出.

...md, 好像问题变得更复杂了😒

## 解决方法

因为在之前的排查过程中, 确认了`su -l`失败就是因为`limits.conf`中添加了`memlock`那两行的问题. 而且如果将`privileged`直接设置为`true`, 不再单独设置`CAP_`字段, 也能够让容器正常启动.

所以现在的目标就是搜索, 除了`CAP_SYS_RESOURCE`, `memlock`是不是还需要其他的内核能力.

ヽ｀、ヽ｀｀、ヽ｀ヽ｀、、ヽ ｀ヽ 、ヽ｀｀ヽヽ｀ヽ、ヽ｀ヽ｀、ヽ｀｀、ヽ 、｀｀、 ｀、ヽ｀  ｀ ヽ｀ヽ、ヽ ｀、ヽ｀｀、ヽ、｀｀、｀、ヽ｀｀、 、ヽヽ｀、｀、、ヽヽ、｀｀😭、 、 ヽ｀、ヽ｀｀、ヽ｀ヽ｀、、ヽ ｀ヽ 、ヽ｀｀ヽヽ｀ヽ、ヽ｀ヽ｀、ヽ｀

经过漫长的搜索, 恰好发现了参考文章1, 其中提到了2种不同的`capabilities`的格式. 该issue中的kubernetes对于`rkt`的运行时, 权限配置是这样的

```yml
      securityContext:
        capabilities:
          add: ["CAP_NET_ADMIN"]
```

而对于`docker`运行时, 权限配置则需要是这样的

```yml
       securityContext:
        capabilities:
          add: ["NET_ADMIN"]
```

后者的`CAP_`前缀没有了, 这时我才突然想起来, `arm`平台与`x86`平台的kubernetes版本是不同的, `arm`平台的版本为1.17.2, `x86`的则是`1.13.2`.

同时在上面出现`RunContainerError`的容器所在主机的`/var/log/message`中, 发现与参考文章1中提到的异常日志

```
Sep 15 19:37:07 ly-xjf-r021110-gyt dockerd: time="2021-09-15T19:37:07.025393994+08:00" level=error msg="Handler for POST /v1.38/containers/21c865364925ff782dc85ab8773d12c65fcf93e7dde31f3160bad26aee3cfa21/start returned error: linux spec capabilities: Unknown capability to add: \"CAP_CAP_SYS_RESOURCE\""
Sep 15 19:37:07 ly-xjf-r021110-gyt kubelet: E0915 19:37:07.050323  189742 remote_runtime.go:213] StartContainer "21c865364925ff782dc85ab8773d12c65fcf93e7dde31f3160bad26aee3cfa21" from runtime service failed: rpc error: code = Unknown desc = failed to start container "21c865364925ff782dc85ab8773d12c65fcf93e7dde31f3160bad26aee3cfa21": Error response from daemon: linux spec capabilities: Unknown capability to add: "CAP_CAP_SYS_RESOURCE"
Sep 15 19:37:07 ly-xjf-r021110-gyt kubelet: E0915 19:37:07.050402  189742 kuberuntime_manager.go:749] container start failed: RunContainerError: failed to start container "21c865364925ff782dc85ab8773d12c65fcf93e7dde31f3160bad26aee3cfa21": Error response from daemon: linux spec capabilities: Unknown capability to add: "CAP_CAP_SYS_RESOURCE"
Sep 15 19:37:07 ly-xjf-r021110-gyt kubelet: E0915 19:37:07.050439  189742 pod_workers.go:190] Error syncing pod b8e6f21e-15fe-11ec-b6be-000af79b1e70 ("esc-210915-143151-master-0_zjjpt-es(b8e6f21e-15fe-11ec-b6be-000af79b1e70)"), skipping: failed to "StartContainer" for "es-cluster" with RunContainerError: "failed to start container \"21c865364925ff782dc85ab8773d12c65fcf93e7dde31f3160bad26aee3cfa21\": Error response from daemon: linux spec capabilities: Unknown capability to add: \"CAP_CAP_SYS_RESOURCE\""
```

😣可恶, 当初我也查看过这个文件的内容, 竟然没发现.

于是我将`capabilities`的配置改为如下

```
        securityContext:
          capabilities:
            add:
            - SYS_RESOURCE
          privileged: false
          procMount: Default
```

然后就可以了.
