参考文章

1. [docker networking namespace not visible in ip netns list](https://stackoverflow.com/questions/31265993/docker-networking-namespace-not-visible-in-ip-netns-list)


查看`docker run`的选项发现, `uts/pid/userns`这3个命名空间可以直接在启动容器时指定, 但`net`和`mount`却没有明确指定.

我觉得这也是docker做隔离最下功夫的部分, `--net`可以指定`container:容器ID`使用某容器的网络, `--volumes-from 容器ID`可以指定某容器的卷.
