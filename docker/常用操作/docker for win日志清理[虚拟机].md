# docker for win日志清理

参考文章

1. [How to Delete Docker Container Log Files (Windows or Linux) ](https://blog.jongallant.com/2017/11/delete-docker-container-log-files/)
    - 介绍了win下MobyLinuxVM虚拟机的存在
2. [How to SSH into the Docker VM (MobyLinuxVM) on Windows ](https://blog.jongallant.com/2017/11/ssh-into-docker-vm-windows/)

默认情况下docker内部进程会将日志打印到标准输出, 这种日志其他最终写入了宿主机的某个文件. 使用`inspect`子命令可以查看.

```
$ docker inspect service_parserhub-serv_1 | grep -i logpath
        "LogPath": "/var/lib/docker/containers/7a394217a0ce7b8d4cd3c452ab683641be8d581bcc5dc3a21bf668644e81d64a/7a394217a0ce7b8d4cd3c452ab683641be8d581bcc5dc3a21bf668644e81d64a-json.log",
```

当我们使用`docker logs`打印了太多日志而影响了我们对错误的排查时, 可以考虑清空这个文件.

但是在windows下, docker服务本质是运行在一个虚拟机中, 名为`MobyLinuxVM`. 镜像, 容器, 挂载卷等都存放在独立的VM文件中, 我们没有办法直接进入这个虚拟机进行操作.

按照参考文章2中提到的方法, 执行如下命令

```
docker run --privileged -it -v /var/run/docker.sock:/var/run/docker.sock jongallant/ubuntu-docker-client 
## 下面这句是在容器内执行的, 目的是清理linux中的日志
docker run --net=host --ipc=host --uts=host --pid=host -it --security-opt=seccomp=unconfined --privileged --rm -v /:/host alpine /bin/sh
chroot /host
```

上面的操作是以能得到和ssh进入容器相同的环境而做了很多步骤, 如果我们仅清空指定容器的日志, 则没必要写这么复杂(事实上, 我并不希望使用别人构建的镜像`jongallant/ubuntu-docker-client`, 而且似乎也没有必要). 

我们只需要在启动容器时将根目录挂载到容器内部, 就可以查看到VM的文件系统了.

```
docker run -it --name vmhost -v /:/host generals/alpine /bin/sh
/ # cd /host
/host # ls
C           bin         d           etc         host_mnt    media       opt         proc        run         sendtohost  sys         usr
D           c           dev         home        lib         mnt         port        root        sbin        srv         tmp         var
/host # cd var/lib/docker/containers/
/host/var/lib/docker/containers # :> 7a394217a0ce7b8d4cd3c452ab683641be8d581bcc5dc3a21bf668644e81d64a/7a394217a0ce7b8d4cd3c452ab683641be8d581bcc5dc3a21bf668644e81d64a-json.log
/host/var/lib/docker/containers #
```

然后再使用`docker logs`就会发现日志已经被清空了.

如果真的想进入`MobyLinuxVM`, 只要执行上面命令的第2条就可以了.

这样执行ps时就可以看到所有容器正在的进程.

```console
$ docker run -it --name vmhost --net=host --ipc=host --uts=host --pid=host --security-opt=seccomp=unconfined --privileged -v /:/host generals/alpine /bin/sh
/ # ps -ef
PID   USER     TIME  COMMAND
    1 root      0:02 /sbin/init text
    2 root      0:00 [kthreadd]
    3 root      0:04 [ksoftirqd/0]
    5 root      0:00 [kworker/0:0H]
    7 root      0:04 [rcu_sched]
    8 root      0:00 [rcu_bh]
```
