# docker build构建镜像时出现Could not resolve host[dns daemon.json resolver.conf]

参考文章

1. [Cannot resolve host on docker build](https://github.com/moby/moby/issues/16600)
    - `daemon.json`配置`{"dns": ["8.8.8.8"]}`
2. [Daemon DNS options](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-dns-options)
    - 官方文档
    - dockerd --dns 8.8.8.8
3. [Docker build "Could not resolve 'archive.ubuntu.com'" apt-get fails to install anything](https://stackoverflow.com/questions/24991136/docker-build-could-not-resolve-archive-ubuntu-com-apt-get-fails-to-install-a)


## 问题描述

下载镜像仓库的 repo 文件时出现"Could not resolve host"域名无法解析的报错.

```log
 => ERROR [ 6/20] RUN curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.163.com/.help/CentOS7-Base-163.repo                                                              10.8s
------
 > [ 6/20] RUN curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.163.com/.help/CentOS7-Base-163.repo:
#0 0.243   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
#0 0.243                                  Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:--  0:00:09 --:--:--     0curl: (6) Could not resolve host: mirrors.163.com; Unknown error
------
```

换了阿里云, 163和清华大学的, 都不行, 手动将这些 repo 下载下来然后用`COPY`指令拷贝进去, 但是在 yum install 还是会报错.

不过宿主机可以 ping 通这些地址, 所以网络本身是没有问题的.

尝试修改了`/etc/resolv.conf`的dns服务器, 但不管是`114.114.114.114`还是`8.8.8.8`都还是不行.

## 解决方法

按照参考文章1, 2, 3中所说, 需要在`/etc/docker/daemon.json`配置中, 显式设置dns.

```json
{
  "dns": ["8.8.8.8"]
}
```

然后重启 docker 服务即可.

```
systemctl restart docker
```
