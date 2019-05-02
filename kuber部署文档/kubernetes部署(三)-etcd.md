---
title: kubernetes部署(三)-etcd
tags: [kubernetes, etcd]
categories: general
---

<!--

# kubernetes部署(三)-etcd

<!tags!>: <!kubernetes!> <!etcd!>

<!keys!>: gMg:Foaepcgywv21

-->


参考文章

1. [官方文档](https://github.com/coreos/etcd/releases/tag/v3.2.0)

本文中是用下载etcd二进制文件的方法安装, 也可以使用yum安装, 但都需要配置诸如`ETCD_LISTEN_CLIENT_URLS`等字段, 以便apiserver连接.

下载[etcd-3.2.0](https://github.com/coreos/etcd/releases/download/v3.2.0/etcd-v3.2.0-linux-amd64.tar.gz), 里面已经有可执行文件`etcd`与`etcdctl`, 其中前者作为服务端, 是一个前端进程, 我们需要为其编写服务脚本(其实还是使用yum装过后把配置文件和启动脚本拷贝过来的...咳)

在这之前, 先为其创建配置文件. `/usr/local/etcd/etc/etcd.conf`. (顺便把两个可执行文件放到`/usr/local/etcd/bin`目录下并把这个目录加到环境变量里). 这就随便你怎么顺眼怎么来了, 不过记得服务脚本中`ExecStart`字段的可执行文件路径要修改成你自定义的才行.

主要修改了以下字段

解开[member]的注释

`ETCD_NAME`: 是这个节点的名称

`ETCD_LISTEN_CLIENT_URLS`: 2379用于客户端连接

`ETCD_LISTEN_PEER_URLS`: 2380用于节点间通信

解开cluster的注释(理论上单节点部署不需要这个字段, 但出了个错, )

```
## /usr/local/etcd/etc/etcd.conf
[member]
ETCD_NAME=kubernetes
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
#ETCD_WAL_DIR=""
#ETCD_SNAPSHOT_COUNT="10000"
#ETCD_HEARTBEAT_INTERVAL="100"
#ETCD_ELECTION_TIMEOUT="1000"
ETCD_LISTEN_PEER_URLS="http://172.32.100.71:2380"
ETCD_LISTEN_CLIENT_URLS="http://172.32.100.71:2379"
#ETCD_MAX_SNAPSHOTS="5"
#ETCD_MAX_WALS="5"
#ETCD_CORS=""
[cluster]
#ETCD_INITIAL_ADVERTISE_PEER_URLS="http://localhost:2380"
# if you use different ETCD_NAME (e.g. test), set ETCD_INITIAL_CLUSTER value for this name, i.e. "test=http://..."
#ETCD_INITIAL_CLUSTER="default=http://localhost:2380"
#ETCD_INITIAL_CLUSTER_STATE="new"
#ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_ADVERTISE_CLIENT_URLS="http://172.32.100.71:2379"
#ETCD_DISCOVERY=""
#ETCD_DISCOVERY_SRV=""
#ETCD_DISCOVERY_FALLBACK="proxy"
#ETCD_DISCOVERY_PROXY=""
#ETCD_STRICT_RECONFIG_CHECK="false"
#ETCD_AUTO_COMPACTION_RETENTION="0"
#
#[proxy]
#ETCD_PROXY="off"
#ETCD_PROXY_FAILURE_WAIT="5000"
#ETCD_PROXY_REFRESH_INTERVAL="30000"
#ETCD_PROXY_DIAL_TIMEOUT="1000"
#ETCD_PROXY_WRITE_TIMEOUT="5000"
#ETCD_PROXY_READ_TIMEOUT="0"
#
#[security]
#ETCD_CERT_FILE=""
#ETCD_KEY_FILE=""
#ETCD_CLIENT_CERT_AUTH="false"
#ETCD_TRUSTED_CA_FILE=""
#ETCD_AUTO_TLS="false"
#ETCD_PEER_CERT_FILE=""
#ETCD_PEER_KEY_FILE=""
#ETCD_PEER_CLIENT_CERT_AUTH="false"
#ETCD_PEER_TRUSTED_CA_FILE=""
#ETCD_PEER_AUTO_TLS="false"
#
[logging]
ETCD_DEBUG="true"
# examples for -log-package-levels etcdserver=WARNING,security=DEBUG
#ETCD_LOG_PACKAGE_LEVELS=""
```


然后创建服务脚本`/usr/lib/systemd/system/etcd.service`.

```
## /usr/lib/systemd/system/etcd.service 
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
WorkingDirectory=/var/lib/etcd/
EnvironmentFile=-/usr/local/etcd/etc/etcd.conf
## etcd服务的启动用户, 需要事先创建etcd用户
User=etcd
# set GOMAXPROCS to number of processors
ExecStart=/bin/bash -c "GOMAXPROCS=$(nproc) /usr/local/etcd/bin/etcd --name=\"${ETCD_NAME}\" --data-dir=\"${ETCD_DATA_DIR}\" --listen-client-urls=\"${ETCD_LISTEN_CLIENT_URLS}\""
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

哦, 它还需要建一个etcd用户, 自行创建吧. `useradd etcd -s /sbin/nologin -m -d /var/lib/etcd`.

其中`-d`参考指定的etcd家目录正好是`etcd`通过`--data-dir`指定的目录, 所以是必需的.

启动.

```
$ systemctl start etcd
$ netstat -nlp | grep 2379
tcp        0      0 172.32.100.71:2379      0.0.0.0:*               LISTEN      16794/etcd
```