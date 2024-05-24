# dockerd服务无法启动-failed to load listeners：no sockets found via socket activation：make sure the service was started by systemd

参考文章

1. [Installing docker-ce through systemd "fails" due to docker.service not finding socket on first try](https://github.com/docker/for-linux/issues/989)
    - "Change fd:// to unix:// in docker.service file locate in etc/systemd/system"
2. [找不到docker.socket解决方法](https://www.cnblogs.com/flasheryu/p/5802531.html)

## 问题描述

某天电脑重启后, 发现`dockerd`服务启动不了了, 很突然.

```log
$ systemctl restart docker
Job for docker.service failed because the control process exited with error code. See "systemctl status docker.service" and "journalctl -xe" for details.
```

查询`journalctl -xe`以及`/var/log/message`中的日志啥也没发现, 倒是`containerd`的报错日志发现很多, 以为是`containerd`的问题, 改了半天配置结果啥也没搞定.

## 解决方法

后来尝试手动执行`/usr/lib/systemd/system/docker.service`中的`ExecStart`命令, 结果有如下输出

```log
$ /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock -l debug
INFO[2023-05-17T21:02:14.294193828+08:00] Starting up
failed to load listeners: no sockets found via socket activation: make sure the service was started by systemd
```

在百度过程中, 发现了参考文章2, 但是ta没说修改了哪些地方, 只是在对比时偶然发现, 上述`dockerd`的`-H`参数好像是空的, `fd://`啥也没有, 猜测有可能是这个原因. 于是将`-H`参数修改为`-H unix:///var/run/docker.sock`, 然后就可以了...

...天呐, 为啥`ExecStart`命令会突然变掉, 难道是我之前误操作的?😡
