---
title: kubernetes部署(五)-kuber组件编译
tags: [kubernetes]
categories: general
---

<!--

# kubernetes部署(五)-kuber组件编译

<!tags!>: <!kubernetes!>

<!keys!>: uuJqwmbU68ezv&gz

-->


参考文章

1. [kubernetes1.6 安装master（二）](http://blog.csdn.net/u010278923/article/details/71126246)

参考文章1中描述了使用官方提供的方法直接下载适合当前系统的二进制发布版的方法, 直接得到编译好的可执行文件. kubernetes官方的git仓库release页面没有提供下载链接, 而是需要执行[在线安装脚本](https://get.k8s.io/).

但为了指定版本, 我还是觉得手动编译比较好(其实我只是没找到这个方法, 手动编译心好累...X﹏X)

> 2G内存不太够, 4G可以. 

------

首先安装go环境.

下载[go1.8](https://dl.gocn.io/golang/1.8.3/go1.8.3.linux-amd64.tar.gz). 解压, 将生成的`go`目录整个放在`/usr/local/go`目录下. 配置如下环境变量到`/etc/profile`

```bash
export GOROOT=/usr/local/go
export PATH=$PATH:GOROOT/bin
export GOPATH=/root/gopath
```

在`/root`下创建`gopath`目录, 随便在哪, 只要是`GOPATH`环境变量指定的路径就行.

下载git, 使用`go get`安装第三方包时默认使用git下载. 为git配置代理, 不然github下载很废时间. 这个看个人需要了.

```
$ git config --global http.proxy 'socks5://172.32.100.1:6060'
$ git config --global https.proxy 'socks5://172.32.100.1:6060'
```

然后执行`go get -d k8s.io/kubernetes`下载kuber所需的所有依赖, 下载的文件路径在`$GOPATH/src/k8s.io/kubernetes`. 

我发现它好像只是调用git工具clone了一个`kubernetes`的git仓库, 没有下载其他依赖, 在想能不能直接git clone. 当然, 直接在浏览器中下载`zip`包是不行的, 编译失败.

```
$ cd /root/gopath/src/k8s.io/kubernetes && make
```

生成的可执行文件在`_output/bin`目录下.

Master节点需要的是`kube-apiserver`, `kube-controller-manager`和`kube-scheduler`, 当然, 还有`kubctl`.

Minion节点需要的是`kubelet`和`kube-proxy`.

把它们所需的可执行文件都放在各自的`/usr/local/kubernetes/bin`目录下, 并加入环境变量.

它们各自的服务脚本及解释在`kube服务脚本.md`中. 请查看`kubernetes部署(六)-kuber组件服务脚本`

启动完成后, 在主节点执行如下命令验证

```
$ kubectl -s http://172.32.100.71:8080 get nodes
NAME            STATUS    AGE       VERSION
172.32.100.81   Ready     1m        v1.7.1-beta.0.2+09955ec93bcfc1
172.32.100.91   Ready     1m        v1.7.1-beta.0.2+09955ec93bcfc1
```