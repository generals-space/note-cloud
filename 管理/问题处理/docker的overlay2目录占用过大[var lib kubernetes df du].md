# docker的overlay2目录占用过大[var lib kubernetes df du]

参考文章

1. [Docker overlay2占用大量磁盘空间解决办法](https://blog.csdn.net/m0_67390788/article/details/123869198)

由于容器在运行中持续写入数据, /var/lib/docker 目录把磁盘占满了, 使用du排查的时候, 发现其下最大的目录是 overlay2.

但是进入到这个目录下, 这里面并不是容器id, 而是容器各layer分层的id. 使用du可以发现占用空间过大的分层, 但是如何反查出所属的容器id?

按照参考文章1中所说, 在使用docker inspect命令的结果中, 包含了该容器的分层信息, 可以使用如下命令反查

```
docker ps -q | xargs docker inspect --format '{{.State.Pid}}, {{.Name}}, {{.GraphDriver.Data.WorkDir}}' | grep "layer分层id"
```
