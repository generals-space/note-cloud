# dockerd-代理设置

参考文章

1. [docker 使用 socks5代理](http://www.jianshu.com/p/fef11e46ebf1)
2. [Docker - 国内镜像的配置及使用](http://www.cnblogs.com/anliven/p/6218741.html)
3. [Docker Hub 镜像站点](https://cr.console.aliyun.com/#/accelerator)

```
Environment=HTTP_PROXY=http://172.32.100.1:1081
Environment=HTTPS_PROXY=http://172.32.100.1:1081
Environment=NO_PROXY=localhost,127.0.0.1,m1empwb1.mirror.aliyuncs.com,registry.cn-hangzhou.aliyuncs.com,daocloud.io
```

需要重启 docker 服务.

> 这里的代理端口要是 http 类型, 而不能是 socks5 类型的.

> 可以在`/etc/docker/daemon.json`文件中设置国内的镜像源(不过有时不好用, 还是得用这个方法).

如果出现了如下错误, 还是试试`daemon.json`吧.

```console
$ d pull pandoc/core
Using default tag: latest
Error response from daemon: Get https://registry-1.docker.io/v2/: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
```
