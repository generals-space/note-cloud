# mysql容器数据丢失

参考文章

1. [mysql的docker镜像中如何创建数据库](http://dockone.io/question/887)

问题描述:

使用mysql官方docker镜像, 新建一个容器A后在其中创建数据库并存储数据. 然后使用`docker commit`将容器A保存为一个镜像B. 以镜像B为基础启动容器C, 但是C中没有保存之前在A中创建的数据.

原因分析:

mysql-server的[Dockerfile](https://hub.docker.com/r/mysql/mysql-server/~/dockerfile/)有这样的一行

```
VOLUME /var/lib/mysql
```

意味着使用这个镜像时, 容器的`/var/lib/mysql`目录会被映射到宿主机上docker工作目录(默认为`/var/lib/docker`)下的某个目录. 具体映射到哪个目录, 可以通过 `docker inspect containerID` 查看.

```
$ docker inspect 原mysql容器ID | grep -i volume
"VolumeDriver": "",
 "VolumesFrom": null,
     "Source": "/var/lib/docker/volumes/023bb9aa4c35bd12625f89768fbfd86b73f0c8286fbc2d504e921872886b0e70/_data",
 "Volumes": {
$ cd /var/lib/docker/volumes/023bb9aa4c35bd12625f89768fbfd86b73f0c8286fbc2d504e921872886b0e70/
$ ls
_data
```

如果不做更改, 那么也就意味着你写入的数据会被直接写入到宿主机的该目录中, 而且不随容器的销毁而销毁。同样, commit的时候, 该目录的内容也不会被加入到镜像中。所以使用commit出来的镜像, 你就无法看到先前的数据了, 因为commit命令不会将挂载卷中的数据commit到镜像中.

解决方法:

将原容器的`_data`目录复制到目标容器的挂载卷下, 重启容器即可.

建议读一下[mysql-server](https://hub.docker.com/r/mysql/mysql-server/)中`Where to Store Data`一节，会对你了解如何存储mysql数据有所帮助.
