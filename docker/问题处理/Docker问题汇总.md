# Docker问题汇总

## 1. Docker容器bash中输入中文

**参考文章**

[linux终端不能输入中文解决方法 ](http://blog.sina.com.cn/s/blog_5c4dd3330100cpmm.html)

[在Docker容器bash中输入中文](http://blog.shiqichan.com/Input-Chinese-character-in-docker-bash/)

------

docker容器内的bash无论如何都无法输入中文, 不管是在启动容器时打开bash, 还是以服务形式启动容器后再通过 `nsenter`工具进入容器之后显示的bash. 不管是在什么情况下输入甚至粘贴, 不是出现乱码, 回车无反应甚至根本无法上屏. 而且输出时中文全都是乱码.

尝试在容器内 `/root`家目录新建 `.inputrc`文件, 添加以下内容

```shell
set meta-flag on
set convert-meta off
set input-meta on
set output-meta on
```

重启容器发现可以在bash命令行上输入中文, 但是回车发现与预期结果不同, 而且输出时中文依然是乱码. 尝试设置`locale`, 不管将环境变量LANG设置为 `LANG=en_US.UTF-8`还是 `LANG=zh_CN.UTF-8`都不起作用.

------

真正的解决方法是, 在启动容器时传入 `env`参数

```shell
docker run -i -t ubuntu env LANG=C.UTF-8 /bin/bash
```

或是在Dockerfile文件中写入如下行

```shell
ENV LANG=C.UTF-8
```

## 2. CentOS:7的docker镜像使用systemctl

### 问题描述

使用`docker.io/centos:7`的docker官方镜像, 配置`php-fpm`的服务脚本完成. 但是使用`systemctl`命令操作`php-fpm`时出现如下问题. (服务脚本的建立参考[这里](http://www.centoscn.com/CentOS/config/2015/0507/5374.html))

```shell
[root@e13c3d3802d0 /]# systemctl start php-fpm
Failed to get D-Bus connection: Operation not permitted
```

### 原因分析及解决方法

参考下面两篇文章的解释.

[docker Failed to get D-Bus connection 报错](http://welcomeweb.blog.51cto.com/10487763/1735251)

[如何在Docker CentOS容器中使用Systemd](http://www.codesec.net/view/434721.html)

这是CentOS:7镜像中的一个bug, 目前无法修复, 只能等到7.2版镜像. 这个bug的原因是因为`dbus-daemon`没能启动. 其实`systemctl`并不是不可以使用. 将你dockerfile的`CMD`或者`ENTRYPOINT`设置为`/usr/sbin/init`即可, 容器会在运行时将`dbus`等服务启动起来. 然后再执行`systemctl`命令即可运行正常.

```shell
docker run --privileged  -e "container=docker"  -v /sys/fs/cgroup:/sys/fs/cgroup -d docker.io/centos:7  /usr/sbin/init
```

> 注意: 上面的命令要完全执行...少一个都会被坑的很惨.

`--privileged`: systemd 依赖于`CAP_SYS_ADMIN capability`. 意味着运行Docker容器需要获得 `privileged`(这不利于一个base image);

`-v`选项: systemd 依赖于访问`cgroups filesystem`; systemd 有很多并不重要的文件存放在一个docker容器中, 如果不删除它们会产生一些错误; 原来以为`-v`选项只是一个共享目录而已, 然而我擅自去掉这个选项后, `systemctl start php-fpm`执行时卡住了, 而且并没有执行成功. 加上这个选项才可以.

另外, 最好加上`-d`选项, 这个是我自己加上的. 交互时启动时竟然提示输入容器密码...而且尝试了很多遍都不对(也不是宿主机的密码)...这真是个奇妙的经历.

```shell
...
CentOS Linux 7 (Core)
Kernel 4.5.5-300.fc24.x86_64 on an x86_64

325606d855b9 login: root
Password:
Login incorrect

325606d855b9 login:
Password:
Login incorrect
...
```

## 3. docker容器无法作为服务启动

参考文章

[Docker为什么刚运行就退出了?](http://blog.simcu.com/archives/467)

### 3.1 问题描述&原因分析

有些容器如果不使用`-it`选项并搭配执行`/bin/bash`命令无法以服务形式保持启动状态, 都是立刻结束, 使用`docker logs 容器ID`也没有错误日志, `docker ps -a`可以看到`Exited (0)`, 说明并未出错而是正常能出.

因为容器是否长久运行, 与`docker run`指定的命令有关. **使用-d选项使Docker容器后台运行, 就需要这个指定命令必须是一个前台进程**.

这个是docker的机制问题, 比如普通的web容器, 以nginx和fpm为例. 正常情况下, 我们配置启动服务只需要启动响应的service即可, 例如

```
$ service nginx start && service php5-fpm start
```

这样做, nginx和fpm均为daemon模式运行, 如果在`docker run`中指定这样的命令, 就会导致docker前台没有运行的应用. 这样的容器, 后台启动后, 会立即自杀, 因为它觉得没事可做了.

### 3.2 解决方法

#### 3.2.1 

最佳的解决方案是, 将你要运行的程序以前台进程的形式运行(如果可以的话). 如果你的容器需要同时启动多个进程, 那么也只需要, 或者说只能将其中一个挂起到前台即可. 比如上面所说的web容器,我们只需要将启动指令修改为:

```
service php5-fpm start && nginx -g "daemon off;"
```

这样, fpm会在容器中以后台进程的方式运行, 而nginx则挂起进程至前台运行. 这样, 就可以保持容器不会认为没事可做而退出, 并且容器本身会因`-d`选项的存在以服务模式运行.

#### 3.2.2 

对于有一些你可能不知道怎么保持前台运行的程序, 提供一个投机方案: 在启动的命令之后, 添加类似于`tail`, `top`这种可以前台运行的程序, 这里特别推荐`tail`, 然后持续输出你的log文件.

还是以上文的web容器为例, 还可以写成如下 

```
service nginx start && service php5-fpm start && tail -f /var/log/nginx/error.log
```

#### 3.2.3

还有一个比较蠢的方法: 使用`-it`选项并执行`/bin/bash`命令, 进入容器shell. 然后退出, 此时容器将会停止. 使用`docker ps -a`查看刚才运行的容器ID, 再使用`docker start 容器ID`将会使其进入服务状态, `docker ps`可以看到它依然在运行, 而且命令还是`/bin/bash`. 然后就可以通过`nsenter`等工具进入容器了.

不过, 这对需要在命令行执行启动服务的命令的情况不适用, 因为`docker start`这个容器后, 服务还是默认停止的状态. 对Dockerfile文件中存在有`CWD`命令的镜像也不会起作用...当做日常开发的小伎俩玩玩吧.

## 4. Ubuntu14.04安装nsenter

ubuntu14的`util-linux`版本为2.20, 但想要进入docker容器, 不能低于`2.24`. 需要手动编译安装. 安装命令如下, 注意要首先安装依赖包.

```shell
sudo apt-get install autopoint autoconf libtool automake
wget https://www.kernel.org/pub/linux/utils/util-linux/v2.24/util-linux-2.24.tar.gz
tar xzvf util-linux-2.24.tar.gz
cd util-linux-2.24
./configure --without-ncurses
make && make install
```

## 4. mysql容器数据丢失

参考文章

[mysql的docker镜像中如何创建数据库](http://dockone.io/question/887)

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

## 5. 

```
docker pull daocloud.io/centos:6
6: Pulling from centos
32c4f4fef1c6: Extracting [==================================================>] 68.74 MB/68.74 MB
failed to register layer: ApplyLayer exit status 1 stdout:  stderr: symlink gawk /bin/awk: operation not supported
```

情境描述: CentOS7的虚拟机, 编译安装的docker, 版本为`1.14.rc2`, containerd与dockerd服务成功启动, `docker search`也可以正常使用, 但pull操作时出现上述错误.

原因分析: 当时dockerd启动的`graph`参数设置为了从windows宿主机共享的目录, 可能是由于docker需要linux的联合文件系统的特性支持而共享目录本质还是NTFS才会报错.

解决办法: 尝试将`graph`参数指向了linux本地任一目录, 重启dockerd, 再次pull时正常.

## 6. docker服务无法启动

参考文章

1. [Error starting daemon](https://github.com/moby/moby/issues/34680)

docker版本

```
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

```
Error starting daemon: couldn't create plugin manager: error setting plugin manager root to private: invalid argument
```

参考文章1中众说纷纭, 看起来也有点道理, 什么指定`--make-private`参数重新挂载硬盘, 没试. 有一位`Faqa`的方法十分非主流, 但我很有眼光的觉得他说的是对的...把docker路径由`/disk1/docker`转到`/disk1/docker_root/docker`...即, docker的路径不能在硬盘的第一级子目录...

OK, 果然解决了.
