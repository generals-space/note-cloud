# Docker私有仓库搭建(二)

参考文章

1. [官方文档 - Docker Registry](https://docs.docker.com/registry/)

2. [Registry镜像的Dockerfile](https://github.com/docker/distribution-library-image/blob/master/Dockerfile)

3. [官方文档 - 详细配置](https://docs.docker.com/registry/configuration/#list-of-configuration-options)

4. [centos7 docker1.12安装私有仓库](http://blog.csdn.net/kadiya2011/article/details/443344337)

5. [搭建Docker私有仓库--自签名方式](http://www.cnblogs.com/li-peng/p/6511331.html)

6. [Registry as a pull through cache](https://docs.docker.com/registry/recipes/mirror/)

7. [如何使用Docker开源仓库建立代理缓存仓库](http://dockone.io/article/774?utm_source=tuicool&utm_medium=referral)

8. [docker registry v2 配置 (REGISTRY_PROXY_REMOTEURL) 解释 + docker pull/push 动作简析](http://www.open-open.com/lib/view/open1456893590796.html)

9. [Docker私有仓库Registry V2 搭建(附带CA证书自签名)](https://my.oschina.net/yyflyons/blog/656278)

10. [搭建Docker私有仓库Registry-v2](http://www.tuicool.com/articles/6jEJZj)

11. [Docker认证自签名证书-Using self-signed certificates](https://docs.docker.com/registry/insecure/#using-self-signed-certificates)

12. [Verify repository client with certificates](https://docs.docker.com/engine/security/certificates/)

## 1. registry认识

`registry`为docker官方提供的搭建私有镜像库的解决方案. 它是一个docker镜像, 基于`alpine`镜像(`alpine`是一个轻量级的linux系统, 一般用于路由器, 防火墙等). `registry`容器启动后, 只有`registry`作为服务被执行.

```
1c5a7ccde40f:/bin# ps -ef
PID   USER     TIME   COMMAND
    1 root       5:35 registry serve /etc/docker/registry/config.yml
   50 root       0:00 -ash
   65 root       0:00 ps -ef
```

`registry`可执行文件在`/bin/registry`, 可以随便启一个registry镜像进行看看. 它本质上是一个web服务器, 与python的flask, nodejs的express, java的tomcat类似, 镜像的查询上传和下载操作都通过http协议完成, 于是可信任仓库配置中需要指定https证书也变得容易理解了. 至于镜像缓存和健康检查, 哼, 不就和nginx一样么?

> ok, 上面的都是我猜的, 各位自己理解...

我曾经尝试过把这个文件拷贝出来, 直接在CentOS7宿主机上运行, 但是得到了如下结果.

```
$ ./registry --help
-bash: ./registry: /lib/ld-musl-x86_64.so.1: bad ELF interpreter: No such file or directory
```

看来是因为编译链接和运行环境不同, 不知道跨平台运行这个坑踩下去有多深, 反正我不踩. 我不明白为什么要把`registry`放到docker镜像里, 还是`alpine`这种小众镜像, 直接拿来部署为仓库服务器不好吗? 我感觉了来自docker官方深深的恶意.

## 2. insecure/纯http仓库

仓库地址: 192.168.166.220

首先要pull官方的registry镜像...

```
$ docker pull registry:latest
```

版本的话自行到docker hub官方查看, 截止2017-07-05貌似最高版本是2.6.

然后启动...

好吧先别启动(如果你只是先玩玩就随便你了)

```
$ docker run -d -p 5000:5000 --name registry registry:latest
```

```
$ docker images
REPOSITORY                                             TAG                 IMAGE ID            CREATED             SIZE
daocloud.io/nginx                                      latest              3448f27c273f        7 weeks ago         109.4 MB
$ docker tag 3448f27c273f 192.168.166.220:5000/nginx:latest
$ docker push 192.168.166.220:5000/nginx:latest
The push refers to a repository [192.168.166.220:5000/nginx]
Get https://192.168.166.220:5000/v1/_ping: net/http: TLS handshake timeout
```

没有证书, 怎么破?

你需要在客户端的docker服务启动时添加`--insecure-registry=192.168.166.220:5000`这个选项(或者不加等号, 直接写`--insecure-registry 192.168.166.220:5000`也可以的), 不管你是直接在docker的systemd启动脚本中直接添加这个选项, 还是在`/etc/sysconfig/docker`文件中解开`INSECURE_REGISTRY`字段的注释, 只要在ps时能看到dockerd的启动参数中有`--insecure-registry`这个选项就行.

![](https://gitee.com/generals-space/gitimg/raw/master/a744cd1fc262339d069cc3d96e20cd59.png)

然后再执行`push`操作就可以了.

> `--insecure-registry`可以写多个哦, 不过不是通过逗号分隔的, 而是写成`--insecure-registry=192.168.166.220:5000 --insecure-registry=192.168.166.221:5000`这种(同样, 不加等号的`--insecure-registry 192.168.166.220:5000 --insecure-registry 192.168.166.221:5000`, 亲测可行哦).

这只是上传下载操作, 搜索操作还是不会去寻找私有仓库的地址的. docker提供了另外一个选项`--add-registry`, 通过指定它为自己的私有仓库地址可以优先搜索私有镜像库中的镜像. 但是经实验, 这个选项是无效的. 官方文档中也没有相关的解释. 反而是Red Hat对此有一篇文章[Docker Experimental Features in Red Hat Enterprise Linux](https://access.redhat.com/articles/1354823)

没有找到对此合适的解决方案, 取代的是使用curl通过restful API查询.

```
$ curl 192.168.166.220:5000/v2/_catalog
{"repositories":["nginx"]}
```

这只是镜像种类, 不包括标签信息, 需要再次手动查询.

```
$ curl 192.168.182.3:5000/v2/nginx/tags/list
{"name":"nginx","tags":["latest","1.0"]}
```

还算可以接受.

## 3. 使用自签名证书

接下来是证书的配置, 需要销毁重建一个容器.

docker的https证书同浏览器访问https网站是同一个道理, 需要由第三方机构(沃通, startSSL等)签发的合法证书才行, 并且你需要有一个顶级域名. 当然你可以申请一个1-2年的免费证书, 但也需要有一个顶级域名才行, docker官方文档中提到一个[letsencrypt](https://docs.docker.com/registry/configuration/#letsencrypt)方案, 我瞄了一眼就知道那很麻烦, 而且我英文不好, 看不太懂. 于是找到了参考文章5.

本来嘛, 私有仓库都是自己人用的, 自己人用当然信得过了. 我自己生成一个证书, 让客户端都添加上信任不就好了?

### 3.1 生成自签名证书

参考文章11中的docker官方文档中有生成自签名证书的方法. 这里随便生成一个证书, 以`registry.sky-mobi.com`为镜像域名(因为我准备直接用443当作服务端口, 所以没加端口好, 它其实等同于`www.example.com:5000`这种形式).

```
$ openssl req -newkey rsa:4096 -nodes -sha256 -keyout registry.sky-mobi.com.key -x509 -days 365 -out registry.sky-mobi.com.crt
Generating a 4096 bit RSA private key
............................................++
..........................................................................................................................++
writing new private key to 'registry.sky-mobi.com.key'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:CN                                                ## 国家(随便写)
State or Province Name (full name) []:ZheJiang                                      ## 省份(随便写)
Locality Name (eg, city) [Default City]:Hangzhou                                    ## 城市(随便写)
Organization Name (eg, company) [Default Company Ltd]:sky-mobi                      ## 公司(随便写)
Organizational Unit Name (eg, section) []:IT                                        ## 部门(随便写)
Common Name (eg, your name or your server's hostname) []:registry.sky-mobi.com      ## 目标域名(必须与仓库地址一致)
Email Address []:domain.sky-mobi.com                                                ## 联系邮箱(随便写)
$ ls
registry.sky-mobi.com.crt  registry.sky-mobi.com.key 
```

我们先试试这两个证书的作用, 把这两个文件放在`/tmp/certs`目录中, 通过`-v`参数挂载到要启动的容器中.

然后通过`REGISTRY_HTTP_TLS_CERTIFICATE`与`REGISTRY_HTTP_TLS_KEY`两个环境变量加载它们.

```
$ docker run -d -p 5000:5000 \
    --name registry \
    --restart=always \      
    -v /tmp/certs:/certs \
    -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/registry.sky-mobi.com.crt \
    -e REGISTRY_HTTP_TLS_KEY=/certs/registry.sky-mobi.com.key_nopwd \
    registry:2
```

然后在docker客户端主机hosts文件中写入`registry.sky-mobi.com`, 试着push一下.

```
$ docker push registry.sky-mobi.com:5000/nginx:1.0.0
The push refers to a repository [registry.sky-mobi.com:5000/nginx]
Get https://registry.sky-mobi.com:5000/v1/_ping: x509: certificate is valid for sky-mobi.com, not registry.sky-mobi.com
```

此时curl查询操作也是这样的结果.

```
$ curl https://registry.sky-mobi.com:5000/v2/_catalog
curl: (60) Peer's certificate issuer has been marked as not trusted by the user.
More details here: http://curl.haxx.se/docs/sslcerts.html
...
```

没关系, 有心理准备.

### 3.2 添加信息

因为这个证书不是权威机构生成的, 所以我们需要让客户端信任这个证书. 按照参考文章9和参考文章10, 把我们生成的`.crt`文件直接拷贝到客户端的`/etc/docker/certs.d/你的仓库地址[:端口]/`目录下就行了(名称随便哦, 没关系), 不用重启docker.

```
$ pwd
/etc/docker/certs.d
$ ls
redhat.com  redhat.io  registry.sky-mobi.com:5000
$ cd registry.sky-mobi.com:5000/
$ ls
registry.sky-mobi.com.crt
$ docker push registry.sky-mobi.com:5000/nginx:1.0.0
The push refers to a repository [sky-mobi.com:5000/nginx]
8963368d3c63: Pushed 
404361ced64e: Pushed 
1.0.0: digest: sha256:38fe56c7f45f7007e6d06ed473c015760bf24d5f048e430af8ce07e75029d630 size: 740
```

这样, `docker push`操作就可以直接进行了. 但是要注意, 这只是针对docker服务本身的信任, 此时使用curl通过restful API查询依然会有问题.

```
$ curl https://registry.sky-mobi.com/v2/_catalog
curl: (60) Peer's certificate issuer has been marked as not trusted by the user.
More details here: http://curl.haxx.se/docs/sslcerts.html
...
```

有什么办法让我们签发的证书在系统层面生效呢?

根据参考文章5的提示, 其实也正是参考文章11中官方的说法, 把我们生成的后缀为`.crt`的文件拷贝到`/etc/pki/ca-trust/source/anchors`目录下. 然后执行`update-ca-trust`命令. 最后重启docker. 再试一遍.

```
$ docker push registry.sky-mobi.com:5000/nginx:1.0.0
The push refers to a repository [sky-mobi.com:5000/nginx]
8963368d3c63: Pushed 
404361ced64e: Pushed 
1.0.0: digest: sha256:38fe56c7f45f7007e6d06ed473c015760bf24d5f048e430af8ce07e75029d630 size: 740
```

good.

这样, restful API就可以通过curl查询镜像了, 注意要加上**https前缀**哦.

```
$ curl https://registry.sky-mobi.com:5000/v2/_catalog
{"repositories":["nginx"]}
$ curl https://registry.sky-mobi.com:5000/v2/nginx/tags/list
{"name":"nginx","tags":["1.0.0"]}
```

完美.

不过, 同样是不正规的信任, 究竟自己创建自签名证书然后让客户端信任方便, 还是原来的直接配置`--insecure-registry`选项方便, 就看你取舍了.

关于用户名密码什么的, 我是真心不想做, 烦. 哪一个仓库这么干过? nodejs? pypi? yum? 
