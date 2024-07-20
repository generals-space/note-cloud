原文链接

[K8S 弃用 Docker 了？Docker 不能用了？别逗了！](https://moelove.info/2020/12/03/K8S-%E5%BC%83%E7%94%A8-Docker-%E4%BA%86Docker-%E4%B8%8D%E8%83%BD%E7%94%A8%E4%BA%86%E5%88%AB%E9%80%97%E4%BA%86/)

## 概览

### 2013 年

Docker 是在 2013 年的 PyCon 上首次正式对外公布的. 它带来了一种先进的软件交付方式, 即, 通过容器镜像进行软件的交付. 工程师们只需要简单的 `docker build` 命令即可制作出自己的镜像, 并通过 `docker push` 将其发布至 DockerHub 上. 通过简单的 `docker run` 命令即可快速的使用指定镜像启动自己的服务. 

通过这种办法, 可以有效的解决软件运行时环境差异带来的问题, 达到其 **Build once, Run anywhere** 的目标. 

从此 **Docker 也基本成为了容器的代名词, 并成为容器时代的引领者.**

### 2014 年

2014 年 Google 推出 Kubernetes 用于解决大规模场景下 Docker 容器编排的问题. 

**这是一个逻辑选择, 在当时 Docker 是最流行也是唯一的运行时.** Kubernetes 通过对 Docker 容器运行时的支持, 迎来了大量的用户. 

同时, Google 及 Kubernetes 社区与 Docker 也在进行着密切的合作, 在其官方博客上有如下内容：

```
We’ll continue to build out the feature set, while collaborating with the Docker community to incorporate the best ideas from Kubernetes into Docker.
```

[An update on container support on Google Cloud Platform](https://cloudplatform.googleblog.com/2014/06/an-update-on-container-support-on-google-cloud-platform.html)

```
Kubernetes is an open source manager for Docker containers, based on Google’s years of experience using containers at Internet scale. Docker is delivering the full container stack that Kubernetes schedules into, and is looking to move critical capabilities upstream and align the Kubernetes framework with Libswarm.
```

[Welcome Microsoft, RedHat, IBM, Docker and more to the Kubernetes community](https://cloudplatform.googleblog.com/2014/07/welcome-microsoft-redhat-ibm-docker-and-more-to-the-kubernetes-community.html)

并在同一个月的 DockerCon 上发布演讲, 介绍了 Kubernetes 并受到了广泛的关注. 

此时 **Docker Inc. 也发布了其容器编排工具, libswarm （也就是后来的 swarmkit）.**

### 2015 年

2015 年 OCI （Open Container Initiative）由 Docker 和其他容器行业领导者共同成立（它也是 Linux 基金会旗下项目）

OCI 主要包含两个规范：

- 运行时规范（runtime-spec）：容器运行时, 如何运行指定的 文件系统上的包
- 容器镜像规范（image-spec）：如何创建一个 OCI 运行时可运行的文件系统上的包

Docker 把它自己的容器镜像格式和 runtime ( **现在的 runc** ) 都捐给了 OCI 作为初始工作. 

### 2016 年

2016 年 6 月, Docker v1.12 发布, **带来了 Docker 在多主机多容器的编排解决方案, Docker Swarm**. 这里也需要注意的是, Docker v1.12 的设计原则：

- Simple Yet Powerful （简单而强大）
- Resilient（弹性）
- Secure（安全）
- Optional Features and Backward Compatibility（可选功能及向后兼容）

所以你可以通过配置自行选择是否需要使用 Docker Swarm, 而无需担心有什么副作用. 

2016 年 12 月, **Kubernetes 发布 CRI （Container Runtime Interface）**, 这当中一部分原因是由于 Kubernetes 尝试支持另一个由 CoreOS 领导的容器运行时项目 rkt, 但是需要写很多兼容的代码之类的, 为了避免后续兼容其他运行时带来的维护工作, 所以发布了统一的 CRI 接口, 凡是支持 CRI 的运行时, 皆可直接作为 Kubernetes 的底层运行时；

当然, **Kubernetes 也是在 2016 年逐步取得那场容器编排战争的胜利的**. 

### 2017 年

2017 年, Docker 将自身从 v1.11 起开始引入的容器运行时 [containerd](https://github.com/containerd/containerd/) 捐给了 CNCF

2017 年, Docker 的网络组件 libnetwork 增加了 CNI 的支持; 同时通过[使用 Docker 为 Docker Swarm 提供的 ipvs 相关的代码](https://github.com/kubernetes/kubernetes/pull/46580), 也在 Kubernetes 中实现了基于 IPvs 的 service 负载均衡. 不过在 v1.18 中开始移除了相关的依赖. 

同年 11 月, [Kubernetes 中新增了 containerd 的支持](https://kubernetes.io/blog/2017/11/containerd-container-runtime-options-kubernetes/)

![cri-containerd](https://gitee.com/generals-space/gitimg/raw/master/2022/0347d4715e0679a6b0f9cdba9792324e.png)

### 2018 年

2018 年, [Kubernetes 的 containerd 集成, 正式 GA](https://kubernetes.io/blog/2018/05/24/kubernetes-containerd-integration-goes-ga/)

之前版本的架构：

![containerd 1.0 cri-containerd](https://gitee.com/generals-space/gitimg/raw/master/2022/846a5f3aa629b402465b1058e2c1a960.png)

新的架构：

![containerd 1.1 cri-containerd](https://gitee.com/generals-space/gitimg/raw/master/2022/3a09ae69a557cd79be2b25082acf1825.png)

### 2019 年

2019 年, 上文中提到的另一个容器运行时项目 rkt 被 CNCF 归档, 终止使命了; 2019 年 Mirantis 收购 Docker 的企业服务. 

### 2020 年

时间回到今年, Docker 主要被误会的两件事：

1. Docker Inc. 修改 DockerHub 的定价和 TOS. 国内争论较多的主要是关于合规性的问题（但是被标题党带歪了, 免不了恐慌）；
2. Kubernetes 宣布开始进入废弃 dockershim 支持的倒计时, 被人误以为 Docker 不能再用了；

## 说明

关于 DockerHub 修改定价和 TOS 的事情, 这里就不再多说了, 毕竟 DockerHub 目前大家仍然用的很欢乐, 远不像当初那些标题党宣称的那样. 

重点来说一下第二件事情吧. 

Kubernetes 当初选择 Docker 作为其容器运行时, 本身就是因为当时它没有其他的选择, 并且选择 Docker 可为它带来众多的用户. 所以, 开始时, 它便提供了内置的对 Docker 运行时的支持. 

而 Docker 其实创建之初, 并没有考虑到“编排”的这个功能, 当然也没有考虑到 Kubernetes 的存在（因为当时还没有）. 

dockershim 一直都是 Kubernetes 社区为了能让 Docker 成为其支持的容器运行时, 所维护的一个兼容程序. 本次所谓的废弃, 也仅仅是 Kubernetes 要放弃对现在 Kubernetes 代码仓库中的 dockershim 的维护支持. 以便其可以像开始时计划的那样, 仅负责维护其 CRI, 任何兼容 CRI 的运行时, 皆可作为 Kubernetes 的 runtime. 

在 Kubernetes 提出 CRI 时, 有人建议在 Docker 中实现它. 但是这种方式也会带来一个问题, 即使 Docker 实现了 CRI, 但它仍然不是一个单纯的容器运行时, 它本身包含了大量的非 “纯底层容器运行时” 所具备的功能. 

所以后来 自 Docker 中分离出来的 containerd 项目, 作为一个底层容器运行时出现了, 它是 Kubernetes 容器运行时更好的选择. 

Docker 使用 containerd 作为其底层容器运行时以及众多的云厂商及公司在生产环境中使用 containerd 作为其 Kubernetes 的运行时, 这也从侧面验证了 containerd 的稳定性. 

现在 Kubernetes 和 Docker 社区都相信 containerd 已经足够成熟可直接作为 Kubernetes 的运行时了, 而无需再通过 dockershim 使用 Docker 作为 Kubernetes 的运行时了. 这也标志着 Docker 为 Kubernetes 提供一个现代化的容器运行时的承诺最终兑现了. 

而本次事件中, 重点的 dockershim 之后的方向如何呢？Kubernetes 代码仓库中的 dockershim 将会在未来版本中移除, 但是 Mirantis 公司已经和 Docker 达成合作, 在未来会共同维护一份 dockershim 组件, 以便支持 Docker 作为 Kubernetes 的容器运行时. 

Otherwise, if you’re using the open source Docker Engine, the dockershim project will be available as an open source component, and you will be able to continue to use it with Kubernetes; it will just require a small configuration change, which we will document.

[Mirantis 公司宣布将维护 dockershim](https://www.mirantis.com/blog/mirantis-to-take-over-support-of-kubernetes-dockershim-2/)

## Q&A

Q：本次 Kubernetes 放弃对 dockershim 的维护, 到底有什么影响？ 

A：对于普通用户而言, 没有任何影响；对于在 Kubernetes 之上进行开发的工程师, 没什么太大影响；对于集群管理员, 需要考虑是否要在未来版本中, 将容器运行时, 升级为支持 CRI 的运行时, 比如 containerd. 当然, 如果你并不想切换容器运行时, 那也没关系, Mirantis 公司未来会和 Docker 共同维护 dockershim, 并作为一个开源组件提供. 

------

Q: Docker 不能用了吗？ 

A：Docker 仍然是本地开发, 或者单机部署最佳的容器工具, 它提供了更为人性化的用户体验, 并且也有丰富的特性. 目前 Docker 已经和 AWS 达成合作, 可直接通过 Docker CLI 与 AWS 集成. 另外, Docker 也仍然可以作为 Kubernetes 的容器运行时, 并没有立即中止对其支持. 

------

Q：听说 Podman 可以借机上位了？ 

A：想太多. Podman 也并不兼容 CRI, 并且它也不具备作为 Kubernetes 容器运行时的条件. 我个人也偶尔有在用 Podman, 并且我们在 KIND 项目中也提供了对 Podman 的支持, 但实话讲, 它也就是只是一个 CLI 工具, 某些情况下会有些作用, 比如如果你的 Kubernetes 容器运行时使用 cri-o 的情况下, 可以用来本地做下调试. 

## 总结

本文主要介绍了 Docker 和 Kubernetes 的发展历程, 也解释了本次 Kubernetes 仅仅是放弃其对 dockershim 组件的支持. 未来更推荐的 Kubernetes 运行时是 兼容 CRI 的 containerd 之类的底层运行时. 

Mirantis 公司将会和 Docker 共同维护 dockershim 并作为开源组件提供. 

Docker 仍然是一款最佳的本地开发测试和部署的工具. 

