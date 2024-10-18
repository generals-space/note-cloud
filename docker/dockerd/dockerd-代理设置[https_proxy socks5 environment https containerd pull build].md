# dockerd-代理设置

参考文章

1. [docker 使用 socks5代理](http://www.jianshu.com/p/fef11e46ebf1)
2. [Docker - 国内镜像的配置及使用](http://www.cnblogs.com/anliven/p/6218741.html)
3. [Docker Hub 镜像站点](https://cr.console.aliyun.com/#/accelerator)
4. [简述关于containerd设置代理](https://blog.51cto.com/u_15343792/5142108)
    - containerd 的代理设置也一样.
5. [如何优雅的给 Docker 配置网络代理](https://cloud.tencent.com/developer/article/1806455)
    - docker pull 拉取镜像需要在 service.d 中配置`Environment`块
    - docker build 过程中可通过`--build-arg`选项设置`HTTP_PROXY`代理
    - 容器运行中, 可通过 config.json 设置全局的代理, 与启动容器时通过`-e`设置环境变量原理一样.

```conf
Environment=HTTP_PROXY=http://172.32.100.1:1081
Environment=HTTPS_PROXY=http://172.32.100.1:1081
Environment=NO_PROXY=localhost,127.0.0.1,m1empwb1.mirror.aliyuncs.com,registry.cn-hangzhou.aliyuncs.com,daocloud.io
```

需要重启 docker 服务.

> **这里的代理需要是 http 类型, 不能是 socks5 类型的.**

> 可以在`/etc/docker/daemon.json`文件中设置国内的镜像源(不过有时不好用, 还是得用这个方法).

如果出现了如下错误, 还是试试`daemon.json`吧.

```log
$ d pull pandoc/core
Using default tag: latest
Error response from daemon: Get https://registry-1.docker.io/v2/: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
```
