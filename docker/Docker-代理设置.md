# Docker-pull命令代理设置

参考文章

1. [docker 使用 socks5代理](http://www.jianshu.com/p/fef11e46ebf1)

2. [Docker - 国内镜像的配置及使用](http://www.cnblogs.com/anliven/p/6218741.html)

3. [Docker Hub 镜像站点](https://cr.console.aliyun.com/#/accelerator)

```
Environment=HTTP_PROXY=http://172.32.100.1:6060
Environment=HTTPS_PROXY=http://172.32.100.1:6060
Environment=NO_PROXY=localhost,127.0.0.1,m1empwb1.mirror.aliyuncs.com,registry.cn-hangzhou.aliyuncs.com,daocloud.io
```