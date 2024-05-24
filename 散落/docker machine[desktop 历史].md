# docker-machine

参考文章

1. [官方文档](https://docs.docker.com/machine/)
2. [docker之docker-machine用法](https://www.cnblogs.com/jsonhc/p/7784466.html)
3. [appium/appium-docker-android](https://github.com/appium/appium-docker-android)
    - 使用 docker-machine 创建虚拟机, 并通过虚拟机接口让容器拥有 usb 的访问权限.
    - 提到`generic`驱动, 其实就是在已存在的 linux 主机上安装 docker, 需要目标主机的 ip, ssh登录密钥等信息, 没什么意思
4. [Docker Machine](https://www.runoob.com/docker/docker-machine.html)
    - 常用命令列表

## docker-machine 是做什么的?

在 docker 1.12 之前(不包括1.12), 在 win/mac 下运行 docker 服务需要使用 docker-machine 工具.

直到 docker 1.12 的 beta 版本, 才有了[Docker Desktop for Mac](https://docs.docker.com/docker-for-mac/)和[Docker Desktop for Windows](https://docs.docker.com/docker-for-windows/)

1. 在比较旧的 win/mac 主机上, 使用 docker-machine 可能比 docker desktop 更合适.
2. docker-machine 可以批量管理多个主机上的 docker 服务, 而且 docker desktop 只能创建单个 docker 服务, 所有容器运行在同一个局域网内. 而 docker-machine 可以在本地创建多个.
3. docker desktop 封装在虚拟机中, 但是这个虚拟机对底层设备访问权限不足, 比如 usb 设备连接就无法被直接挂载到 docker desktop 创建的容器中.
    - 比如参考文章3就是一个需要这样能力的场景
4. docker desktop 创建的 docker 服务默认运行在一个虚拟机中, 但是通常进入这个虚拟机进行管理操作是非常困难的, docker-machine 就不一样了.

## 关于创建

可以说, docker-machine 就是先创建一个虚拟机, 然后在这个虚拟机上部署 docker 服务. 

docker desktop for win使用的虚拟化技术为`Hyper-V`, 所以需要通过`Hyper-V`驱动创建虚拟机. 或者可以不安装 desktop, 只安装 docker-machine, 然后安装 virtualbox, 通过`virtualbox`驱动创建虚拟机.

而 docker desktop for mac 的虚拟化技术则为`HyperKit`, 但是 docker-machine 还不支持直接通过`HyperKit`创建虚拟机, 所以只能使用`virtualbox`驱动, 同样要求本地事先存在`virtualbox`.

docker-machine 的安装步骤中, 第1步就写着先安装 docker, 不过看ta文档中的介绍是不需要的, 毕竟是和 docker desktop 平级的东西.

### mac

在一个只安装了 desktop 的 mac 系统中, 执行如下命令失败了, 应该是因为我没有安装`virtualbox`吧.

```log
$ docker-machine create --driver virtualbox default
Creating CA: /Users/general/.docker/machine/certs/ca.pem
Creating client certificate: /Users/general/.docker/machine/certs/cert.pem
Running pre-create checks...
Error with pre-create check: "exit status 126"
```

安装`virtualbox`后, 再次执行创建命令

```log
$ docker-machine create --driver virtualbox virtualbox
Running pre-create checks...
(virtualbox) Image cache directory does not exist, creating it at /Users/general/.docker/machine/cache...
(virtualbox) No default Boot2Docker ISO found locally, downloading the latest release...
(virtualbox) Latest release for github.com/boot2docker/boot2docker is v19.03.12
(virtualbox) Downloading /Users/general/.docker/machine/cache/boot2docker.iso from https://github.com/boot2docker/boot2docker/releases/download/v19.03.12/boot2docker.iso...
Creating machine...
(virtualbox) Copying /Users/general/.docker/machine/cache/boot2docker.iso to /Users/general/.docker/machine/machines/virtualbox/boot2docker.iso...
(virtualbox) Creating VirtualBox VM...
(virtualbox) Creating SSH key...
(virtualbox) Starting the VM...
(virtualbox) Check network to re-create if needed...
(virtualbox) Found a new host-only adapter: "vboxnet1"
(virtualbox) Waiting for an IP...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with boot2docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env virtualbox
```

创建成功, 查看虚拟机列表

```
$ docker-machine ls
NAME         ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER      ERRORS
virtualbox   -        virtualbox   Running   tcp://192.168.99.100:2376           v19.03.12
```

此时打开`virtualbox`, 可以在主界面上看到如下信息.

![](https://gitee.com/generals-space/gitimg/raw/master/93d71b0f28303c3edfba8652aeafa768.png)

> 注意: 由于是直接使用`docker-machine`通过`virtualbox`创建的虚拟机, 所以此时完全不需要 docker desktop 启动.

> `docker-machine`创建的虚拟机与 docker desktop 的不互通, 镜像与容器也是不互通的, 所以需要单独下载镜像.

------

将本地目录拷贝到目标 docker 虚拟机中

```
docker-machine scp -r ./scripts/ virtualbox:/root/
```
