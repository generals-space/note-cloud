# 容器进程与宿主机PID的映射关系

参考文章

1. [在docker宿主机上查找指定容器内运行的所有进程的PID](https://www.cnblogs.com/keithtt/p/7591097.html)
    - `/sys/fs/cgroup/memory/docker/${containerID}/cgroup.procs`
2. [PID mapping between docker and host](https://stackoverflow.com/questions/33328841/pid-mapping-between-docker-and-host)
    - `/proc/${宿主机PID}/status`, 不过这是通过宿主机`PID`得到`namespace id`以及容器内该进程`pid`的方式, 作用有点鸡肋.
3. [在Docker容器中运行的进程的主机中的PID是什么？](http://codingdict.com/questions/44979)
    - 与参考文章2讲和是同一个东西
4. [Kubernetes 教程：根据 PID 获取 Pod 名称](https://zhuanlan.zhihu.com/p/164421055)
    - `/proc/${PID}/cgroup`存储着pod UID
    - `/proc/${PID}/mountinfo`存储着pod UID
5. [Kubernetes 教程：根据 PID 获取 Pod 名称](https://www.cnblogs.com/ryanyangcs/p/13384118.html)
    - 同参考文章4

在宿主机上使用`ps -ef`可以看到该主机上所有容器内的进程信息, 但是没办法和指定容器对应起来, 这篇文章就是讨论如何在宿主机上, 通过进程`PID`查找一个其所在的容器的方法.

## 1. `cgroup.procs`文件

按照参考文章1所说, 可以使用`cgroup`中的`cgroup.procs`来查找.

对于裸docker容器, 可以直接查看`/sys/fs/cgroup/memory/docker/${containerID}/cgroup.procs`的文件内容, 如下

```console
$ cat cgroup.procs
58384
61533
$ ps -ef | grep 58384
root      58384  58365  0 15:36 ?        00:00:01 tail -f /etc/profile
$ ps -ef | grep 61533
root      61533  58365  0 18:45 pts/0    00:00:00 sh
```

而对于kuber中的容器, 则可能还需要先查找一下 Pod 的 UID.

1. `Guaranteed`: `/sys/fs/cgroup/memory/kubepods.slice/kubepods-pod${podUID}.slice/docker-${containerID}.scope/cgroup.procs`
2. `Burstable`: `/sys/fs/cgroup/memory/kubepods.slice/kubepods-burstable.slice/kubepods-pod${podUID}.slice/docker-${containerID}.scope/cgroup.procs`
3. `BestEffort`: `/sys/fs/cgroup/memory/kubepods.slice/kubepods-besteffort.slice/kubepods-pod${podUID}.slice/docker-${containerID}.scope/cgroup.procs`

## 2. docker top

使用`docker top`可以查看容器内的进程列表, 而且进程`PID`是按宿主机的`PID`显示的.

```console
$ docker top 470a404d9538
UID                 PID                 PPID                C                   STIME               TTY                 TIME                CMD
root                58384               58365               0                   15:36               ?                   00:00:01            tail -f /etc/profile
root                61533               58365               0                   18:45               pts/0               00:00:00            sh
```

------

如果你仔细观察进程`PID`与`PPID`的话, 就会发现ta们是基于同一个`PPID`创建的.

```console
$ ps -ef| grep 58365
root      58365   1285  0 15:36 ?        00:00:01 containerd-shim -namespace moby -workdir /var/lib/containerd/io.containerd.runtime.v1.linux/moby/470a404d95381d66c3eced7c02594a14589c66849fac50ef5332d74d3f83f9ae -address /run/containerd/containerd.sock -containerd-binary /usr/bin/containerd -runtime-root /var/run/docker/runtime-runc -systemd-cgroup
root      58384  58365  0 15:36 ?        00:00:01 tail -f /etc/profile
root      61533  58365  0 18:45 pts/0    00:00:00 sh
```

这样, 我们可以反推, 得到获取容器进程的父进程的方法.
