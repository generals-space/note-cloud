# docker history查看镜像构建过程[dockerfile]

参考文章

1. [查看docker镜像的修改历史记录](https://www.jianshu.com/p/9d3d16691329)
    - 从上到下时间依次久远
2. [通过docker history查看镜像构建过程（即dockerfile）](https://www.cnblogs.com/cooper-73/p/9830371.html)
    - `--no-trunc`和`--format`参数的使用.
3. [Dockerfile的最佳实践](https://juejin.im/post/6844903922830671885)
    - 只有`RUN`,`COPY`,`ADD`指令会创建层, 其他指令创建**临时中间层**, 并不增加构建的大小.
4. [“docker image history” shows <missing> on image name](https://forums.docker.com/t/docker-image-history-shows-missing-on-image-name/33948)
    - 引用了参考文章5, 解释了镜像信息丢失显示`<missing>`的原因
5. [Explaining Docker Image IDs](https://windsock.io/explaining-docker-image-ids/)

```
docker history 镜像ID
docker image history 镜像ID
```

以下以`registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:devops`镜像为例

```
[root@k8s-master-01 manifests]# docker history registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:devops --no-trunc
IMAGE                                                                     CREATED             CREATED BY                                                                                                                                                                                                            SIZE                COMMENT
sha256:e5d8f10fcec2116498fa2c1b40b463f140c681a89fa60c512578ad47b53b0f61   7 months ago        /bin/sh -c #(nop)  CMD ["tail" "-f" "/etc/profile"]                                                                                                                                                                   0B
<missing>                                                                 7 months ago        /bin/sh -c yum install -y iproute net-tools     && yum clean all     && rm -rf /var/cache/yum                                                                                                                         28.1MB
<missing>                                                                 7 months ago        /bin/sh -c #(nop)  ENV LANG=C.UTF-8                                                                                                                                                                                   0B
<missing>                                                                 7 months ago        /bin/sh -c #(nop)  LABEL email=generals.space@gmail.com                                                                                                                                                               0B
<missing>                                                                 7 months ago        /bin/sh -c #(nop)  LABEL author=general                                                                                                                                                                               0B
<missing>                                                                 7 months ago        /bin/sh -c #(nop)  CMD ["tail" "-f" "/etc/profile"]                                                                                                                                                                   0B
<missing>                                                                 7 months ago        /bin/sh -c \cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime                                                                                                                                                    528B
<missing>                                                                 7 months ago        /bin/sh -c yum update -y     && yum install -y gcc gcc-c++ glibc-common make vim telnet     && yum clean all     && rm -rf /var/cache/yum                                                                             214MB
<missing>                                                                 7 months ago        /bin/sh -c rm -f /etc/yum.repos.d/*     && curl http://mirrors.aliyun.com/repo/Centos-7.repo -o /etc/yum.repos.d/Centos-7.repo     && curl http://mirrors.aliyun.com/repo/epel-7.repo -o /etc/yum.repos.d/epel.repo   3.19kB
<missing>                                                                 7 months ago        /bin/sh -c #(nop)  ENV LANG=C.UTF-8                                                                                                                                                                                   0B
<missing>                                                                 7 months ago        /bin/sh -c #(nop)  LABEL email=generals.space@gmail.com                                                                                                                                                               0B
<missing>                                                                 7 months ago        /bin/sh -c #(nop)  LABEL author=general                                                                                                                                                                               0B
<missing>                                                                 11 months ago       /bin/sh -c #(nop)  CMD ["/bin/bash"]                                                                                                                                                                                  0B
<missing>                                                                 11 months ago       /bin/sh -c #(nop)  LABEL org.label-schema.schema-version=1.0 org.label-schema.name=CentOS Base Image org.label-schema.vendor=CentOS org.label-schema.license=GPLv2 org.label-schema.build-date=20191001               0B
<missing>                                                                 11 months ago       /bin/sh -c #(nop) ADD file:45a381049c52b5664e5e911dead277b25fadbae689c0bb35be3c42dff0f2dffe in /                                                                                                                      203MB
[root@k8s-master-01 manifests]# docker images | grep devops
registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7                      devops              e5d8f10fcec2        7 months ago        445MB
```

首先, 第一列为构建时其各层镜像的id, `e5d8f10fcec2`即为该镜像本身的id.

下面的行中, 第一列为`<missing>`, 表示其依赖的基础镜像不在这台主机上了. 实际上, 就算当前主机上其依赖的镜像, 也可能显示`<missing>`, 因为这些信息显示的是构建过程中每个中间层, 上传到`registry`的镜像信息中并不包含这些信息. 所以当这些构建中间层被清理过后, 就算存在依赖的基础镜像也不会再显示了.

`CREATE BY`这一列, 表示的就是`dockerfile`中每一行的操作了, 仔细分辨就会发现, 除了`RUN`指令, 其他如`COPY`, `LABEL`的行都有一个`#(nop)`的标记.

后面的`SIZE`表示当前层的变动大小, 比如`yum`安装的层数值就比较大, 而`ENV`, `LABEL`操作的层都是`0B`.

最底层的`ADD file:xxx in /`, 应该是最初的依赖镜像, 因为之后我依赖自建的镜像进行构建时, 没有再出现过这样的情况.
