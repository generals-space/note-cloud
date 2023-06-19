# centos 6.x中安装旧版docker[yum rpm]

参考文章

1. [CentOS 6.8上安装 docker.io](https://developer.aliyun.com/article/585178)
2. [epel-release-6-8.noarch.rpm](http://dl.fedoraproject.org/pub/archive/epel/6/x86_64/)
3. [Installing Docker on CentOS 6 after removal of docker-io](https://stackoverflow.com/questions/55134196/installing-docker-on-centos-6-after-removal-of-docker-io)

20230607

参考文章1中的epel包已经找不到了.

参考文章1中能找到, 但是已经没有 docker-io 的包了, 官方将其从仓库中移除了.

普通的epel仓库的版本只有1.5的, `docker-1.5-5.el6.x86_64.rpm`.

按照参考文章3中的说法, 可以从"get.docker.com"直接找rpm包安装.

```
$ yum install https://get.docker.com/rpm/1.7.1/centos-6/RPMS/x86_64/docker-engine-1.7.1-1.el6.x86_64.rpm

$ docker --version
Docker version 1.7.0, build 0baf609
```
