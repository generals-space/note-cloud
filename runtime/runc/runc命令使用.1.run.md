# runc命令使用.1.run

runc 版本: 3e425f80a8c931f88e6d94a8c831b9d5aa481657

```bash
# create the top most bundle directory
mkdir /mycontainer
cd /mycontainer

# create the rootfs directory
mkdir rootfs

# export busybox via Docker into the rootfs directory
## docker create 创建一个未启动的容器
## docker export 可以导出容器的整个文件系统
## 与 docker save 不同, ta的目标是 container 而非 image, 
## 导出的tar包中也不包括分层信息.
## tar -C 可以将tar包解压到目标目录
docker export $(docker create busybox) | tar -C rootfs -xvf -
```

## 启动方式

先生成`config.json`

```
cd /mycontainer
runc spec
```

## 不修改 config.json 直接 run

```log
$ cd /mycontainer
$ runc run c01
INFO[0000] setupIO: detach: false, sockpath:
INFO[0000] wait for communicate with runc init golang part
INFO[0000] child process in init()
INFO[0000] os getwd: /home/project/mycontainer/rootfs
INFO[0000] tty.forward pid1: 10208, notifySocket: <nil>
sh-5.0#
```

> `c01`可以是任意字符串.

如果未修改上面生成的`config.json`文件, `run`子命令会打开一个`sh`会话.

新开一个终端, 可以使用`runc list`子命令查看正在运行的容器.

```console
$ runc list
ID     PID    STATUS     BUNDLE          CREATED                           OWNER
c01    897    running    /mycontainer    2020-04-15T09:46:03.463907749Z    root
```

> 使用`runc`启动的 container, 是没有办法使用`docker ps`查看到的.

当退出`sh`会话时, 容器也会结束, 而且貌似没有办法再找回了. 不像docker, stop后还可以重新start.
