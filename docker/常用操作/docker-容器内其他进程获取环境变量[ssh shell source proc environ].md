# docker-容器内其他进程获取环境变量

参考文章

1. [docker容器中的环境变量](https://www.cnblogs.com/xuxinkun/p/10531091.html)
2. [xargs with export is not working](https://stackoverflow.com/questions/44364059/xargs-with-export-is-not-working)
3. [/bin/sh source from stdin (from other program) not file](https://superuser.com/questions/272485/bin-sh-source-from-stdin-from-other-program-not-file)

容器的环境变量只有主进程(1号进程)可以获取到, 其他进程是不可以的. 另外, 如果容器内置了ssh服务, 那么通过ssh登录容器终端也是无法获取到环境变量的.

使用`docker exec {containerID} env`即可查看容器中生效的环境变量

```ini
$ docker exec 984 env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/java/default/bin
TERM=xterm
AUTHORIZED_KEYS=**None**
JAVA_HOME=/usr/java/default
HOME=/root
...
```

进入到容器中, 查看进程的环境变量, 可以通过`/proc`下进行查看

```
cat /proc/{pid}/environ
```

因此, 容器中的环境变量也可以通过在容器中查看1号进程的环境变量来获取, 可以通过执行`cat /proc/1/environ |tr '\0' '\n'`命令进行查看.

## 批量加载

虽然上述命令可以查看, 但是我们还是更希望能够自动加载(比如写到`/root/.bashrc`文件里), 这样在ssh登录时就不会手动执行了. 于是在上述命令的基础上, 尝试了如下改动.

## xargs+export

```bash
$ ## cat /proc/1/environ |tr '\0' '\n' | xargs export
$ cat /proc/1/environ |tr '\0' '\n' | xargs -i export {}
xargs: export: No such file or directory
```

按照参考文章2的说法, 是因为`export`是bash的内置命令, 所以不能这么用, 推荐使用`source`和`.`

## source

我首先想到的是, 将`cat /proc/1/environ |tr '\0' '\n'`结果写入文件后, 调用`source`加载, 然后再将该文件删除, 或者直接写入到`/root/.bashrc`就行.

但是我还是想要一个可以在行内执行的命令, 于是找到了参考文章3.

```bash
source <(cat /proc/1/environ |tr '\0' '\n')
```

> `<(xxx)`这个用法可以注意一下.

## 最终

```
source <(cat /proc/1/environ |tr '\0' '\n' | grep -v ' ')
```
