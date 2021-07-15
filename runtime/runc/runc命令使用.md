# runc使用

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

### 1. 不修改 config.json 直接 run

```
cd /mycontainer
runc run mycontainerid
```

> `mycontainerid`可以是任意字符串.

如果未修改上面生成的`config.json`文件, `run`子命令会打开一个`sh`会话.

新开一个终端, 可以使用`runc list`子命令查看正在运行的容器.

```console
$ runc list
ID              PID         STATUS      BUNDLE             CREATED                          OWNER
mycontainerid   89791       running     /mycontainer       2020-04-15T09:46:03.463907749Z   root
```

> 使用`runc`启动的 container, 是没有办法使用`docker ps`查看到的.

当退出`sh`会话时, 容器也会结束, 而且貌似没有办法再找回了. 不像docker, stop后还可以重新start.

### 2. 修改 config.json 先 create 再 start

默认的`config.json`会创建 sh 会话, 这样只能用`run`子命令启动.

```
$ runc create mycontainerid
ERRO[0000] cannot allocate tty if runc will detach without setting console socket
cannot allocate tty if runc will detach without setting console socket
```

接下来我们修改`config.json`中的配置, 使之可以像`docker run -d`一样实现后台启动.

```json
{
	"ociVersion": "1.0.1-dev",
	"process": {
		"terminal": false,
		"args": [
			"tail", "-f", "/etc/hosts"
		]
	}
}
```

```bash
cd /mycontainer
runc create mycontainerid
```

执行`create`后就可以使用`list`子命令查看, 此时容器处于`created`状态.

```console
$ runc list
ID              PID         STATUS      BUNDLE             CREATED                          OWNER
mycontainerid   90414       created     /mycontainer   2020-04-15T10:28:44.760800845Z   root
```

使用`ps`可以查看目标pid.

```console
$ ps -ef | grep 90414
root      90414      1  0 18:28 ?        00:00:00 runc init
```

看来`runc init`是一个类似守护进程的命令.

然后启动容器.

```
runc start mycontainerid
```

此时使用`list`命令再看, 容器已处于`running`状态了.

```console
$ runc list
ID              PID         STATUS      BUNDLE             CREATED                          OWNER
mycontainerid   90414       running     /mycontainer   2020-04-15T10:28:44.760800845Z   root
```

最主要是的, `ps`的目标进程发生了变化, 我们看不到上面的`runc init`进程了.

```console
$ ps -ef | grep 90414
root      90414      1  0 18:28 ?        00:00:00 tail -f /etc/hosts
```

> 看来是由`init`进程`fork`了一个子进程, 然后`exec`了目标命令啊.

可以使用如下命令进入到容器终端, 与`docker exec`同理.

```console
$ runc exec -t mycontainerid sh
/ # ls
bin   dev   etc   home  proc  root  sys   tmp   usr   var
```

使用`runc kill`可以实现`docker stop`的功能, 使 container 变成`stopped`的状态.

```console
$ runc kill mycontainerid KILL
$ runc list
ID              PID         STATUS      BUNDLE             CREATED                          OWNER
mycontainerid   0           stopped     /tmp/mycontainer   2020-04-15T10:28:44.760800845Z   root
```

第2种方式可以实现更多的设置, 比如在`create`后, `start`前, 可以设置 container 的网络.

