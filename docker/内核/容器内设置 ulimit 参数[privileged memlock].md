# 容器内设置 ulimit 内核参数[privileged]

参考文章

1. [rlimit support](https://github.com/kubernetes/kubernetes/issues/3595)
2. [给容器设置内核参数](https://tencentcloudcontainerteam.github.io/2018/11/19/kernel-parameters-and-container/)
3. [配置 Kubernetes 中 Docker 容器的 ulimit](https://blog.csdn.net/yellowfruit/article/details/108479895)
    - coredump 文件过大, 严重影响系统性能, 所以需要禁止进程在容器内生成 core 文件, 消除对整个系统的影响。
    - docker 服务通过`systemd`托管, `ulimit`的参数设置对`systemd`托管的服务是无效的, 因此不能直接修改
4. [elasticsearch 官网](https://www.elastic.co/guide/en/elasticsearch/reference/7.3/setting-system-settings.html#limits.conf)
5. [容器内的ulimits](https://zhuanlan.zhihu.com/p/144230003)
    - `CAP_SYS_RESOURCE`: `ulimit`所需的内核能力
6. [capabilities(7) — Linux manual page](https://man7.org/linux/man-pages/man7/capabilities.7.html)
    - `capabilities`参考手册.

参考文章2中给出了使用`systemd`给`dockerd`进程本身设置`ulimit`参数, 以及通过`dockerd`选项为启动的容器设置默认的`ulimit`参数的方法, 不过我们需要的不是这个.

同时ta也给出了两个设置容器内`ulimit`参数的方法, 如下.

## 1. `docker run`的`--ulimit`选项

```
docker run -d --ulimit nofile=20480:40960 nproc=1024:2048 容器名
```

> 冒号前面是 soft limit, 后面是 hard limit.
> 
> limit 值只接受数值类型, `unlimited`需要设置成`-1`.

## 2. kuber 中设置`ulimit`

## 3. memlock(max locked memory), `ulimit`命令

在一个普通容器里执行`ulimit`是可行的, 就像在参考文章3中所说的那样.

```
[root@2ab83fdb2871 /]# ulimit -a
max locked memory       (kbytes, -l) 82000
open files                      (-n) 1048576
[root@2ab83fdb2871 /]# ulimit -n 65535
[root@2ab83fdb2871 /]# ulimit -a
max locked memory       (kbytes, -l) 82000
open files                      (-n) 65535
```

`elasticsearch`配置文件中有一个`bootstrap.memory_lock`选项, 默认为`false`, 改为`true`开启内存锁定可以提高性能. 需要设置`max locked memory`值, 这个值表示可以锁定的内存的大小.

但是在容器内修改`max locked memory`时, 报了如下错误

```
[root@2ab83fdb2871 /]# ulimit -l unlimited
bash: ulimit: max locked memory: cannot modify limit: Operation not permitted
```

看起来像是`core file size`, `open files`这些比较...不那么危险? 但是`max locked memory`就比较危险了.

> 虽然设置`unlimited`失败了, 但其实当设置值大于默认的`82000`时都会失败, 设置值小于`82000`是可以成功的...

------

在`docker run`时使用`--ulimit`选项可以

```console
$ docker run -it --name ulimit --ulimit memlock=-1:-1 registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7-devops bash
[root@3b5fe1388e5a /]# ulimit -a | grep locked
max locked memory       (kbytes, -l) unlimited
```

不过这种方法貌似还有点问题.

```
$ docker run -it --name ulimit --ulimit memlock=100:100 registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7-devops bash
[root@78f42b22f00b /]# ulimit -a | grep locked
max locked memory       (kbytes, -l) 0
```

设置其他值时, 直接为0了, 貌似以`82000`这个值为分界线...???

> 当然, 先使用`--privileged`启动容器, 再执行`ulimit -l unlimited`也是可以的...

------

虽然说`ulimit`命令设置的参数只在当前会话中有效, 但是在此会话下的子进程也可以生效.

```console
$ docker run -it --name ulimit --privileged registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7-devops bash
[root@a4d27f9b6c45 /]# ulimit -l unlimited
[root@a4d27f9b6c45 /]# ulimit -a
max locked memory       (kbytes, -l) unlimited
[root@a4d27f9b6c45 /]# bash
[root@a4d27f9b6c45 /]# ulimit -a
max locked memory       (kbytes, -l) unlimited
[root@a4d27f9b6c45 /]# useradd general
[root@a4d27f9b6c45 /]# su -l general
[general@a4d27f9b6c45 ~]$ ulimit -a
max locked memory       (kbytes, -l) unlimited
```

## 4. memlock(max locked memory), `/etc/security/limit.conf`文件

启动一个`--privileged`的容器, 修改`/etc/security/limit.conf`文件, 添加如下内容

```
* soft memlock unlimited
* hard memlock unlimited
```

我之前一篇笔记说修改这个文件是即时生效的, 不需要重启. 但是在容器里修改后就没什么反应, 使用`ulimit -a`查看, 还是原来的样子.

又由于在容器中`elasticsearch`需要使用一个普通用户启动, 所以我又建了个用户, 竟然可以.

```
[root@e598f55a2320 security]# ulimit -a
max locked memory       (kbytes, -l) 82000
[root@e598f55a2320 security]# useradd general
[root@e598f55a2320 security]# su -l general
[general@e598f55a2320 ~]$ ulimit -a
max locked memory       (kbytes, -l) unlimited
```

> **别想新建一个bash再看了, `/etc/security/limit.conf`这方法对root就不生效, 重启容器也一样, 所以容器内的主进程最好通过普通用户运行.**

------

不使用`--privileged`, 修改`/etc/security/limit.conf`文件呢?

```console
[root@01bd4c7e71d8 security]# useradd general
[root@01bd4c7e71d8 security]# su -l general
su: cannot open session: Permission denied
```

别说`root`了, 连普通用户也切换不了了.

## 5. CAP_SYS_RESOURCE

如果实在需要修改容器内的`memlock`, 也最好不要直接指定`privileged`, 而是通过`capabilities`指定内核能力.

```console
$ docker run -it --name ulimit --cap-add SYS_RESOURCE registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7-devops bash
[root@a282fae2b94f /]# ulimit -l unlimited
[root@a282fae2b94f /]# ulimit -a | grep locked
max locked memory       (kbytes, -l) unlimited
```

```yaml
        securityContext:
          privileged: false
          capabilities:
            add: ["CAP_SYS_RESOURCE"]
```

经过验证, 正常使用.
