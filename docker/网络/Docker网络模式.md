# Docker网络模式

参考文章

1. [docker的四种网络模式](http://www.cnblogs.com/frankielf0921/p/5822699.html)

host模式 ：

docker run --net=host

docker使用的网络实际上和宿主机一样, 能看到宿主机所有的网卡信息.


container模式：

```
docker run --net=container:容器id/容器名称
```

多个容器使用共同的网络，看到的ip是一样的。
    
none 模式

--net=none

这种模式下，不会配置任何网络。

bridge模式

--net=网桥名称(默认是docker0)

此模式会为每个容器分配一个独立的network namespace