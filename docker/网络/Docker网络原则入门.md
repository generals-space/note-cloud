# Docker网络原则入门

原文链接

[Docker网络原则入门：EXPOSE, -p, -P, -link](http://dockone.io/article/455)

翻译自

[A Brief Primer on Docker Networking Rules: EXPOSE, -p, -P, –link](A Brief Primer on Docker Networking Rules: EXPOSE, -p, -P, –link)

> 构建多容器应用程序, 需要定义网络参数来设置容器间的通信, 可以通过EXPOSE或者-expose暴露端口、使用-p发布特定端口, 还可以用-link等等来实现, 这些方法可能会得到一样的效果, 但是这些方法之间是否有不同, 应该选择什么样的方法, 将是本文讨论的重点内容.

如果你已经构建了一些多容器的应用程序, 那么肯定需要定义一些网络规则来设置容器间的通信. 有多种方式可以实现：可以通过`--expose`参数在运行时暴露端口, 或者在Dockerfile里使用`EXPOSE`指令. 还可以在`docker run`的时候通过`-p`或者`-P`参数来发布端口. 或者通过`--link`链接容器. 虽然这些方式几乎都能达到一样的结果, 但是它们还是有细微的区别. 那么到底应该使用哪一种呢？

> 使用-p或者-P来创建特定端口绑定规则最为可靠, EXPOSE可以看做是容器文档化的方式, 谨慎使用--link的方式.

在比较这些不同方式之前, 我们先分别了解细节.

## 1. 通过EXPOSE或者-expose暴露端口

有两种方式可以用来暴露端口：要么在Dockerfile里用`EXPOSE 目标端口号`定义, 要么在`docker run`时指定`--expose=目标端口号`. 这两种方式作用相同, 但是, `--expose`可以接受端口范围作为参数, 比如 `--expose=2000-3000`. 但是, `EXPOSE`和`--expose`都不依赖于宿主机器. 默认状态下, 这些规则并不会使这些端口可以通过宿主机来访问.

基于`EXPOSE`指令的上述限制, Dockerfile的作者一般在包含`EXPOSE`规则时都只将其作为 **哪个端口提供哪个服务的提示**. 使用时, 还要依赖于容器的操作人员进一步指定网络规则. 和`-P`参数联合使用的情况, 下文会进一步阐述. 不过通过`EXPOSE`命令文档化端口的方式十分有用.

**本质上说, `EXPOSE`或者`--expose`只是为其他命令提供所需信息的元数据, 或者只是告诉容器操作人员有哪些已知选择.**

实际上, 在运行时暴露端口和通过Dockerfile的指令暴露端口, 这两者没什么区别. 在这两种方式启动的容器里, 通过`docker inspect $container_id | $container_name`查看到的网络配置是一样的：

```json
"NetworkSettings": {
  "PortMapping": null,
  "Ports": {
      "目标端口号/tcp": null
  }
},
"Config": {
  "ExposedPorts": {
      "目标端口号/tcp": {}
  }
}
```

可以看到端口被标示成已暴露, 但是没有定义任何映射. 注意这一点, 因为我们查看的是发布端口.

*ProTip：使用运行时标志`--expose`是附加的, 因此会在Dockerfile的`EXPOSE`指令定义的端口之外暴露添加的端口. *

## 2. 使用-p映射特定端口

可以使用`-p`参数显式将一个或者一组端口从容器里绑定到宿主机上, 而不仅仅是提供一个端口. 注意这里是 **小写的p**, 不是大写. 因为该配置依赖于宿主机器, 所以Dockerfile里没有对应的指令, 这是运行时才可用的配置. -p参数有几种不同的格式：

- 宿主机IP:宿主机端口:docker容器端口

- 宿主机IP::docker容器端口

- 宿主机端口:docker容器端口

- docker容器端口

实际中, 可以忽略ip或者hostPort, 但是必须要指定需要暴露的`containerPort(容器端口)`. 另外, 所有这些发布的规则都默认为tcp. 如果需要udp, 需要在最后加以指定, 比如`-p 1234:1234/udp`. 如果只是用命令`docker run -p 8080:3000 my-image`运行一个简单的应用程序, 那么容器里运行在3000端口的服务在宿主机的8080端口也就可用了. 端口不需要一样, 但是 **在多个容器都暴露端口时, 必须注意避免端口冲突**.

避免冲突的最佳方法是 **让Docker自己分配hostPort**. 在上述例子里, 可以选择`docker run -p 3000 my_image`来运行容器, 而不是显式指定宿主机端口. 这时, Docker会帮助选择一个宿主机端口. 运行命令`docker port $container_id | $container_name`可以查看Docker选出的端口号. 除了端口号, 该命令只能显示容器运行时端口绑定信息. 还可以通过在容器上运行`docker inspect`查看详细的网络信息, 在定义了端口映射时, 这样的信息就很有用. 该信息在`Config`、`HostConfig`和`NetworkSettings`部分. 我们查看这些信息来对比不同方式搭建的容器间的网络区别.

*ProTip：可以用`-p`参数定义任意数量的端口映射. *

## 3. -expose/EXPOSE和-p对比

为了更好得理解两者之间的区别, 我们使用不同的端口设置来运行容器.

运行一个很简单的应用程序, 会在curl它的时候打印'hello world'. 称这个镜像为`no-exposed-ports`:

```
FROM ubuntu:trusty
MAINTAINER Laura Frank <laura.frank@centurylink.com>
CMD while true; do echo 'hello world' | nc -l -p 8888; done
```

实验时注意使用的是Docker主机, 而不是`boot2docker`. 如果使用的是`boot2docker`, 运行本文示例命令前先运行`boot2docker ssh`.

注意, 我们使用-d参数运行该容器, 因此容器一直在后台运行. （端口映射规则只适用于运行着的容器）：

```
$ docker run -d --name no-exposed-ports no-exposed-ports
e18a76da06b3af7708792765745466ed485a69afaedfd7e561cf3645d1aa7149
```

这儿没有太多的信息, 只是回显了容器的ID, 提示服务已经成功启动. 和预期结果一样, 运行`docker port no-exposed-ports`和`docker inspect no-exposed-ports`时没显示什么信息, 因为我们既没有定义端口映射规则也没有发布任何端口.

因此, 如果我们发布一个端口会发生什么呢, `-p`参数和`EXPOSE`到底有什么区别呢？

还是使用上文的`no-exposed-ports`镜像, 在运行时添加`-p`参数, 但是不添加任何`expose`规则. 在`config.ExposedPorts`里重新查看`--expose`参数或者`EXPOSE`指令的结果.

```
$ docker run -d --name no-exposed-ports-with-p-flag -p 8888:8888 no-exposed-ports
c876e590cfafa734f42a42872881e68479387dc2039b55bceba3a11afd8f17ca
$ docker port no-exposed-ports-with-p-flag
8888/tcp -> 0.0.0.0:8888
```

太棒了！我们可以看到可用端口. 注意默认这是tcp. 我们到网络设置里看看还有什么信息：

```json
"Config": {
  [...]
  "ExposedPorts": {
      "8888/tcp": {}
  }
},
"HostConfig": {
  [...]
  "PortBindings": {
      "8888/tcp": [
          {
              "HostIp": "",
              "HostPort": "8888"
          }
      ]
  }
},
"NetworkSettings": {
  [...]
  "Ports": {
      "8888/tcp": [
          {
              "HostIp": "0.0.0.0",
              "HostPort": "8888"
          }
      ]
  }
}
```

注意`Config`部分的`ExposedPorts`字段. 这和我们使用`EXPOSE`或者`--expose`暴露的端口是一致的. Docker会隐式暴露已经发布的端口. 已暴露端口和已发布端口的区别在于已发布端口在宿主机上可用, 而且我们可以在`HostConfig`和`NetworkSettings`两个部分都能看到已发布端口的信息.

所有发布（`-p`或者`-P`）的端口都暴露了, 但是并不是所有暴露（`EXPOSE`或`--expose`）的端口都会发布.

## 4. 使用-P和EXPOSE发布端口

因为`EXPOSE`通常只是作为记录机制, 也就是告诉用户哪些端口会提供服务, Docker可以很容易地把Dockerfile里的EXPOSE指令转换成特定的端口绑定规则. 只需要在运行时加上`-P`参数, Docker会自动为用户创建端口映射规则, 并且帮助避免端口映射的冲突.

添加如下行到上文使用的Web应用Dockerfile里：

```
EXPOSE 1000
EXPOSE 2000
EXPOSE 3000
```

构建镜像, 命名为exposed-ports.

```
docker build -t exposed-ports .
```

再次用-P参数运行, 但是不传入任何特定的-p规则. 可以看到Docker会将EXPOSE指令相关的每个端口映射到宿主机的端口上：

```
$ docker run -d -P --name exposed-ports-in-dockerfile exposed-ports
63264dae9db85c5d667a37dac77e0da7c8d2d699f49b69ba992485242160ad3a
$ docker port exposed-ports-in-dockerfile
1000/tcp -> 0.0.0.0:49156
2000/tcp -> 0.0.0.0:49157
3000/tcp -> 0.0.0.0:49158
```

很方便, 不是么？

...**再说一遍, expose指令暴露的端口无法让其他主机通过宿主机的IP访问得到**

## 5. --link怎么样呢？

你可能在多容器应用程序里使用过运行时参数 `--link name：alias`来设定容器间关系. 虽然`--link`非常易于使用, 几乎能提供和端口映射规则和环境变量相同的功能. 但是最好将`--link`当做服务发现的机制, 而不是网络流量的门户.

`--link`参数唯一多做的事情是会使用源容器的主机名和容器ID来更新新建目标容器（使用`--link`参数创建的容器）的`/etc/hosts`文件.

当使用--link参数时, Docker提供了一系列标准的环境变量, 如果想知道细节的话可以查看相应文档.

虽然`--link`对于需要隔离域的小型项目非常有用, 它的功能更像服务发现的工具. 如果项目中使用了编排服务, 比如Kubernetes或者Fleet, 很可能就会使用别的服务发现工具来管理关系. 这些编排服务可能会不管理Docker的链接, 而是管理服务发现工具里包含的所有服务, 在Panamax项目里使用的很多远程部署适配器正是做这个的.

## 6. 找到平衡

哪一种网络选择更为适合, 这取决于谁（或者哪个容器）使用Docker运行的服务. 需要注意的是一旦镜像发布到`Docker Hub`之后, 你无法知道其他人如何使用该镜像, 因此要尽可能让镜像更加灵活. 如果你只是从Docker Hub里取得镜像, 使用-P参数运行容器是最方便迅速的方式, 来基于作者的建议创建端口映射规则. 记住每一个发布的端口都是暴露端口, 但是反过来是不对的.

## 7. 快速参考

|    命令    |                                    功能                                    |
|:--------:|:------------------------------------------------------------------------:|
|  EXPOSE  |                       记录服务可用的端口, 但是并不创建与宿主机之间的端口映射                       |
| --expose |                        运行时暴露端口, 但是并不创建与宿主机之间的端口映射                        |
|    -p    |                 创建容器时同时创建端口映射规则, 比如`-p 宿主机IP:宿主机端口:容器端口`                 |
|    -P    |                   将`Dockerfile`中暴露的所有容器端口全部动态映射到宿主机的端口                   |
|  --link  | 在`消费`和`服务`容器之间创建链接, 这会创建一系列环境变量, 并在消费者容器的`/etc/hosts`文件里田间入口项. 必须暴露或发布端口 |
