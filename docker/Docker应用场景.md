# Docker应用配置场景

## 1. 修改docker镜像存储位置

docker默认将存储目录设置在`/var/lib/docker`, 如果`/usr`分区所在目录空间不足, 或其他什么原因, 希望将此目录设置为其他(示例中假定为`/opt/docker/data`)可以使用如下方法.

### 1.1 建立软链接

```
## 首先停止docker服务, 不再对原目录进行写入
$ systemctl stop docker
## 执行目录迁移
$ mkdir -p /opt/docker
$ mv /var/lib/docker /opt/docker/data
$ ln -s /opt/docker/data /var/lib/docker
## 重启docker
$ systemctl start docker
```

### 1.2 修改docker服务启动参数

涉及到docker服务启动项的有两个文件

#### 1.2.1 `/usr/lib/systemd/system/docker.service`: docker的服务脚本.

```
...
ExecStart=/bin/sh -c '/usr/bin/docker-current daemon \
          --exec-opt native.cgroupdriver=systemd \
          $OPTIONS \
          $DOCKER_STORAGE_OPTIONS \
          $DOCKER_NETWORK_OPTIONS \
          $ADD_REGISTRY \
          $BLOCK_REGISTRY \
          $INSECURE_REGISTRY \
          2>&1 | /usr/bin/forward-journald -tag docker'
...
```

#### 1.2.2 `/etc/sysconfig/docker`: docker启动项配置文件

```
OPTIONS='--selinux-enabled --log-driver=journald'
DOCKER_CERT_PATH=/etc/docker
# ADD_REGISTRY='--add-registry registry.access.redhat.com'

# BLOCK_REGISTRY='--block-registry'

INSECURE_REGISTRY='--insecure-registry'

DOCKER_TMPDIR=/opt/docker/data/tmp

# LOGROTATE=false
```

可以看到, `docker.service`文件中docker的启动项中引用了如`$OPTIONS`, `$ADD_REGISTRY`等变量, 而这些变量的值就是在`/etc/sysconfig/docker`中定义的.

我们在后者的`$OPTIONS`中加入`-g=/opt/docker/data`, 就可以使docker的存储路径变为`/opt/docker/data`.

```
OPTIONS='--selinux-enabled --log-driver=journald -g=/opt/docker/data'
```

然后启动docker服务, 查看当前存储根目录, 是否已经写入数据.
