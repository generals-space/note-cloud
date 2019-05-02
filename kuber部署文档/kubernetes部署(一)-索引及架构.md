---
title: kubernetes部署(一)-索引及架构
tags: [kubernetes]
categories: general
---

<!--

# kubernetes部署(一)-索引及架构

<!tags!>: <!kubernetes!>

<!keys!>: q25mt)dJizbvfaCz

-->

Master: 172.32.100.71 kube-apiserver + scheduler + controller-manager + kubectl

minion1: 172.32.100.81 kube-proxy + kubelet

minion2: 172.32.100.91 kube-proxy + kubelet

系统: Centos7.2

kubernetes: 1.7(最新发布)

docker: 1.12.6

flannel: 0.7.1

因为是摸索阶段, 所以没有启用apiserver与etcd的https证书验证, 不过貌似走了很多弯路. 也许按照标准模式反而会更简单一点(谁知道呢, 坑都踩过了 ╮(╯▽╰)╭)

## 1. docker

kuber集群Master节点与minion节点都要安装docker, 这一步是必须的. CentOS7的yum源中docker版本为`1.12.6`, 正好是kuber所支持的. 单纯作为容器, 这一版本的docker的功能已经基本健全, 更高版本只是实现了集群化的辅助功能(通信, 高可用等), 在使用kuber管理集群时, 不必再使用docker本身提供的功能, 反正也不好用(好吧...我在瞎说, 因为我没用过).

## 2. etcd

etcd部署在Master节点, 当然也可以在kuber集群之外另建集群, 只要apiserver能够连接就可以. 为了省力, 这里直接装在Master节点. 

> 在生产环境中绝对有必要搭建etcd集群.

## 3. flannel

~~放弃使用OVS, 虽然OVS有与docker结合完成跨主机容器互联的方法, 但并没有与etcd结合使用的解决方案, 与kubernetes结合的不好, 而且过于庞大.~~

关于在flannel在kubernetes中的作用. 

我们知道, 默认docker在宿主机中分配的网段为`docker0`接口上的`172.17.0.0/24`, 该主机上所有docker容器必然会在这个网段中. 不同宿主机的不同docker容器很可能会同时获得相同的`172.17.0.1`这个IP, 多主机间容器通信时不会允许这样的情况出现. 

flannel就是用来重新定义每台宿主机上的`docker0`的IP的, 它会使不同的宿主机上的docker获得不同的网段: `172.17.1.0/24`与`172.17.2.0/24`等.

## 4. kuber组件编译

这一步在任何机器上都可以执行, 因为生成的可执行文件可以直接运行在其他机器上. 

也可以直接从官方下载.

请查看`kubernetes部署(五)-kuber组件编译`

## 5. dns插件

dns服务在kuber集群中以插件形式存在, 它的作用这里就介绍了. 关于pods和service创建的先后问题, 等我把集群搭建完成后再说吧.

请查看`kubernetes部署(七)-dns插件.md`

## 6. dashboard插件

...我一直觉得, 如果看不到UI界面, kuber集群就不算搭建完. 我还觉得, 大部分人都是这么认为的.