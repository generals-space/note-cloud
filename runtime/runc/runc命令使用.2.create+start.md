# runc命令使用.2.create+start

默认的`config.json`会创建 sh 会话, 这样只能用`run`子命令启动.

```
$ runc create c02
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
runc create c02
```

执行`create`后就可以使用`list`子命令查看, 此时容器处于`created`状态.

```log
$ runc list
ID     PID    STATUS      BUNDLE          CREATED                           OWNER
c02    904    created     /mycontainer    2020-04-15T10:28:44.760800845Z    root
```

使用`ps`可以查看目标pid.

```log
$ ps -ef | grep 904
root      904      1  0 18:28 ?        00:00:00 runc init
```

看来`runc init`是一个类似守护进程的命令.

然后启动容器.

```
runc start c02
```

此时使用`list`命令再看, 容器已处于`running`状态了.

```log
$ runc list
ID     PID    STATUS     BUNDLE          CREATED                           OWNER
c02    904    running    /mycontainer    2020-04-15T10:28:44.760800845Z    root
```

最主要是的, `ps`的目标进程发生了变化, 我们看不到上面的`runc init`进程了.

```log
$ ps -ef | grep 904
root      904      1  0 18:28 ?        00:00:00 tail -f /etc/hosts
```

> 看来是由`init`进程`fork`了一个子进程, 然后`exec`了目标命令啊.

可以使用如下命令进入到容器终端, 与`docker exec`同理.

```log
$ runc exec -t c02 sh
/ # ls
bin   dev   etc   home  proc  root  sys   tmp   usr   var
```

使用`runc kill`可以实现`docker stop`的功能, 使 container 变成`stopped`的状态.

```log
$ runc kill c02 KILL
$ runc list
ID     PID    STATUS     BUNDLE              CREATED                           OWNER
c02    0      stopped    /tmp/mycontainer    2020-04-15T10:28:44.760800845Z    root
```

第2种方式可以实现更多的设置, 比如在`create`后, `start`前, 可以设置 container 的网络.
