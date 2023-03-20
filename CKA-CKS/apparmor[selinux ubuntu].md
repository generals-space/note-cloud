
<!--
<!key!>: {8ac2d6d9-b4a7-4551-87a0-cb2af8f0c81f}
<!link!>: {note-devops:a4fdeb45-7d50-48ac-851b-015c97eb340d}
-->

参考文章

1. [云原生|kubernetes|apparmor的配置和使用](https://blog.csdn.net/alwaysbefine/article/details/128661805)
2. [使用 AppArmor 限制容器对资源的访问](https://kubernetes.io/zh-cn/docs/tutorials/security/apparmor/)
    - 官方文档

## 引言

APPARMor是一个Linux操作系统的内核安全模块，和selinux功能基本是一致的。

Ubuntu和openSUSE才有AppAromor，centos，Redhat这些操作系统是使用SeLinux的。

## 在 kube 环境中应用

apparmor的应用对象为 Pod, 可以让Pod引用其所在主机上`/etc/apparmor.d/`目录下中配置定义的profile规则.

以下规则没有定义匹配的服务, 将会应用到所有进程.

```conf
# 文件名: deny.conf(名称随意)

#include <tunables/global>

profile k8s-apparmor-example-deny-write flags=(attach_disconnected) {
  #include <abstractions/base>

  file,

  # 拒绝所有文件写入
  deny /** w,
}
```

使用如下命令加载(必须)

```
apparmor_parser /etc/apparmor.d/deny.conf
```

> 由于Pod被调度到哪台主机并不是固定的, 所以这些规则需要在集群中各个节点上都存在, 且都加载.

```pod
apiVersion: v1
kind: Pod
metadata:
  name: hello-apparmor
  annotations:
    # 告知 Kubernetes 去应用 AppArmor 配置 "k8s-apparmor-example-deny-write"。
    # 请注意，如果节点上运行的 Kubernetes 不是 1.4 或更高版本，此注解将被忽略。
    container.apparmor.security.beta.kubernetes.io/hello: localhost/k8s-apparmor-example-deny-write
spec:
  containers:
  - name: hello
    image: busybox
    command: [ "sh", "-c", "echo 'Hello AppArmor!' && sleep 1h" ]
```

注解中配置格式如下

```
container.apparmor.security.beta.kubernetes.io/container名称: localhost/profile名称
```

## 检测生效

进入目标容器, 尝试写入文件, 会发现会失败.

```console
controlplane $ k exec -it hello-apparmor sh
/ # touch abc
touch: abc: Permission denied
/ # echo test >> /etc/hosts 
sh: can't create /etc/hosts: Permission denied
```

> 不过这个`deny.conf`只在容器中生效, 而在宿主机上, 虽然已经使用`apparmor_parser`加载了该配置, 但是并没有限制用户的写入行为, 为什么???

## FAQ

### apparmor profile not found

```console
controlplane $ k get pod -owide
NAME             READY   STATUS                 RESTARTS   AGE   IP            NODE     NOMINATED NODE   READINESS GATES
hello-apparmor   0/1     CreateContainerError   0          5s    192.168.1.3   node01   <none>           <none>
controlplane $ k describe pod hello-apparmor | tail
Events:
  Type     Reason     Age               From               Message
  ----     ------     ----              ----               -------
  Warning  Failed     2s (x3 over 14s)  kubelet            Error: failed to get container spec opts: failed to generate apparmor spec opts: apparmor profile not found k8s-apparmor-example-deny-write
```

pod所在主机上, `/etc/apparmor.d/`目录下必须存在目标`profile`, 且需要使用`apparmor_parser`事先加载.
