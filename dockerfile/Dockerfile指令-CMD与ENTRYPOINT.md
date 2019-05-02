# Dockerfile指令-CMD与ENTRYPOINT

参考文章

`CMD`与`ENTRYPOINT`两者有相似之处, 都是镜像在实例化成容器时执行的命令.

首先要明白, `-d`(服务模式)和`-it`(交互模式)两者后面都可以指定要执行的命令.

以服务模式启动容器, 容器中运行`tail -f`命令.

```
$ docker run -d  centos:6 tail -f /etc/yum.conf
```

> 虽然看起来`/bin/bash`也是长驻进程, 但它不能运行在`-d`模式下. 若要运行, 需得使用`-dt`. tty模式.

------

进入正题, 如果一个镜像在封装之时没有通过`CMD`或是`ENTRYPOINT`指令指定在启动时执行的命令, 就只能在命令行启动时手动指定.

即, `CMD`与`ENTRYPOINT`的目的就是可以在容器启动时不必手动执行这一过程.

这其实也要看镜像的, 像centos, ubuntu这种官方系统镜像是不会为你指定启动命令的, 我们得自己定义. 而其他与nginx, mysql这种第三方应用级镜像, 则一般会带有启动命令, 目的是简化操作, 直接使用.

`CMD`与`ENTRYPOINT`的区别在于: 

- 使用`CMD`创建的镜像, 如果在启动容器时手动指定了命令, 会覆盖`CMD`指定的命令;

- 而使用`ENTRYPOINT`创建的镜像, 我们加的命令会成为`ENTRYPOINT`后追加的**参数**(是参数哦, 不是两条命令都执行哦), 不会发生覆盖.

## 指令格式

```
CMD ping www.baidu.com 
CMD ["/bin/ping", "www.baidu.com"]
```

同样, `ENTRYPOINT`也有2种可用的格式.

```
ENTRYPOINT ping www.baidu.com 
ENTRYPOINT ["/bin/ping","www.baidu.com"]
```

> md, 列表里必须用双引号...

> 注意1: `ENTRYPOINT ping www.baidu.com`这种格式依然会导致被命令行指定的命令覆盖, 第2种列表格式则不会.

> 注意2: 做这种命令覆盖的实验时, 不建议使用`echo`这种命令, 因为容器内外标准输出可能被阻断, 尽量写入到文件.

> 注意3: shell命令格式的CMD与ENTRYPOINT, 可以在其中通过`$变量名`引用`ENV`或是`-e`选项指定的环境变量, 但是数组形式不可以, 它会把`$变量名`当成字符串处理.

## `ENTRYPOINT`与参数追加

其实`ENTRYPOINT`这种所谓的不会覆盖的命令很鸡肋. 以如下dockerfile为例

```
FROM centos:6
ENTRYPOINT ["tail", "-f", "/etc/passwd"]
```

容器在启动时执行`tail -f /etc/passwd`, 这种情况下在`docker run`时就最好使用`-d`选项了. 并且我在启动命令末尾又通过指定`tail -f /etc/yum.conf`命令. 结果等容器启动后, 进入容器查看进程发现了如下情况

```
[root@deabb9078762 ~]# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 07:16 ?        00:00:00 tail -f /etc/passwd tail -f /etc/yum.conf
root         7     0  0 07:16 ?        00:00:00 su - root
root         8     7  0 07:16 ?        00:00:00 -bash
root        19     8  0 07:16 ?        00:00:00 ps -ef
```

这很让人无语, 因为这相当于在命令行指定的命令直接附加到`ENTRYPOINT`所指定的命令后面, 完全没有预想的那么神(两条命令都执行).

## CMD与ENTRYPOINT配合使用

`CMD`指令还有一种使用方法, 可以为`ENTRYPOINT`提供参数, 就和上面`ENTRYPOINT ["tail", "-f", "/etc/passwd"]`与`tail -f /etc/yum.conf`的合并一样.

```
CMD ["参数1", "参数2"]
```

`ENTRYPOINT`本身也可以包含参数, 但是由于`CMD`指定的值可以被命令行覆盖, 你可以把那些可能需要变动的参数写到`CMD`里而把那些不需要变动的参数写到`ENTRYPOINT`里面例如：

```
FROM centos:6
ENTRYPOINT ["top", "-b"]
CMD ["-c"]  
```

这样, 如果在`docker run`时不指定任何额外参数, 容器启动时就会执行`top -b -c`, `-c`即为默认参数, 如果指定自定义的参数, 这样`CMD`里的参数(这里是`-c`)就会被覆盖掉而`ENTRYPOINT`里的不被覆盖, 正好实现了**默认参数**的功能.

> `ENTRYPOINT`与`CMD`同时使用时, 貌似只能都用列表形式(没试过都用shell命令格式的, 但一个列表一个shell命令的肯定不行)
