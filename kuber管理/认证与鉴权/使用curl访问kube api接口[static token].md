# 使用curl访问kube api接口

参考文章

1. [使用 Kubernetes API 访问集群](https://kubernetes.io/zh/docs/tasks/administer-cluster/access-cluster-api/#%E4%B8%8D%E4%BD%BF%E7%94%A8-kubectl-%E4%BB%A3%E7%90%86)
    - 这篇文章中所说的, 不使用kubectl代理的访问方式, 其实就是取得`kube-system`命名空间下的, 名为`default`的`serviceAccount`的token, 然后获取该`default`用户的权限.
2. [用户认证](https://kubernetes.io/zh/docs/reference/access-authn-authz/authentication/)

apiserver本质是一个http服务器, 无论是kubectl, 还是operator, 最终都是通过http api进行通信.

apiserver支持多种token认证: static token, bootstrap token, service account token等.

token本质上就是一种密码, ta绑定了某个用户, 以至于客户端在使用某一token发起请求时, 服务端可以从token得到其对应的用户, 然后赋予该请求对应用户的权限.

## static token

static token是最简单的一种token, token内容和其绑定的用户都可以直接写在文本文件中, 如下

```
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx,admin,1
79e5a71195692f0fd989275b3ddfb4a4,admin,1
```

1. token字符串: 可为自定义的字符串, 长度不限, 没有固定格式
2. Role/ClusterRole 对象, 需要自行绑定权限
3. id: 目前不清楚有何作用, 本来 随机, 两行token可以指定同一个id....

将上述内容写入`/etc/kubernetes/pki/token_auth_file`, 然后开启apiserver的静态token认证方式.

上面两个token都可以使用.

```
curl -k -H 'Content-Type: application/json' -H 'Authorization: Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' 'https://127.0.0.1:6443/api/v1/namespaces/kube-system/pods'
curl -k -H 'Content-Type: application/json' -H 'Authorization: Bearer 79e5a71195692f0fd989275b3ddfb4a4' 'https://127.0.0.1:6443/api/v1/namespaces/kube-system/pods'
```
