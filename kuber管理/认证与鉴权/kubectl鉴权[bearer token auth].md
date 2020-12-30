# kubectl鉴权[bearer token auth]

参考文章

1. [配置kubectl客户端通过token方式访问kube-apiserver](https://www.cnblogs.com/tianshifu/p/7841007.html)
2. [在 kubectl 中使用 Service Account Token](https://blog.csdn.net/kwame211/article/details/78981403)
    - 使用SA对象中的使用的`service-account-token`类型的`secret`资源中的token值, 作为`kubectl`的认证依据.
3. [访问控制](https://kubernetes.feisky.xyz/extension/auth)
    - Kubernetes 对 API 访问提供了三种安全访问控制措施: 认证、授权和 Admission Control
    - **认证**解决用户是谁的问题, **授权**解决用户能做什么的问题, **Admission Control**则是资源管理方面的作用
4. [Kubernetes 中的用户与身份认证授权](https://jimmysong.io/kubernetes-handbook/guide/authentication.html)
    - kuber中认证的多种手段

kuber集群对于API访问, 在认证阶段提供了多种方式(见参考文章4), 开发者在访问 kuber API 时可以自行选择认证手段, 只要一种成功就算成功. 

`kubectl`本身支持其中的3种, ta的`config`子命令有一个`set-credentials`参数, ta主要配置`kubeconfig`文件中的`user`字段, 如下

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

而`set-credentials`可以指定不同的认证方式, 使用`-h`选项可以查看到如下帮助信息

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

可以看到上面提到了3种认证方式

1. Client-certificate: 最常用的证书/密钥对, 一般使用kubeadm创建的集群都会生成此种方式的配置文件.
2. Bearer Token:
3. Basic Auth: 就是用户名密码.

其中`Bearer Token`与`Basic Auth`不可同时使用.

## Bearer Token

```
kubectl config set-credentials user02 --token=7176d48e4e66ddb3557a82f2dd316a93
```

上面的命令可以设置用户`user02`的认证字段, 但是要让`apiserver`认可此token, 需要在apiserver的启动参数中, 通过`--token-auth-file`指向一个token文件, 其格式如下

```
7176d48e4e66ddb3557a82f2dd316a93,user02,1
```

不只如此, 参考文章2中给出了使用`SA`引用的`Secret`类型资源中的`token`字段作为kubectl配置中的`user.token`值的示例. 有以下几点需要注意:

1. `kubernetes.io/service-account-token`类型的`secret`资源对象, 其中的`token`数据要经过`base64`解密才能填写到`kubeconfig`的`user.token`字段;
2. 各`ns`下默认创建的名为`default`的SA是没有绑定`Role/RoleBinding`权限的, 需要手动创建, 否则没有任何权限.
