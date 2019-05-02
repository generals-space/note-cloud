# CentOS7搭建Docker私有仓库

参考文章

[CentOS7搭建Docker私有仓库](http://www.centoscn.com/CentosServer/ftp/2015/0426/5280.html)

## 1. 环境准备

- Docker私有仓库服务端: A

- Docker客户端: B

> 注意关掉两者的防火墙与SELinux

## 2. 服务端A

### 2.1 安装仓库镜像

Docker私有仓库服务端需要首先安装`docker`, 然后拉取docker官方提供的`registry`镜像(即仓库镜像).

```shell
yum install -y docker
docker pull docker.io/registry
```

### 2.2 启动仓库容器

```shell
docker run -d --name docker-registry -p 5000:5000 -v /opt/docker/registry:/tmp/registry docker.io/registry
```

- `-p`是端口映射, 访问A上的5000端口就相当于访问仓库容器的5000端口.

- `-v`挂载共享分区, 这里注意, `/opt/registry`为仓库容器宿主机上的路径, 可以随意指定; `/tmp/registry`是仓库容器内部的路径, **不要更改**, 因为之后`docker push`的镜像都将保存在这个地方, 这样就可以不进入容器也能看到上传的镜像了.

### 2.3 测试查询与上传

#### 2.3.1 查询

```
curl 127.0.0.1:5000/v1/search
{"num_results": 0, "query": "", "results": []}
```

私有仓库结果为空, 因为没有提交新镜像到仓库中.

#### 2.3.2 上传

从官网上下载CentOS6的镜像用来实验(也可以是本地的任何镜像)

```
docker pull docker.io/centos:6
```

给它打个标签, 注意这个标签名, 也就是`127.0.0.1:5000/centos6:base`是有意义的, 尤其是`127.0.0.1:5000`, 这其实是镜像所属的仓库地址, 就类似与那个`docker.io`. 镜像名称与标签名就随便了, 你开心就好.

```
docker tag centos6的镜像ID 127.0.0.1:5000/centos6:base
```

然后修改Docker配置文件, Docker新版需要 `SSL Auth`, 解决就是两种方法,一种就是在私有仓库上使用SSL, 需要安装签发证书,另一种就是强制使用普通http方式...签发证书那么麻烦, 果断放弃. 注意重启docker服务.

```
vim /etc/sysconfig/docker
...
## 此句解开注释, 并修改为如下
OPTIONS='--selinux-enabled --insecure-registry=127.0.0.1:5000'
...

## 注意重启docker服务
systemctl restart docker
```

下面尝试上传...


```
docker push 127.0.0.1:5000/centos6
The push refers to a repository [127.0.0.1:5000/centos6] (len: 1)
Sending image list
Pushing repository 127.0.0.1:5000/centos6 (1 tags)
511136ea3c5a: Image successfully pushed
...
2aeb2b6d9705: Image successfully pushed
Pushing tag for rev [2aeb2b6d9705] on {http://127.0.0.1:5000/v1/repositories/centos6/tags/base}
```

再次查询

```
curl 127.0.0.1:5000/v1/search
{"num_results": 1, "query": "", "results": [{"description": "", "name": "library/centos6"}]}
```

在宿主机上`/opt/Registry`中查看存储文件.

```shell
$ tree /opt/registry/repositories/
/opt/registry/repositories/
└── library
    └── centos6
        ├── _index_images
        ├── json
        ├── tag_latest
        └── taglatest_json

2 directories, 4 files
```

## 3. 客户端

想要将客户端的镜像上传到私有仓库, 当然也是需要安装docker的. 这里不多描述了. 要做的只有修改docker的配置文件.

```
vim /etc/sysconfig/docker
...
## 此句解开注释, 并修改为如下
OPTIONS='--selinux-enabled --insecure-registry=私有仓库的IP:5000'
...

## 注意重启docker服务
systemctl restart docker
```

然后就能`pull`与`push`了.

------

## 扩展

### 1.

```
docker push 127.0.0.1:5000/centos6:base
The push refers to a repository [127.0.0.1:5000/centos6]
Get https://127.0.0.1:5000/v1/_ping: EOF
```

```
Error: Invalid registry endpoint https://127.0.0.1:5000/v1/: Get https://127.0.0.1:5000/v1/_ping: EOF. If this private registry supports only HTTP or HTTPS with an unknown CA certificate, please add `--insecure-registry 127.0.0.1:5000` to the daemon's arguments. In the case of HTTPS, if you have access to the registry's CA certificate, no need for the flag; simply place the CA certificate at /etc/docker/certs.d/127.0.0.1:5000/ca.crt
```

几乎所有私有仓库的pull与push的问题都是由于https认证方式的问题. 除去这个, 私有仓库的搭建还是很简单的. 这些问题按照上面的方法可以解决.

### 2.

有一种情况, 当docker并不是通过yum或apt-get方式安装的时候. 比如, 我是按照DaoCloud安装的官方最新版docker. 搭建私有仓库时发现`/etc/sysconfig`目录下没有`docker`文件, 这就很尴尬了. 最初我自行创建`docker`文件, 只添加了`OPTIONS='--selinux-enabled --insecure-registry=私有仓库的IP:5000'`一行. 重启之后并没有效果.

按照[这篇文章](https://forums.docker.com/t/docker-private-registry-ping-attempt-failed/4868/3)中`nikx`的提示, 修改docker的启动脚本, 具体如下.

```
vim /usr/lib/systemd/system/docker.service

# ExecStart=/usr/bin/docker daemon
ExecStart=/usr/bin/docker daemon -H tcp://127.0.0.1:2375 -H unix:///var/run/docker.sock --insecure-registry 私有仓库的IP:5000
```

**注意私有仓库的客户端与服务端都要进行此项修改**.

------

**更新**

通过yum安装的docker, 其启动脚本类似如下

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

而通过daocloud安装的docker, 启动脚本为

```
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
ExecStart=/usr/bin/docker daemon
```

前者的诸如`$OPTIONS`, `$INSECURE_REGISTRY`是通过`/etc/sysconfig/docker`文件中`EnvironmentFile`字段定义的, 但后者的service文件不支持这种方式, 所以只能把对应的选项直接写到`ExecStart`字段中.

另外, docker服务支持的选项参数及其格式可以通过如下命令查看

```
/usr/bin/docker daemon --help
```

### 3. 如何删除私有仓库的镜像?

参考文章

[如何删除私有仓库的镜像](http://www.aixchina.net/Question/123917)

```
curl -X DELETE localhost:5000/v1/repositories/ubuntu/tags/latest
```

ubuntu是镜像名, latest是版本名

### 4. docker仓库容器无法启动

系统版本: CentOS7

docker版本: 1.10.3

使用`docker logs 容器ID`得到如下输出.

```
docker logs 容器ID
...
Traceback (most recent call last):
  File "/usr/local/lib/python2.7/dist-packages/gunicorn/arbiter.py", line 507, in spawn_worker
    worker.init_process()
...
  File "/usr/local/lib/python2.7/dist-packages/docker_registry/toolkit.py", line 330, in wrapper
    os.remove(lock_path)
OSError: [Errno 2] No such file or directory: './registry._setup_database.lock'
[2016-07-23 13:49:03 +0000] [13] [INFO] Worker exiting (pid: 13)
[2016-07-23 13:49:03 +0000] [17] [INFO] Worker exiting (pid: 17)
[2016-07-23 13:49:04 +0000] [12] [INFO] Worker exiting (pid: 12)
[2016-07-23 13:49:04 +0000] [16] [INFO] Worker exiting (pid: 16)
[2016-07-23 13:49:04 +0000] [1] [INFO] Shutting down: Master
[2016-07-23 13:49:04 +0000] [1] [INFO] Reason: Worker failed to boot.
```

比较容易识别的是`OSError`这一行. 按照[issue #796](https://github.com/docker/docker-registry/issues/796)所说, 启动registry容器时, 在命令行加入`-e GUNICORN_OPTS=[--preload]`可以解决.

```
docker run -d --name docker-registry -p 5000:5000 -e GUNICORN_OPTS=[--preload]  -v /opt/docker/registry:/tmp/registry docker.io/registry
```

### 5. docker服务

docker服务支持的选项参数及其格式可以通过如下命令查看

```
/usr/bin/docker daemon --help
```

### 6. 添加多个私有仓库(待检验, 可能有问题)

查看添加信任仓库的语法

```
docker daemon --help | grep registry
  --disable-legacy-registry                Do not contact legacy registries
  --insecure-registry=[]                   Enable insecure registry communication
  --registry-mirror=[]                     Preferred Docker registry mirror
```

可以看到`--insecure-registry`选项的取值类型为数组, 添加多个信任仓库, 只需要使用如下方式即可.

```
--insecure-registry=[sky.generals.space:5000,docker.generals.space:5000]
```
