
参考文章

1. [Write permissions on volume mount with security context fsgroup option](https://discuss.kubernetes.io/t/write-permissions-on-volume-mount-with-security-context-fsgroup-option/16524)
    - `hostPath`与`nfs`, 并不兼容`securityContext`(`emptyDir`卷可以).
2. [Kubernetes: how to set VolumeMount user group and file permissions](https://stackoverflow.com/questions/43544370/kubernetes-how-to-set-volumemount-user-group-and-file-permissions)
3. [k8s中pod目录访问权限不足](https://blog.csdn.net/kylezhou1992/article/details/125452190)
4. [k8s pod挂载hostPath执行写时报错Permission denied](https://www.cnblogs.com/v-fan/p/16034960.html)
    - hostPath卷挂载规则DirectoryOrCreate: 如果给定的路径没有任何东西存在, 那将根据需要在此创建一个空目录, 权限设置为0755, 与kubelet拥有相同的用户与组
    - hostPath卷挂载规则FileOrCreate: 如果给定的路径没有任何东西存在, 那将根据需要在此创建一个空文件, 权限设置为0644, 与kubelet拥有相同的用户与组

## 场景描述

开发环境中, 容器中的进程一般都是通过 root 用户启动的, dockerfile不会额外添加`USER`指令, kube yaml中也不会添加`Pod.spec.securityContext.runAsUser`指定一个普通用户.

但是生产环境中一般都有安全规范, 禁止进程使用 root 权限运行, 所以需要进行改造.

在原本的dockerfile中, 在`CMD`指令前添加一行`USER`指令即可.

```dockerfile
USER nonroot
CMD ["xxx"]
```

有一种场景是, 容器内的进程要通过`hostPath`挂载的目录, 将日志/数据落盘保存. 但 kubernetes 映射到容器内部的目录, 属主仍然是 root 用户, 普通用户启动的进程, 没有这个目录的写入权限, 就会报错.

## 排查思路

最开始找到了参考文章2, 我本来想着, 能否找到让 kubernetes 在将 hostPath 映射到容器中时, 自动赋予指定的属主或权限, 比如

```yaml
kind: Pod
spec:
  securityContext:
    runAsUser: 1001
    runAsGroup: 1001
    fsGroup: 1001
```

> 这里假设`nonroot`用户的uid为1001.

希望`fsGroup`能达到我的目的, 但是没用. 

按照参考文章1中的解释, `hostPath`与`nfs`, 并不兼容`securityContext`, 无法使用后者指定的用户/组映射目录(`emptyDir`卷可以, 不过我的场景不能使用`emptyDir`)

参考文章2中采纳了`fsGroup`的答案, 是因为ta挂载的目录是`AWS`卷.

## 解决方案

最终选择了参考文章2中, 第2高票的答案, 这也是参考文章3中的解决方法. 即使用`initContainers`, 在实际的业务容器启动前, 先以root用户身份启动, 并将待挂载的目录的属主修改为业务容器的普通用户, 然后退出即可.

```yaml
kind: Pod
spec:
  initContainers:
  - name: init
    ## 这里最好用与业务容器相同的镜像, 为了保证该容器也存在`nonroot`用户.
    image: myredis:latest
    securityContext:
      ## 指定 root 用户执行
      runAsUser: 0
    command:
    - chown
    - -R
    - nonroot:nonroot
    - /var/lib/redis
    volumeMounts:
    - name: data
      mountPath: /var/lib/redis
  containers:
  - name: redis
    image: myredis:latest
    volumeMounts:
    - name: data
      mountPath: /var/lib/redis
  volumes:
  - name: data
    hostPath:
      path: /data
      type: DirectoryOrCreate
```

有两个规则:

1. init容器最好用与业务容器相同的镜像, 为了保证该容器也存在`nonroot`用户;
2. init容器中指定 root 用户执行, 以免 dockerfile 中已经用`USER`指定了`nonroot`普通用户运行;

------

但是还有一个问题, `initContainer`与`containers`并不存在网络/文件系统的共享机制, 毕竟ta们并不是同一个容器, 就算在`containers`下的多个容器之间, 也只是共享网络, 文件系统还是隔离的. 那么前者将目录属主修改为普通用户, 是在何处生效的呢?

其实在`init`修改`hostPath`(或是其子目录)后, 宿主机上该目录的属主也发生了变化, 双方存在的关联就是`uid`. 

假设容器内`nonroot`的`uid`为1001, 宿主机上可能并不存在这个用户, 但很有可能存在一个`uid`为1001的用户`nonroot2`, 不过没关系, 宿主机上的目录属主变成了`nonroot2`. 而业务容器挂载并映射该目录后, 会把其属主对应的`uid`1001, 当作`nonroot`.
