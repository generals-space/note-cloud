参考文章

1. [Docker已经再见，替代 Docker 的五种容器选择](https://cloud.tencent.com/developer/article/1422822)
    - [apache/mesos](https://github.com/apache/mesos) C++
    - [rkt/rkt](https://github.com/rkt/rkt) 该项目已结束, CoreOS公司最早发起.
    - [docker/docker-ce](https://github.com/docker/docker-ce)
    - LXC 容器. 不支持与 kuber 整合, 没有实现 OCI 的标准
2. [opencontainers/runc](https://github.com/opencontainers/runc)
    - 之前 docker 旗下的 [libcontainer](https://github.com/docker-archive/libcontainer)
3. [containerd/containerd](https://github.com/containerd/containerd)
    - 之前 docker 旗下的 [containerd](https://github.com/docker-archive/containerd)
    - 后来集成了自家的 [containerd/cri](https://github.com/containerd/cri)
4. [Docker背后的标准化容器执行引擎——runC](https://blog.csdn.net/HarmonyCloud_/article/details/125999479)
    - runc 把原本的 libcontainer 当成一个包放到自己仓库里了, 其实主要功能还是由这个包实现的.
5. [docker，containerd，runc，docker-shim之间的关系](https://blog.51cto.com/zhangxueliang/4945674)
    - 配图真是清晰明了
    - docker-shim 应该改成 containerd-shim
6. [K8S 弃用 Docker 了？Docker 不能用了？别逗了！](https://moelove.info/2020/12/03/K8S-%E5%BC%83%E7%94%A8-Docker-%E4%BA%86Docker-%E4%B8%8D%E8%83%BD%E7%94%A8%E4%BA%86%E5%88%AB%E9%80%97%E4%BA%86/)
    - dockershim 一直都是 Kubernetes 社区为了能让 Docker 成为其支持的容器运行时，所维护的一个兼容程序。 
    - 2016 年, docker 发布 swarm, 向上发展, kubernetes 发布 CRI 标准, 向下发展 - 基本相当于双方正式开战了...
7. [kubernetes真要放弃docker吗?](https://zhuanlan.zhihu.com/p/333367514)
    - dockershim 之后会在 kubernetes 之外独立维护[cri-dockerd](https://github.com/Mirantis/cri-dockerd)
    - kubernetes 创立之初, docker已经是容器领域事实的老大了，kubernetes想要发展壮大，就必须对docker大力支持，所以当时就在kubelet上开发了docker shim。
    - 有人说kubernetes现在翅膀硬了，就要甩开docker，这种说法也能说得过去。
8. [终于可以像使用 Docker 一样丝滑地使用 Containerd 了](https://zhuanlan.zhihu.com/p/364206329)
    - Kubernetes 虽然制定了容器运行时接口（CRI）标准，但早期能用的容器运行时只有 Docker，而 Docker 又不适配这个标准，于是给 Docker 开了后门，花了大量的精力去适配它。
    - 后来有了更多的容器运行时可以选择后，Kubernetes 就不得不重新考量要不要继续适配 Docker 了，因为每次更新 Kubelet 都要考虑与 Docker 的适配问题。
9. [Docker，containerd，CRI，CRI-O，OCI，runc 分不清？看这一篇就够了](https://zhuanlan.zhihu.com/p/490585683)
    - 图不错, 解释了 cri-o 与 runc 的关系.
10. [名词解释：OCI、CRI、ContainerD、CRI-O以及runC](https://zhuanlan.zhihu.com/p/468495520)

docker info: 19.03.5

相关的可执行文件有:

```console
$ ls /usr/bin/ | grep docker
docker
dockerd
docker-init
docker-proxy
$ ls /usr/bin/ | grep container
containerd
containerd-shim
$ ls /usr/bin/ | grep runc
runc
```

```console
$ ps -ef | grep containerd
root       1258      1  0 05:58 ?        00:00:44 /usr/bin/containerd
root       1263      1  1 05:58 ?        00:06:17 /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock
root       1942   1258  0 05:58 ?        00:00:04 containerd-shim -namespace moby -workdir /var/lib/containerd/io.containerd.runtime.v1.linux/moby/ad7a7a5952fd0f1b6637d49cdb673d73b73b65d750f21a734b133d8e07e25b98 -address /run/containerd/containerd.sock -containerd-binary /usr/bin/containerd -runtime-root /var/run/docker/runtime-runc -systemd-cgroup
```

可以看到, `dockerd`与`containerd`是并列的, `containerd-shim`(每启动一个docker容器都会启动一个shim进程)则是`containerd`的子进程, 容器中的`CMD/ENTRYPOINT`执行命令是由`containerd-shim`启动执行的. 

如下, 容器`ad7a7a5952fd0`中的nginx进程就是对应`container-shim`的子进程.

```console
$ ps -ef | grep 1942
root       1942   1258  0 05:58 ?        00:00:04 containerd-shim -namespace moby -workdir /var/lib/containerd/io.containerd.runtime.v1.linux/moby/ad7a7a5952fd0f1b6637d49cdb673d73b73b65d750f21a734b133d8e07e25b98 -address /run/containerd/containerd.sock -containerd-binary /usr/bin/containerd -runtime-root /var/run/docker/runtime-runc -systemd-cgroup
root       1961   1942  0 05:58 ?        00:00:00 nginx: master process nginx -g daemon off;
```

`dockerd`在启动时可以指定`runc`的实现, 使用`docker info`也可以查到`containerd`和`runc`的版本信息.

## 纯 docker 

```
+----------+        +-----------+ grpc  +-----------+
|docker-cli| -----> |  dockerd  | ----> | containerd|
+----------+        +-----------+       +-----┬-----+
                                              | exec
                                    ┌─────────┴─────────┐
                            +-------↓-------+   +-------↓-------+
                            |containerd-shim|   |containerd-shim|
                            +-------┬-------+   +-------┬-------+
                                    | exec              | exec
                              +-----↓-----+       +-----↓-----+
                              |    runc   |       |    runc   |
                              +-----------+       +-----------+
```

> 这里把`dockerd`和`containerd`放在了同一级, 不过其实ta们是存在调用顺序的.

## docker + kubernetes.v1.24-(1.24之前 )

```
                                  +-----------+   
                                  |  kubelet  |   
                                  +-----┬-----+   
                                        |         
                              +---------↓--------+
                              |  GenericRuntime  |
                              +---------┬--------+
                          ┌─────────────┴─────────────┐ 
                    +-----↓------+              +-----↓----+ 
                    | dockershim |              | cri-shim | 
                    +-----┬------+              +-----┬----+ 
                          |                           |
                          |              +------------------------+
                          |              | containerd |    rkt    |
                          |              +------------------------+
                          |
+----------+        +-----↓-----+ grpc  +-----------+
|docker-cli| -----> |  dockerd  | ----> | containerd|
+----------+        +-----------+       +-----┬-----+
                                              | exec
                                    ┌─────────┴─────────┐
                            +-------↓-------+   +-------↓-------+
                            |containerd-shim|   |containerd-shim|
                            +-------┬-------+   +-------┬-------+
                                    | exec              | exec
                              +-----↓-----+       +-----↓-----+
                              |    runc   |       |    runc   |
                              +-----------+       +-----------+
```

kubernetes 最开始出现时, 是与 docker 强绑定的(当时没有其他容器化实现), kubelet 与 dockerd 直接通信.

后来才出现了 docker 以外的其他 runtime, 如 runv, rkt. 

2016年, kubernetes 官方发布了 cri 接口规范, 规范所有运行时接口. 但此时 docker 也发布了 swarm, 进行容器编排. 一个由上往下, 一个由下向上, 都向对方发起正义的背刺😂.

docker 没有理会这个 cri, kubernetes 官方只能自己写了个`dockershim`包, 给 docker 服务提供了 cri 接口. 

kubelet 在启动时, 会先创建与 dockerd 服务(/var/run/docker.sock)的连接对象. 然后启动名为 dockershim 的 grpc server, kubelet 对容器的各种操作, 都是向该 grpc server 发出请求(就是调用 grpc 服务中提供的 Service 的函数), dockershim 服务会将请求转发给 dockerd.

`GenericRuntime`是一个通用接口, 可以与任何实现了 cri 接口的 runtime 通信, 我们可以自行指定一个其他实现了 CRI 接口的 runtime, 把 dockerd 替换掉.

## docker + kubernetes.v1.24+(1.24及之后)

```
                                                        +-----------+   
                                                        |  kubelet  |   
                                                        +-----┬-----+   
                                                              |         
                                                    +---------↓--------+
                                                    |  GenericRuntime  |
                                                    +---------┬--------+
                                              ┌───────────────┴───────────────┐
                                              |                               |
+----------+        +-----------+ grpc  +-----↓-----+                     +---↓---+
|docker-cli| -----> |  dockerd  | ----> | containerd|                     | cri-o |
+----------+        +-----------+       +-----┬-----+                     +---┬---+
                                              | exec                          |
                                    ┌─────────┴─────────┐                     |
                            +-------↓-------+   +-------↓-------+             |
                            |containerd-shim|   |containerd-shim|             |
                            +-------┬-------+   +-------┬-------+             |
                                    | exec              | exec                |
                              +-----↓-----+       +-----↓-----+         +-----↓-----+
                              |    runc   |       |    runc   |         |    runc   |
                              +-----------+       +-----------+         +-----------+
```

1.24的修改, 其实就是把 dockershim 从 kubelet 源码中移除了, 直接与 containerd 服务进行通信(因为 containerd 实现了 CRI, ta 集成了自家的 [containerd/cri](https://github.com/containerd/cri)), 不再让 dockerd 这中间商赚差价了.

可以说, kubernetes 发达后, 就一脚把 docker 踹开了. 倒是 containerd 是 docker 开源的, 捐给 CNCF 组织后, 实现了 CRI, 也有点格局大了的意思.

CRI-O也是一个CRI的实现，它来自于Red Hat/IBM.
