# Docker网络模式

参考文章

1. [docker的四种网络模式](http://www.cnblogs.com/frankielf0921/p/5822699.html)

1. host模式: docker使用的网络实际上和宿主机一样, 能看到宿主机所有的网卡信息.
    - `docker run --net=host`
2. container模式: 多个容器使用共同的网络, 看到的ip是一样的. 
    - kuber 中的 pause 容器就是用的这个原理
    - `docker run --net=container:容器id/容器名称`
3. bridge模式: 此模式会为每个容器分配一个独立的network namespace
    - `--net=网桥名称(默认是docker0)`
4. none 模式: 这种模式下, 不会配置任何网络. 
    - `--net=none`

使用`docker info`可以查看跨主机通信的网络模型实现, 以`docker-ce`为例:

```
$ docker info
Server:
...
 Plugins:
  ...
  Network: bridge host ipvlan macvlan null overlay
```

注意其中的`ipvlan`, `macvlan`, `overlay`.
