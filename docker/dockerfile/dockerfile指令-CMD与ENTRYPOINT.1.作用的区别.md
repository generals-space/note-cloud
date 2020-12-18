# dockerfile指令-CMD与ENTRYPOINT.1.作用的区别

参考文章

1. [一行代码的变更让我陷入无尽加班，Dockerfile的ENTRYPOINT的两种格式](https://www.pkslow.com/archives/docker-entrypoint-issue)
    - `entrypoints`两种格式: `exec`格式与`shell`格式
    - exec格式可以接受参数，而shell格式会忽略参数
    - shell格式相当于在前面还要再添加`/bin/sh -c`，所以app启动的进程ID不是1
2. [Demystifying ENTRYPOINT and CMD in Docker](https://aws.amazon.com/cn/blogs/opensource/demystifying-entrypoint-cmd-docker/)
    - docker 会自动将字符串形式的`ENTRYPOINT`与`CMD`转换成数组形式. 比如
    - `ENTRYPOINT /usr/bin/httpd -DFOREGROUND` -> `["/bin/sh", "-c", "/usr/bin/httpd -DFOREGROUND"]`

## 引言 - dockerfile 的 CMD 指令与 docker run

我们知道, 在执行`docker run`时, 有两种形式可以启动容器

```
docker run -d centos:7 tail -f /etc/os-release
docker run -it centos:7 /bin/bash
```

上述两条命令分别表示`-d`(服务模式)和`-it`(交互模式)的启动方式.

其中`tail -f /etc/os-release`和`/bin/bash`其实就相当于`dockerfile`中`CMD`指令的作用.

比如, 如下两种方式启动容器就是等价的.

**第1种**

```dockerfile
docker run -d centos7:test tail -f /etc/os-release
```

**第2种**

```dockerfile
FROM centos:7
CMD tail -f /etc/os-release
```

```
## 构建为 centos7:test 镜像
docker build -f dockerfile -t centos7:test .
docker run -d centos7:test
```

------

如果`dockerfile`中存在`CMD`指令, 又在`docker run`的时候指定了命令, 该怎么办?

这种情况下, 后者的内容将会覆盖前者. 以上面第2种情况为例.

```
docker run -d centos7:test tail -f /etc/hostname
```

这样启动的容器中, 运行的将是`tail -f /etc/hostname`, 而`tail -f /etc/os-release`将无效.

或者下面这种, 也是一样的.

```
docker run -it centos7:test /bin/bash
```

这会直接进入容器的bash终端, `CMD`中的`tail`不会执行.

## ENTRYPOINT

有了`CMD`, dockerfile 又提供了`ENTRYPOINT`.

`CMD`与`ENTRYPOINT`的区别在于: 

- 使用`CMD`创建的镜像, 如果在启动容器时手动指定了命令, 会覆盖`CMD`指定的命令;
- 而使用`ENTRYPOINT`创建的镜像, 我们加的命令会成为`ENTRYPOINT`后追加的**参数**(是参数哦, 不是两条命令都执行哦), 不会覆盖.

以如下dockerfile为例

```dockerfile
FROM centos:7
ENTRYPOINT ["tail", "-f", "/etc/os-release"]
```

> 注意: 本例中不使用`ENTRYPOINT tail -f" /etc/os-release`格式.

```
docker build -f dockerfile -t centos7:test .
docker run -d centos7:test tail -f /etc/hostname
```

进入容器查看进程得到如下输出

```console
[root@deabb9078762 ~]# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 07:16 ?        00:00:00 tail -f /etc/os-release tail -f /etc/hostname
root         7     0  0 07:16 ?        00:00:00 su - root
root         8     7  0 07:16 ?        00:00:00 -bash
root        19     8  0 07:16 ?        00:00:00 ps -ef
```

> 这种场景必须使用`-d`启动容器, 如果使用`docker run -it centos7:test /bin/bash`, `ENTRYPOINT`中的`tail`命令会把`/bin/bash`的内容也打印出来的...

## ENTRYPOINT+CMD 默认参数

`CMD`指令还有一种使用方法, 可以为`ENTRYPOINT`提供参数, 就和上面`ENTRYPOINT ["tail", "-f", "/etc/os-release"]`与`tail -f /etc/hostname`的合并一样.

```
CMD ["参数1", "参数2"]
```

`ENTRYPOINT`本身也可以包含参数, 但是由于`CMD`指定的值可以被命令行覆盖, 你可以把那些可能需要变动的参数写到`CMD`里而把那些不需要变动的参数写到`ENTRYPOINT`里面例如：

```dockerfile
FROM centos:7
ENTRYPOINT ["top", "-b"]
CMD ["-c"]
```

这样, 如果在`docker run`时末尾不指定任何额外参数, 容器启动时就会执行`top -b -c`, `-c`即为默认参数. 

```
docker run -d centos7:test 
```

如果末尾指定了参数, 这样`CMD`里的参数(这里是`-c`)就会被覆盖掉, 而`ENTRYPOINT`里的不被覆盖, 正好实现了**默认参数**的功能.

```
docker run -d centos7:test -p
```

> `ENTRYPOINT`与`CMD`同时使用时, 只能都用列表形式(没试过都用shell命令格式的, 但一个列表一个shell命令的肯定不行)
