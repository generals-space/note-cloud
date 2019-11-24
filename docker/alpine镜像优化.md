## 设置apk的阿里云镜像站点

```bash
sed -i 's/http:\/\/dl-cdn.alpinelinux.org/https:\/\/mirrors.aliyun.com/' /etc/apk/repositories
```

## 软件介绍

### `ca-certificates`: 

通用的证书文件, CentOS下一般随系统安装. 

在alpine的`/etc/ssl`目录下默认有cert.pem(文件), certs(目录), openssl.cnf(文件), x509v3.cnf(文件), 其中`certs`目录为空. 

在golang程序里请求https接口时, 会报这个错误`x509: certificate signed by unknown authority`, 但是在浏览器中这个网站的证书是合法的, 这是由于alpine系统中没有安装通用机构的https根证书. 这样会让用户认为所有的网站都是不安全的.

安装`ca-certificates`包后就会发现这个目录下多了很多证书. 再发起https请求就可以了.

### `tzdata`

时区文件, 默认`/usr/share`目录下是没有`zoneinfo`目录的, 也就没有办法调整系统时区了. 安装这个包后可以链接到`/etc/localtime`来修改.

### `busybox-extras`

Alpine 镜像中的 telnet在 3.7 版本后被转移至 busybox-extras 包中, ping, nc等命令也在这个包里.

### `alpine-sdk`

alpine-sdk与ubuntu的build-essential包作用相似, 是开发的头文件和链接库.

