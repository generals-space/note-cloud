---
title: kubernetes部署(四)-flannel
tags: [kubernetes, flannel]
categories: general
---

<!--

# kubernetes部署(四)-flannel

<!tags!>: <!kubernetes!> <!flannel!>

<!keys!>: POd18cl?kzijbpgn

-->

参考文章

1. [DockOne技术分享（十八）：一篇文章带你了解Flannel](http://dockone.io/article/618)

2. [浅析flannel与docker结合的机制和原理](http://www.cnblogs.com/xuxinkun/p/5696031.html)

貌似也能用yum装, 但还是那个原因, 自定义性更强, 版本可控, 所以我选二进制安装.

下载[flannel-0.7.1](https://github.com/coreos/flannel/releases/download/v0.7.1/flannel-v0.7.1-linux-amd64.tar.gz). 解压, 里面只有一个可执行文件`flanneld`和一个脚本`mk-docker-opts.sh`.

**所有主机**(Master和Minion)上传`flanneld`和`mk-docker-opts.sh`, 这里指定`/usr/local/flannel/bin`, 并将其添加到环境变量.

```bash
## /usr/lib/systemd/system/flanneld.service 
[Unit]
Description=Flanneld overlay address etcd agent
After=network.target
After=network-online.target
Wants=network-online.target
After=etcd.service
Before=docker.service

[Service]
Type=notify
EnvironmentFile=-/usr/local/flannel/etc/flanneld.conf
EnvironmentFile=-/etc/sysconfig/docker-network
# -etcd-endpoints: etcd地址
# -etcd-prefix: etcd网络配置存储路径
ExecStart=/usr/local/flannel/bin/flanneld \
	-etcd-endpoints=http://172.32.100.71:2379 \
	-etcd-prefix=/sky-mobi.com/network \
    -subnet-file=/run/flannel/subnet.env \
    -log_dir=/var/log \
    -logtostderr

ExecStartPost=/usr/local/flannel/bin/mk-docker-opts.sh -k DOCKER_NETWORK_OPTIONS -d /etc/sysconfig/docker-network -f /run/flannel/subnet.env
Restart=on-failure

[Install]
WantedBy=multi-user.target
RequiredBy=docker.service
```

在上述服务脚本中

`flanneld`启动后会生成`/run/flannel/subnet.env`文件, 里面保存着flannel从etcd分配到的子网信息. 内容如下

```
$ cat /run/flannel/subnet.env 
FLANNEL_NETWORK=172.17.0.0/16
FLANNEL_SUBNET=172.17.38.1/24
FLANNEL_MTU=1472
FLANNEL_IPMASQ=false
```

`ExecStartPost`字段指定的是`flanneld`服务启动完成后执行的操作. 它会执行`mk-docker-opts.sh`脚本, 通过`-f`传入`subnet.env`文件路径, `-k`参数指定docker服务可识别的键名, 脚本把`subnet.env`中的参数转化成docker服务可以识别的格式, 然后输出到`-d`参数所指定的文件中. 其内容为

```
$ cat /etc/sysconfig/docker-network 
DOCKER_OPT_BIP="--bip=172.17.38.1/24"
DOCKER_OPT_IPMASQ="--ip-masq=true"
DOCKER_OPT_MTU="--mtu=1472"
DOCKER_NETWORK_OPTIONS=" --bip=172.17.38.1/24 --ip-masq=true --mtu=1472"
```

其中有效的字段为`DOCKER_NETWORK_OPTIONS`, 正是`mk-docker-opts.sh`通过`-k`指定的字符串. 至于为什么是它, 查看`docker`的服务脚本就知道了.

```ini
$ cat /usr/lib/systemd/system/docker.service 
...省略
[Service]
...省略
EnvironmentFile=-/etc/sysconfig/docker-network
...省略
ExecStart=/usr/bin/dockerd-current \
          --add-runtime docker-runc=/usr/libexec/docker/docker-runc-current \
          --default-runtime=docker-runc \
          --exec-opt native.cgroupdriver=systemd \
          --userland-proxy-path=/usr/libexec/docker/docker-proxy-current \
          $OPTIONS \
          $DOCKER_STORAGE_OPTIONS \
          $DOCKER_NETWORK_OPTIONS \
          $ADD_REGISTRY \
          $BLOCK_REGISTRY \
          $INSECURE_REGISTRY
...省略
```

呶, 就是那个`EnvironmentFile`和`$DOCKER_NETWORK_OPTIONS`的效果了. 当然简单点, 直接写规则文件路径也行, 随便你了.

啊, 还有. 在启动`flanneld`之前, 我们需要在`etcd`中增加一个键, 表示可供kubernetes集群中docker可使用的子网范围. 这里我们添加的是`/sky-mobi.com/network`, 与`flanneld.conf`文件中的`FLANNEL_ETCD_PREFIX`字段一致即可.

```
$ etcdctl --endpoints 'http://172.32.100.71:2379' set /sky-mobi.com/network/config '{ "Network": "172.17.0.0/16" }'
{ "Network": "172.17.0.0/16" }
$ etcdctl --endpoints 'http://172.32.100.71:2379' ls
/sky-mobi.com
```

在各个节点上启动flannel, 然后重启docker.

使用`ip a`查看Master与Minion的`docker0`网络信息, 如果已经变成了`172.17.x.1`, 而且各不相同, 基本上就成了.

再次测试, 在两个Minion节点上各启动一个docker容器, 查看其ip, 然后在容器中互ping, 能ping通就可以了.

记得**关闭防火墙, SELinux, 开启转发**(这个应该是启动flanneld后自动打开的), **桌面版系统最好关闭NetworkManager服务(忠告啊...)**.

```
$ sysctl -a | grep ipv4.ip_forward
net.ipv4.ip_forward = 1
$ systemctl stop NetworkManager
$ systemctl disable NetworkManager
Removed symlink /etc/systemd/system/multi-user.target.wants/NetworkManager.service.
Removed symlink /etc/systemd/system/dbus-org.freedesktop.NetworkManager.service.
Removed symlink /etc/systemd/system/dbus-org.freedesktop.nm-dispatcher.service.
$ setenforce 0
$ systemctl stop firewalld
$ systemctl disable firewalld
Removed symlink /etc/systemd/system/dbus-org.fedoraproject.FirewallD1.service.
Removed symlink /etc/systemd/system/basic.target.wants/firewalld.service.
```