# docker-容器内其他进程获取环境变量

参考文章

1. [docker容器中的环境变量](https://www.cnblogs.com/xuxinkun/p/10531091.html)

使用`docker exec {containerID} env`即可查看容器中生效的环境变量

```console
$ docker exec 984 env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/java/default/bin
TERM=xterm
AUTHORIZED_KEYS=**None**
JAVA_HOME=/usr/java/default
HOME=/root
...
```

进入到容器中，查看进程的环境变量，可以通过/proc下进行查看

```
cat /proc/{pid}/environ
```

因此，容器中的环境变量也可以通过在容器中查看1号进程的环境变量来获取。可以通过执行`cat /proc/1/environ |tr '\0' '\n'`命令进行查看.
