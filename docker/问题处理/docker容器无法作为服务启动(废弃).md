# docker容器无法作为服务启动(废弃)

参考文章

1. [Docker为什么刚运行就退出了?](http://blog.simcu.com/archives/467)

### 1. 问题描述&原因分析

有些容器如果不使用`-it`选项并搭配执行`/bin/bash`命令无法以服务形式保持启动状态, 都是立刻结束, 使用`docker logs 容器ID`也没有错误日志, `docker ps -a`可以看到`Exited (0)`, 说明并未出错而是正常能出.

因为容器是否长久运行, 与`docker run`指定的命令有关. **使用-d选项使Docker容器后台运行, 就需要这个指定命令必须是一个前台进程**.

这个是docker的机制问题, 比如普通的web容器, 以nginx和fpm为例. 正常情况下, 我们配置启动服务只需要启动响应的service即可, 例如

```
$ service nginx start && service php5-fpm start
```

这样做, nginx和fpm均为daemon模式运行, 如果在`docker run`中指定这样的命令, 就会导致docker前台没有运行的应用. 这样的容器, 后台启动后, 会立即自杀, 因为它觉得没事可做了.

### 2. 解决方法

#### 2.1 

最佳的解决方案是, 将你要运行的程序以前台进程的形式运行(如果可以的话). 如果你的容器需要同时启动多个进程, 那么也只需要, 或者说只能将其中一个挂起到前台即可. 比如上面所说的web容器,我们只需要将启动指令修改为:

```
service php5-fpm start && nginx -g "daemon off;"
```

这样, fpm会在容器中以后台进程的方式运行, 而nginx则挂起进程至前台运行. 这样, 就可以保持容器不会认为没事可做而退出, 并且容器本身会因`-d`选项的存在以服务模式运行.

#### 2.2 

对于有一些你可能不知道怎么保持前台运行的程序, 提供一个投机方案: 在启动的命令之后, 添加类似于`tail`, `top`这种可以前台运行的程序, 这里特别推荐`tail`, 然后持续输出你的log文件.

还是以上文的web容器为例, 还可以写成如下 

```
service nginx start && service php5-fpm start && tail -f /var/log/nginx/error.log
```

#### 2.3

还有一个比较蠢的方法: 使用`-it`选项并执行`/bin/bash`命令, 进入容器shell. 然后退出, 此时容器将会停止. 使用`docker ps -a`查看刚才运行的容器ID, 再使用`docker start 容器ID`将会使其进入服务状态, `docker ps`可以看到它依然在运行, 而且命令还是`/bin/bash`. 然后就可以通过`nsenter`等工具进入容器了.

不过, 这对需要在命令行执行启动服务的命令的情况不适用, 因为`docker start`这个容器后, 服务还是默认停止的状态. 对Dockerfile文件中存在有`CWD`命令的镜像也不会起作用...当做日常开发的小伎俩玩玩吧.
