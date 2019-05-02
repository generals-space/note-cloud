# docker代理缓存

```
$ docker run -d \
    -p 443:443 \
    --name registry \
    --restart=always \
    -v /opt/docker/data:/var/lib/registry \
    -v /opt/docker/config.yml:/etc/docker/registry/config.yml \
    -v /opt/docker/certs:/certs \
    registry:2
```

`-p 443:443`与下面config.yml配置文件中`http` > `addr`字段的配置有关, 让registry在容器内监听443端口, 同时将宿主机的443也映射给仓库容器. 如果是http形式的, 也可以写成80对80.

`--restart=always`为docker服务的自动重启机制, 在这个机制下, docker服务本身重启后都会立刻启动镜像仓库容器.

`-v /opt/docker/data:/var/lib/registry`: 挂载了上传的镜像的存储路径, 因为registry镜像会将接收到的镜像放置在容器本身的`/var/lib/registry`目录下, 使用`-v`参数挂载`/opt/docker/data`到容器, 可以将容器信息存储在仓库容器外面, 更方便查看.

`-v /opt/docker/config.yml:/etc/docker/registry/config.yml`: 指定了仓库容器内registry服务的配置文件路径(`/etc/docker/registry/config.yml`), 同样挂载到宿主机的某个目录, 方便修改. 内容如下.

为了实现代理缓存的目的, 通过`proxy`字段配置pull源.

```yml
version: 0.1
log:
    fields:
        service: registry
storage:
    cache:
        blobdescriptor: inmemory
    filesystem:
        ## 镜像存储路径
        rootdirectory: /var/lib/registry
    delete:
        ## 允许删除操作
        enabled: true
http:
    addr: :443
    host: https://registry.sky-mobi.com
    tls:
        ## 加载证书路径, 也可以通过在容器启动时分别通过`REGISTRY_HTTP_TLS_CERTIFICATE`与
        ## `REGISTRY_HTTP_TLS_KEY`两个环境变量指定
        certificate: /certs/registry.sky-mobi.com.crt
        key: /certs/registry.sky-mobi.com.key_nopwd
    headers:
        X-Content-Type-Options: [nosniff]
health:
    storagedriver:
        enabled: true
        interval: 10s
        threshold: 3
## proxy:
##     remoteurl: https://registry-1.docker.io
##     username: generals.space@gmail.com
##     password: 
## proxy:
##     remoteurl: https://so8cv4zt.mirror.aliyuncs.com
## proxy:
##     remoteurl: https://registry.docker-cn.com
```

未实现.

按照官方文档与参考文章7, proxy字段的确是可以实现代理缓存的, 我自己单独搭建docker的代理缓存也是可以的(客户端需要配置docker服务的`--registry-mirror=https://registry.sky-mobi.com`选项). 但是push时会出现`Retrying`的情况. 如下

```
$ docker push registry.sky-mobi.com/k8s-dns-kube-dns-amd64:1.14.4
The push refers to a repository [registry.sky-mobi.com/k8s-dns-kube-dns-amd64]
8963368d3c63: Retrying in 1 second 
404361ced64e: Retrying in 1 second 
unsupportede: Preparing 
```

但是删除proxy后, push功能又可以了, 总之就是私有仓库的推送功能与缓存功能不能同时启用...我尝试了docker官方仓库, 阿里云的私有仓库, 和docker的中国区仓库, 但push依然如故, 参考文章8中提到, 猜测push时也会与中央仓库进行数据校验...放弃了.