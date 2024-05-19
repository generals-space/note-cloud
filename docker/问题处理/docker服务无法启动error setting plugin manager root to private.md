# docker服务无法启动error setting plugin manager root to private

参考文章

1. [Error starting daemon](https://github.com/moby/moby/issues/34680)

docker版本

```log
$ docker version
Client:
 Version:      17.09.1-ce
 API version:  1.32
 Go version:   go1.8.3
 Git commit:   19e2cf6
 Built:        Thu Dec  7 22:23:40 2017
 OS/Arch:      linux/amd64

Server:
 Version:      17.09.1-ce
 API version:  1.32 (minimum version 1.12)
 Go version:   go1.8.3
 Git commit:   19e2cf6
 Built:        Thu Dec  7 22:25:03 2017
 OS/Arch:      linux/amd64
 Experimental: false
```

场景描述:

docker的存储路径在一块独立硬盘上, 不过fstab文件貌似有点问题, 本来docker服务运行的好好的, 重启了下系统, 那块硬盘没能自动挂载, 由于docker是开机启动, 所以`docker images`镜像都不见了. 

关闭docker, 手动挂载硬盘, 启动docker...失败了. 查看日志, 报了如下错误.

```log
Error starting daemon: couldn't create plugin manager: error setting plugin manager root to private: invalid argument
```

参考文章1中众说纷纭, 看起来也有点道理, 什么指定`--make-private`参数重新挂载硬盘, 没试. 有一位`Faqa`的方法十分非主流, 但我很有眼光的觉得他说的是对的...把docker路径由`/disk1/docker`转到`/disk1/docker_root/docker`...即, docker的路径不能在硬盘的第一级子目录...

OK, 果然解决了.
