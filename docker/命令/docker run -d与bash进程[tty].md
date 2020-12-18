# docker run -d与bash进程

使用`docker run`时可以分别使用`-d`(服务模式)和`-it`(交互模式), 以及末尾追加的命令启动容器. 如下

```
docker run -d registry.cn-hangzhou.aliyuncs.com/generals-space/centos7 tail -f /etc/yum.conf
docker run -it registry.cn-hangzhou.aliyuncs.com/generals-space/centos7 /bin/bash
```

虽然看起来`/bin/bash`也是长驻进程, 但它不能运行在`-d`模式下. 

```
docker run -d registry.cn-hangzhou.aliyuncs.com/generals-space/centos7 /bin/bash
```

上面执行虽然不会报错, 但容器会立刻退出.

若要运行, 需得使用`-dt`(tty模式).

```
docker run -dt registry.cn-hangzhou.aliyuncs.com/generals-space/centos7 /bin/bash
```

这样容器就可以运行了, 之后可以使用`docker exec`进入该容器.

...其实可以用常规的`docker run -it 镜像名 /bin/bash`, 然后直接关掉终端, 这个容器就一直在后台运行了.
