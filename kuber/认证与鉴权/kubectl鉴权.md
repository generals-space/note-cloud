# kubectl鉴权

参考文章

1. [配置kubectl客户端通过token方式访问kube-apiserver](https://www.cnblogs.com/tianshifu/p/7841007.html)

之前一直不清楚这三种认证方式是从哪里来的, 直到今天才有了这样的认识.

`kubectl`工具的`config`子命令有一个`set-credentials`参数, ta主要配置了`kubeconfig`文件中的`user`字段, 如下

```yaml
users:
- name: user01
  user:
    client-certificate: client.crt 
    client-key: client.key 
- name: user02
  user:
    token: 7176d48e4e66ddb3557a82f2dd316a93
```

指定`-h`选项可以查看到如下帮助信息

```
$ kubectl config set-credentials -h

Sets a user entry in kubeconfig

 Specifying a name that already exists will merge new fields on top of existing values.

  Client-certificate flags:
  --client-certificate=certfile --client-key=keyfile

  Bearer token flags:
    --token=bearer_token

  Basic auth flags:
    --username=basic_user --password=basic_password

 Bearer token and basic auth are mutually exclusive.
```

可以看到上面提到了3种认证方式:

1. Client-certificate: 最常用的证书/密钥对, 一般使用kubeadm创建的集群都会生成此种方式的配置文件.
2. Bearer Token:
3. Basic Auth: 就是用户名密码.

其中`Bearer Token`与`Basic Auth`不可同时使用.

```
kubectl config set-credentials user02 --token=7176d48e4e66ddb3557a82f2dd316a93
```

上面的命令可以设置用户`user02`的认证字段, 但是要让apiserver认可此token, 需要在apiserver的启动参数中, 通过`--token-auth-file`指向一个token文件, 其格式如下

```
7176d48e4e66ddb3557a82f2dd316a93,user02,1
```
