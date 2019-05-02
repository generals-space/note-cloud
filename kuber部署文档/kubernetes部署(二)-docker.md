---
title: kubernetes部署(二)-docker
tags: [kubernetes, docker]
categories: general
---

<!--

# kubernetes部署(二)-docker

<!tags!>: <!kubernetes!> <!docker!>

<!keys!>: p65CzzMweptjrrj{

-->



CentOS7下直接使用yum安装, 反正kubernetes也没对更高版本的docker做过兼容测试, 而且yum源中的docker版本绝对够了.

```
$ yum install docker -y
```

在启动docker服务之前, 需要修改docker数据存放路径, 默认在`/var/lib/docker`, 按照斯凯生产环境的部署风格, `/var`所在目录一般是`/`, 空间较小, 需要手动将其指定到`/opt`, `app`这样的大容量分区.

编辑`/usr/lib/systemd/system/docker.service`, 在`ExecStart`字段添加启动参数`--graph=/opt/docker`. 如下

```
ExecStart=/usr/bin/dockerd-current \
          --add-runtime docker-runc=/usr/libexec/docker/docker-runc-current \
          --default-runtime=docker-runc \
          --exec-opt native.cgroupdriver=systemd \
          --userland-proxy-path=/usr/libexec/docker/docker-proxy-current \
          --graph=/opt/docker \
          $OPTIONS \
...
```

ok, 启动docker服务.

```
$ systemctl start docker
```
