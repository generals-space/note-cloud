# Docker资源限制(三)磁盘限制

参考文章

1. [[经验分享] docker的资源隔离---cpu、内存、磁盘限制](https://www.iyunv.com/thread-116572-1-1.html)
2. [DOCKER基础技术：LINUX CGROUP](https://coolshell.cn/articles/17049.html)
    - 磁盘I/O限制

## 1. 空间限制

```
dd if=/dev/zero of=/tmp/ddfile bs=1M count=10240
```

## 2. IO限制

docker容器默认的空间是10G,如果想指定默认容器的大小（在启动容器的时候指定），可以在docker配置文件里通过dm.basesize参数指定，比如`docker -d --storage-opt dm.basesize=20G`指定默认的大小为20G.

