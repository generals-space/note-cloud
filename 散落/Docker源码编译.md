# Docker源码编译

参考文章

[如何编译docker源码](http://mogu.io/compile-docker-121)

[Work with a development container](https://docs.docker.com/v1.11/opensource/project/set-up-dev-env/)

[Docker 1.11 增强功能：直接在runC和containerd构建引擎](http://dockone.io/article/1327)

> 本文编译的是github/docker的master分支, 完成后的版本显示为`1.14.0-dev`

## 1. 对源码编译的认识

Docker源码编译需要预先在目标主机上运行docker服务, 这不免有着`先有蛋还是先有鸡`的问题. 但是官方文档的目的是, 随github中发布的编译脚本(包括`Makefile`, `Dockerfile`还有一系列shell脚本等), 都是为了方便开发者加入docker开源社区, 为其贡献代码的.

所以, 源码编译docker的环境要求是, 运行有docker服务, 安装了`make`工具.

`git`工具其实可有可无, 只要你能得到一份`git clone`的docker仓库代码就行.

> 注意: **不能是直接在网页上下载的zip包, 因为zip包中不包含`.git`对象库, 而编译过程中需要知道当前的编译版本, 所以只能通过`git clone`得到源码包**.

另外, 完全从零编译docker也并非不可能. 但是你首先需要了解自己系统上是否安装有哪种适合docker运行的驱动, 比如用于存储的文件系统选择, 是`aufs`或是`devicemapper`等. 并且需要弄清楚docker源码中给出的`Makefile`, `Dockerfile`文件中的编译步骤...too difficult to implement.

------

正式编译之前的操作其实是利用docker源码中附出的`Dockerfile`文件构建一个镜像, 其中将会安装一系列编译运行docker源码的软件包, 也是方便开发者不必搭建步骤复杂的环境.

然后, 以这个镜像启动一个docker容器, 完成docker的可执行文件的编译. 开发者也是在类似这样的容器中验证他们对代码的修改结果.

------

编译完成后, 会得到`client`与`daemon`两个目录.

从1.11版本开始, docker将集成的docker服务拆分成了一个个组件, 以符合`OCI`标准, 并且增强对集群的支持.

client下的`docker`文件就是平时我们执行`docker info|run|start`等操作的命令行工具; 

`daemon`中会有比较多的可执行文件, 包括最顶层的, 用于管理镜像的服务`dockerd`, 用于启动, 销毁容器的`containerd`, 用于运行容器的微型服务`runc`, 还有一个运行于`containerd`与`runc`之间, 用于容器间通信的服务`containerd-shim`.

## 2. 编译步骤分析及建议

按照官方文档的说明, 只需要在docker源码目录下执行`make build && make binary`就可以得到二进制的可执行文件, 并且`The first build may take a few minutes to create an image`. 编译了整整两天的我表示...呵呵.

天网恢恢啊...

------

`make build`命令所做的是, 使用`Dockerfile`构建一个docker镜像, 下载各种依赖软件的步骤就写在这里面.

为了深刻地理解其编译原理, 建议一边不修改任何东西直接构建, 一边尝试修改它的`Dockerfile`. 因为执行时间超长的步骤一定是需要fanqiang才能迅速搞定的.

这里是一些建议.

1. Dockerfile利用`debian`的镜像为基础, 可以事先利用国内的docker镜像点下载好debian基础镜像, 比如`daocloud.io/debian:jessie`.
2. 以pull下的debian镜像创建容器, 在其中手动安装pip(因为debian镜像自带的源中, pip可能无法使用proxy与trust功能, 之后一系列使用pip安装的依赖将无法使用镜像加速...), 并且设置好pipy镜像源, 可以换成豆瓣的镜像源. 完成后commit成新的镜像, 修改Dockerfile中的`FROM`指令为此镜像.
3. 修改apt镜像源地址, 可以换成阿里云或网易的镜像源, 方便之后执行`apt-get`的相关操作.
4. 创建镜像中会在容器内使用git下载依赖的源码包, 速度太慢, 需要为git设置代理. 可以在安装完`git`工具后添加`RUN`指令, 手动指定代理地址.
5. 容器内同样会使用curl下载一些东西, 所以curl也要设置代理.

除了Dockerfile, 在实际执行中还发现, `contrib/download-frozen-image-v2.sh`中也用到了curl下载, 所以这里也需要修改.

随着docker源码的更新, 一些具体的方法可能不再适用, 但基本原理就是这样 - 镜像与代理.

------

`make build`后执行`make binary`, 生成的二进制文件位置在`bundles`目录, 会生成以当前docker源码所在的版本号为名称的目录, 本文试验中为`1.14.0-dev`.

## 3. 启动

停止正在运行的docker服务, 然后将二进制目录包含的`binary-client`与`binary-daemon`这两个目录下的可执行文件全部copy到`/usr/bin`目录下.

> 注意: 虽然不一定要copy到`/usr/bin`下, 但要保证这些可执行文件在`$PATH`的路径中. 也可以使用软链接等方式.

启动命令, 以下是非`daemon`方式, 可能需要多终端. 另外, `dockerd`服务依赖`containerd`, 所以后者要先启动.

注意: 由于相当于版本升级, 原来docker的存储目录`/var/lib/docker`不能作为新docker服务的存储目录, 可能会出现不兼容的情况.

```
$ docker-containerd --runtime /usr/bin/docker-runc --shim /usr/bin/docker-containerd-shim
$ dockerd -D -g /var/lib/docker_new --containerd /run/containerd/containerd.sock
```

然后使用`docker`命令`pull`新的镜像, `run`启动容器.

```
$ docker pull daocloud.io/nginx
$ docker run -it -name nginx daocloud.io/nginx /bin/bash
root@43f6a9550b88:/# 
```

## FAQ

### 1.

```
$ ./dockerd -D
libcontainerd: containerd health check returned error: rpc error: code = 14 desc = grpc: the connection is unavailable
```

启动dockerd之前需要先启动`containerd`服务, 并且启动dockerd时要指定`containerd`的`.sock`文件路径.

### 2.

```
$ ./dockerd -D 
ERRO[0001] devmapper: Udev sync is not supported. This will lead to data loss and unexpected behavior. Install a dynamic binary to use devicemapper or select a different storage driver.      For more information, see https://docs.docker.com/engine/reference/commandline/daemon/#daemon-storage-driver-option 
ERRO[0001] [graphdriver] prior storage driver "devicemapper" failed: driver not supported 
FATA[0001] Error starting daemon: error initializing graphdriver: driver not supported 
```

网上有很多说是`Udev`什么什么的, 但是docker源码编译的前提是已经有运行docker服务, 并且会延用当前docker的驱动与配置, 所以不太可能是这些原因.

真正原因很可能是在启动`dockerd`的时候没有额外指定`-g`参数, 或是指定了原`docker`的`-g`选项的路径. 新版本的docker可能在目录结构上与之前不同, 所以这两个路径不能指向同一个.

``` 
$ ls
containers image network overlay swarm tmp trust volumes
```

这样看来, docker版本升级可能会导致下载到本地的镜像无法重用. 可以预先将这些镜像push到私有镜像库中, 升级后再pull下载.

### 3. 

```
$ ./docker run -it --name nginx daocloud.io/nginx /bin/bash
./docker: Error response from daemon: containerd-shim not installed on system.
```

问题描述: 从daocloud下载了镜像, `docker run`时报错. 启动`docker-containerd`时明明已经指定了`docker-containerd-shim`与`docker-runc`的绝对路径, 但还是报这个错误.

```
./docker-containerd --runtime /root/Downloads/docker/bundles/latest/binary-daemon/docker-runc --shim /root/Downloads/docker/bundles/latest/binary-daemon/docker-containerd-shim
```

解决办法: 

`docker-containerd-shim`与`docker-runc`需要放在`$PATH`变量中的路径下, 最好把编译完成的文件都放在`/usr/bin`目录下, 或是软链接过去.

待解决
------

lvm与device_mapper关系.

[Linux多路径、LVM的基础--内核Device Mapper机制](http://blog.csdn.net/smstong/article/details/40583129)


## iptables设置

初始时清空`iptables`规则, 启动`dockerd`服务, 将会写入如下规则(不会写入到`/etc/sysconfig/iptables`文件中, 但单纯关闭docker服务不会清空docker相关的规则).

```
#### filter表
$ iptables -L
Chain INPUT (policy ACCEPT)
target     prot opt source               destination         

Chain FORWARD (policy ACCEPT)
target     prot opt source               destination         
DOCKER-ISOLATION  all  --  anywhere             anywhere            
DOCKER     all  --  anywhere             anywhere            
ACCEPT     all  --  anywhere             anywhere             ctstate RELATED,ESTABLISHED
ACCEPT     all  --  anywhere             anywhere            
ACCEPT     all  --  anywhere             anywhere            

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination         

Chain DOCKER (1 references)
target     prot opt source               destination         
ACCEPT     tcp  --  anywhere             172.17.0.3           tcp dpt:6379

Chain DOCKER-ISOLATION (1 references)
target     prot opt source               destination         
RETURN     all  --  anywhere             anywhere  

#### nat表
$ iptables -t nat -L
Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination         
PREROUTING_direct  all  --  anywhere             anywhere            
PREROUTING_ZONES_SOURCE  all  --  anywhere             anywhere            
PREROUTING_ZONES  all  --  anywhere             anywhere            
DOCKER     all  --  anywhere             anywhere             ADDRTYPE match dst-type LOCAL

Chain INPUT (policy ACCEPT)
target     prot opt source               destination         

Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination         
OUTPUT_direct  all  --  anywhere             anywhere            
DOCKER     all  --  anywhere             anywhere             ADDRTYPE match dst-type LOCAL

Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination         
RETURN     all  --  192.168.122.0/24     base-address.mcast.net/24 
RETURN     all  --  192.168.122.0/24     255.255.255.255     
MASQUERADE  tcp  --  192.168.122.0/24    !192.168.122.0/24     masq ports: 1024-65535
MASQUERADE  udp  --  192.168.122.0/24    !192.168.122.0/24     masq ports: 1024-65535
MASQUERADE  all  --  192.168.122.0/24    !192.168.122.0/24    
MASQUERADE  all  --  anywhere             anywhere             ADDRTYPE match src-type LOCAL
POSTROUTING_direct  all  --  anywhere             anywhere            
POSTROUTING_ZONES_SOURCE  all  --  anywhere             anywhere            
POSTROUTING_ZONES  all  --  anywhere             anywhere            
MASQUERADE  tcp  --  172.17.0.3           172.17.0.3           tcp dpt:6379
MASQUERADE  all  --  172.17.0.0/16       !172.17.0.0/16       

Chain DOCKER (2 references)
target     prot opt source               destination         
DNAT       tcp  --  anywhere             anywhere             tcp dpt:6379 to:172.17.0.3:6379

Chain OUTPUT_direct (1 references)
target     prot opt source               destination         

Chain POSTROUTING_ZONES (1 references)
target     prot opt source               destination         
POST_public  all  --  anywhere             anywhere            [goto] 
POST_public  all  --  anywhere             anywhere            [goto] 

Chain POSTROUTING_ZONES_SOURCE (1 references)
target     prot opt source               destination         

Chain POSTROUTING_direct (1 references)
target     prot opt source               destination         

Chain POST_public (2 references)
target     prot opt source               destination         
POST_public_log  all  --  anywhere             anywhere            
POST_public_deny  all  --  anywhere             anywhere            
POST_public_allow  all  --  anywhere             anywhere            

Chain POST_public_allow (1 references)
target     prot opt source               destination         

Chain POST_public_deny (1 references)
target     prot opt source               destination         

Chain POST_public_log (1 references)
target     prot opt source               destination         

Chain PREROUTING_ZONES (1 references)
target     prot opt source               destination         
PRE_public  all  --  anywhere             anywhere            [goto] 
PRE_public  all  --  anywhere             anywhere            [goto] 

Chain PREROUTING_ZONES_SOURCE (1 references)
target     prot opt source               destination         

Chain PREROUTING_direct (1 references)
target     prot opt source               destination         

Chain PRE_public (2 references)
target     prot opt source               destination         
PRE_public_log  all  --  anywhere             anywhere            
PRE_public_deny  all  --  anywhere             anywhere            
PRE_public_allow  all  --  anywhere             anywhere            

Chain PRE_public_allow (1 references)
target     prot opt source               destination         

Chain PRE_public_deny (1 references)
target     prot opt source               destination         

Chain PRE_public_log (1 references)
target     prot opt source               destination   
```

初始情况下进入容器无法访问外网(但是外面可以访问容器内部的服务), 导致yum无法更新, 可以执行如下命令解决.

...呃, 好像执行一次就行了, 以后再清空iptables, 重启`dockerd`也不会再出现无法访问的问题了, 奇怪

```
$ iptables -t nat -A POSTROUTING -s 172.17.0.0/16 ! -d 172.17.0.0/16 -j MASQUERADE
```