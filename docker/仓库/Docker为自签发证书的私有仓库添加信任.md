# Docker为自签发证书的私有仓库添加信任

参考文章

1. [Test an insecure registry](https://docs.docker.com/registry/insecure/#using-self-signed-certificates)
    - Linux/Win下添加证书的方法
2. [Add TLS certificates](https://docs.docker.com/docker-for-mac/#add-tls-certificates)
    - MacOS下添加证书的方法
3. [Adding Self-signed Registry Certs to Docker & Docker for Mac](https://blog.container-solutions.com/adding-self-signed-registry-certs-docker-mac)
4. [苹果电脑Mac添加Docker Nexus自制证书](https://blog.csdn.net/happyfreeangel/article/details/90773354)

有一些自建的私有仓库使用的是自签发的证书, 在使用docker pull时会被(docker客户端)拒绝掉.

```
$ docker pull www.test.com/library/terway-vlan
Error response from daemon: Get https://www.test.com/v2/: x509: certificate signed by unknown authority
```

对证书的信任一般分为2个层面, 一个是应用层面, 比较docker/curl等应用本身可指定是否信任或忽略目标地址的证书, 一种是系统层面, 系统层面的信任是全局的.

Linux下可以在`/etc/docker/certs.d/www.test.com/`目录下放置服务端证书(docker服务不用重启), 这是针对docker应用本身的配置, 使用curl去请求相同的地址时, 仍然会显示证书非法. 或是将证书文件拷贝到`/etc/pki/ca-trust/source/anchors`目录下. 然后执行`update-ca-trust`命令, 这种方法需要重启docker服务才能生效.

Mac下有点不一样, 第一种方法无法生效(见参考文章3), 只能使用第二种. 方法如下

```
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ./www.test.com.crt
```

> 注意: 同样需要重启docker服务.
