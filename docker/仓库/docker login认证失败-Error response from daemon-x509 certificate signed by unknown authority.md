# docker认证失败-Error response from daemon-x509 certificate signed by unknown authority

参考文章

1. [docker下载私有仓库镜像失败：Error response from daemon: Get $ip:5000/v2/: http: server gave HTTP response to HT](https://blog.csdn.net/liurizhou/article/details/88192945)

```console
$ docker login 10.1.11.210
Username: admin
Password: 
Error response from daemon: Get https://10.1.11.210/v2/: x509: certificate signed by unknown authority
```

按照参考文章中的配置, 在`/etc/docker/daeon.json`中添加一个`insecure-registries`字段, 值为列表类型, 如下

```json
{
    "insecure-registries":["xxx.xxx.xxx:5000"]
}
```

重启docker, 再登录就可以了.

```console
$ docker login 10.1.11.210
Username: admin
Password: 
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded
```
