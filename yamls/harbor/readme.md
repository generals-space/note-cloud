# 关于Harbor

## 写在前面

官方文档写的真烂, 不仅烂, 还过时. harbor本身都1.8了, kuber相关的部署配置还是1.2.

下面分析一下harbor的部署流程.

按照官方最主要的部署文档来说, 其实下载一个release版本的压缩包, 解压后执行其中的`install.sh`就可以了. 这个脚本会根据解压后的`harbor.yml`文件中的相关配置, 生成各组件需要的配置文件, 同时也会自动创建密钥对, 从docker hub下载相关镜像, 最后通过docker-compose启动.

需要注意的是, 生成配置文件的步骤并没有在`install.sh`脚本, 而是通过一个`goharbor/prepare`镜像来完成. 

真的是骚操作了, 我本来还想读一下这个脚本的...

不过最终我也没通过官方文档或是阅读脚本把kuber的配置文件搞定...

而是通过执行`install.sh`一遍, 将harbor各组件启动完成, 并且各配置字段都算比较熟悉之后再移植过去的.

> 当前目录下的`docker-compose.yml`就是`install.sh`生成的, 我借鉴了其中的volume配置, 编写了对应的kuber部署配置文件.

我使用NFS提供了存储服务, NFS Server的配置可以参考`exports`文件.

## `harbor.yml`理解

------

首先`hostname`字段是必须要定义的, 同时`https`也要解开相关注释, 因为即使在web界面上无所谓, 可以通过http访问, 但使用`docker login`要求服务必须是https接口(自签发证书也行). 

但是安装脚本并不会帮你生成证书和密钥, 我们需要事先创建好, 然后修改`https.certificate`和`https.private_key`这两处路径.

------

`database.password`字段定义了数据库的密码(postgres超级用户), 安装脚本会自动下载postgres和redis镜像并启动.

如果你想指定各组件使用外部的数据库和redis服务, 需要修改`external_database`和`external_redis`部分.

...我好像直接就把两个`external`的注释解开了, 默认配置没试过, 最后数据库和redis的指向还是安装脚本启动的那两个.

由于整个服务最终是用docker-compose启动的, 所以数据库/redis的host都要写成compose配置中的服务名称.

安装脚本下载的是harbor封装过的数据库镜像, 其中默认创建了3个数据库, 对应harbor, clair和notray 3个服务. 不过后两者我没启用, 先不管ta.

------

最终生成的`docker-compose.yml`同目录的还有一个`common`目录, 是各组件用到的配置文件, 通过volume挂载到容器内部.

------

clair, jobservice, log这3个部分可以不用动.
