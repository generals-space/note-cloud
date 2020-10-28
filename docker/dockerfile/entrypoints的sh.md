# entrypoints的sh

参考文章

1. [一行代码的变更让我陷入无尽加班，Dockerfile的ENTRYPOINT的两种格式](https://www.pkslow.com/archives/docker-entrypoint-issue)
    - `entrypoints`两种格式: `exec`格式与`shell`格式
    - exec格式可以接受参数，而shell格式会忽略参数
    - shell格式相当于在前面还要再添加`/bin/sh -c`，所以app启动的进程ID不是1

对于ENTRYPOINT有两种格式：

exec格式（官方推荐使用）：

```
ENTRYPOINT ["executable", "param1", "param2"]
```

shell格式：

```
ENTRYPOINT command param1 param2
```

这两种不同的格式有一个很大的区别在于：exec格式可以接受参数，而shell格式是会忽略参数的。shell格式相当于在前面还要再添加`/bin/sh -c`，所以app启动的进程ID不是1。

------

参考文章1还是可以读一读的, 但是我遇到的问题有点不太一样.

在 docker 18.03.6 下, 构建如下 dockerfile.

```dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:latest
ENTRYPOINT ["tail", "-f", "/etc/profile"]
```

使用`docker history`子命令查看镜像层历史信息, 会发现为`ENTRYPOINT ["tail" "-f" "/etc/profile"]`

```
$ docker history 镜像名 --no-trunc
IMAGE                   CREATED             CREATED BY                                                      SIZE                COMMENT
sha256:4f47f76be97316   16 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["tail" "-f" "/etc/profile"]      0B
```

然而我在 arm64v8 平台下(docker 19.03.7), 构建同样的`dockerfile`, 其镜像层历史如下

```
$ docker history 镜像名 --no-trunc
IMAGE                   CREATED             CREATED BY                                                                              SIZE                COMMENT
sha256:4f47f76be97316   16 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["/bin/sh", "-c", "[\"tail\" \"-f\" \"/etc/profile\"]"]    0B
```

然后在启动容器的时候, 报`[tail`命令不存在(有个左括号)...对于构建出的`ENTRYPOINT`指令中直接出现`/bin/sh -c`我还是有点惊讶的.

我找到参考文章1, 尝试使用`shell`的格式重新构建.

```dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:latest
ENTRYPOINT tail -f /etc/profile
```

这次的结果为

```
$ docker history 镜像名 --no-trunc
IMAGE                   CREATED             CREATED BY                                                                          SIZE                COMMENT
sha256:4f47f76be97316   16 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["/bin/sh", "-c", "tail -f /etc/profile"]             0B
```

...还是有`/bin/sh -c`(和参考文章1中提到的, 正常情况下`shell`模式的构建结果完全相同了), 不过这次好在目标命令是可以执行成功的, 问题也算是解决了.

------

为了验证参考文章1中所说, 使用 shell 模式构建的镜像, 在启动容器时 pid 不为1的情况, 我试着启动了下, 发现进入容器后ps有如下结果

```
[root@0d80a96f4bf6 /]# ps -ef
UID         PID   PPID  C STIME TTY          TIME CMD
root          1      0  0 15:08 ?        00:00:00 tail -f /etc/profile
root          6      0  3 15:08 pts/0    00:00:00 bash
root         21      6  0 15:08 pts/0    00:00:00 ps -ef
```

...`tail`进程的pid明明就是1, 目前没想明白究竟是哪里和参考文章1有出入.

