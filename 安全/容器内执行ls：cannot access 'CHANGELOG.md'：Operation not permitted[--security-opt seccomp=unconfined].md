# 容器内执行ls：cannot access 'CHANGELOG.md'：Operation not permitted

参考文章

1. [Docker + SLES15 : Unable to access files inside docker container : ls: cannot access '<directory/file name>': Operation not permitted](https://stackoverflow.com/questions/70714357/docker-sles15-unable-to-access-files-inside-docker-container-ls-cannot-ac)
    - 容器镜像与宿主机不是同一个系统.
2. [docker container出現Operation not permitted的錯誤](https://blog.twtnn.com/2021/09/docker-containeroperation-not-permitted.html)

## 问题描述

打完镜像后启动容器测试一下, 结果发现`ls`命令执行报错.

```log
[root@8276303295ba logstash]# ls
ls: cannot access 'CHANGELOG.md': Operation not permitted
## 省略
ls: cannot access 'vendor': Operation not permitted
CHANGELOG.md  Gemfile                 LICENSE     bin     data  logstash-core             logstash_start.sh  tools
CONTRIBUTORS  Gemfile.jruby-1.9.lock  NOTICE.TXT  config  lib   logstash-core-plugin-api  modules            vendor
[root@8276303295ba logstash]# ls -al
ls: cannot access '.': Operation not permitted
ls: cannot access '..': Operation not permitted
ls: cannot access 'CHANGELOG.md': Operation not permitted
## 省略
ls: cannot access 'vendor': Operation not permitted
total 0
d????????? ? ? ? ?            ? .
d????????? ? ? ? ?            ? ..
-????????? ? ? ? ?            ? CHANGELOG.md
## 省略
d????????? ? ? ? ?            ? vendor
```

不只是logstash目录, 其他所有目录都会出现, 而且我是root用户, 所以不是权限的问题.

虽然有报错, 但是文件/目录还是打印出来了.

## 解决方法

虽然我的确如参考文章1中所说, 容器的基础镜像与宿主机并不是同一个系统(宿主机为 centos7, 容器基础镜像为 openeuler:20.03, 都是arm64架构的), 但是也没有要求说 centos 主机不能运行 ubuntu 容器啊.

而且之前也没有出现过这种情况, 而且目前测试环境没有 arm 架构的 openeuler 主机, 所以无法测试.

后来又找到了参考文章2, 该文章提到了`seccomp`, 即使是 root 用户也会受到ta的限制.

我的宿主机系统的确开启了`secocmp`.

```log
$ cat /boot/config-$(uname -r) |grep CONFIG_SECCOMP
CONFIG_SECCOMP=y
```

参考文章2给出了2种解决方法(第3种貌似跟这个问题不相关)

1. `--privileged`开启特权模式
2. `--security-opt seccomp=unconfined`关闭`seccomp`特性

实践第2种方法有效.
