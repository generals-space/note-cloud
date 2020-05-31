# docker的netns路径

参考文章

1. [ip netns命令详解](https://blog.csdn.net/supahero/article/details/100606953)

使用`ip netns add ns01`创建的`netns`会出现在`/var/run/netns`目录下(默认这个目录不存在, 在第一次创建`netns`时自动创建).

```console
$ ip netns add ns01
$ cd /var/run/netns/
$ ls
ns01
```

但 docker 创建出来的容器的 netns 并不在`/var/run/netns`, 而是在`/proc/[pid]/ns/net`, 其中容器内部的进程pid可以通过如下命令查看.

```console
$ docker inspect ad7a7a5952fd | grep Pid
            "Pid": 1843,
            "PidMode": "",
            "PidsLimit": null,
```

然后可以使用如下命令建立一个软链接再进行访问.

```console
$ ip netns ls
$ ln -s /proc/1843/ns/net /var/run/netns/nginx
$ ip netns ls
nginx
$ ip netns exec nginx sh
sh-4.2# ls
```
