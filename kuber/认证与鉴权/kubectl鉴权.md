# kubectl鉴权

参考文章

1. [配置kubectl客户端通过token方式访问kube-apiserver](https://www.cnblogs.com/tianshifu/p/7841007.html)

2. [在 kubectl 中使用 Service Account Token](https://blog.csdn.net/kwame211/article/details/78981403)
    - 使用SA对象中的使用的`service-account-token`类型的`secret`资源中的token值, 作为kubectl的认证值.

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

```console
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

------

参考文章2中给出了使用SA引用的Secret类型资源中的token数据作为kubectl配置中的`user.token`值的示例. 有以下几点需要注意:

1. `kubernetes.io/service-account-token`类型的`secret`资源对象, 其中的`token`数据要经过base64解密才能填写到kubeconfig的`user.token`字段;
2. 各ns下默认创建的名为`default`的SA是没有绑定`Role/RoleBinding`权限的, 需要手动创建, 否则没有任何权限.
