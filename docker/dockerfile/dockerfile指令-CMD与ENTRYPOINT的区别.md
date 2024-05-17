# dockerfile指令-CMD与ENTRYPOINT的区别

参考文章

1. [一行代码的变更让我陷入无尽加班，Dockerfile的ENTRYPOINT的两种格式](https://www.pkslow.com/archives/docker-entrypoint-issue)
    - `entrypoints`两种格式: `exec`格式与`shell`格式
    - exec格式可以接受参数，而shell格式会忽略参数
    - shell格式相当于在前面还要再添加`/bin/sh -c`，所以app启动的进程ID不是1
2. [Demystifying ENTRYPOINT and CMD in Docker](https://aws.amazon.com/cn/blogs/opensource/demystifying-entrypoint-cmd-docker/)
    - docker 会自动将字符串形式的`ENTRYPOINT`与`CMD`转换成数组形式. 比如
    - `ENTRYPOINT /usr/bin/httpd -DFOREGROUND` -> `["/bin/sh", "-c", "/usr/bin/httpd -DFOREGROUND"]`

## 格式

`ENTRYPOINT`和`CMD`有两种格式

### exec格式(官方推荐使用)

```dockerfile
CMD ["executable", "param1", "param2"]
ENTRYPOINT ["executable", "param1", "param2"]
```

> 列表中必须为双引号

### shell格式

```dockerfile
CMD command param1 param2
ENTRYPOINT command param1 param2
```

## 两种格式的区别

> `ENTRYPOINT`与`CMD`同时使用时, 只能都用列表形式(没试过都用shell命令格式的, 但一个列表一个shell命令的肯定不行)

按照参考文章1所说, docker 会自动将字符串形式的`ENTRYPOINT`与`CMD`转换成数组形式.

`ENTRYPOINT ["/usr/bin/httpd", "-DFOREGROUND"]` -> `["/usr/bin/httpd", "-DFOREGROUND"]`

`ENTRYPOINT /usr/bin/httpd -DFOREGROUND` -> `["/bin/sh", "-c", "/usr/bin/httpd -DFOREGROUND"]`

`CMD ["tail" "-f" "/etc/os-release"]` -> `["tail" "-f" "/etc/os-release"]`

`CMD tail -f /etc/os-release` -> `["/bin/sh" "-c" "tail -f /etc/os-release"]`

exec格式与shell格式有一个很大的区别在于: `exec`格式可以接受参数, 而`shell`格式是会忽略参数的. 

比如如下 dockerfile

```dockerfile
FROM centos:7
ENTRYPOINT ["tail", "-f", "/etc/os-release"]
```

构建镜像并按如下方式启动

```
docker build -f dockerfile -t centos7:test .
docker run -d centos7:test tail -f /etc/hostname
```

容器内的进程如下

```console
[root@deabb9078762 ~]# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 07:16 ?        00:00:00 tail -f /etc/os-release tail -f /etc/hostname
```

但是如果 dockerfile 的内容为

```dockerfile
FROM centos:7
ENTRYPOINT tail -f /etc/os-release
```

那么在用与上面同样的方式构建镜像, 启动容器后, 容器内的进程将是这样的

```console
[root@deabb9078762 ~]# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 07:16 ?        00:00:00 tail -f /etc/os-release
```

你会发现`docker run`末尾指定的`tail -f /etc/hostname`根本不生效...

## shell 格式容器进程 pid 不为 1 ?

shell格式相当于在前面还要再添加`/bin/sh -c`, 所以app启动的进程ID不是1.

为了验证参考文章1中所说, 使用 shell 模式构建的镜像, 在启动容器时 pid 不为1的情况, 我试着启动了下, 发现进入容器后ps有如下结果

```console
[root@0d80a96f4bf6 /]# ps -ef
UID         PID   PPID  C STIME TTY          TIME CMD
root          1      0  0 15:08 ?        00:00:00 tail -f /etc/profile
root          6      0  3 15:08 pts/0    00:00:00 bash
root         21      6  0 15:08 pts/0    00:00:00 ps -ef
```

...`tail`进程的pid明明就是1, 目前没想明白究竟是哪里和参考文章1有出入.

## shell 格式的优势

shell命令格式的`CMD`与`ENTRYPOINT`, 可以在其中通过`$变量名`引用`ENV`或是`-e`选项指定的环境变量, 但是数组形式不可以, 它会把`$变量名`当成字符串处理.

不过即使这样, 使用`ENTRYPOINT ["/docker-entrypoint.sh"]`, 然后在`/docker-entrypoint.sh`脚本中可以同时得到`CMD`参数, 也可以获取环境变量信息, 感觉更为灵活一点.
