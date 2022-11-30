# Secret对象认识[generic tls]

参考文章

1. [Kubernetes中的Secret配置](https://www.cnblogs.com/Leslieblog/p/10158429.html)
    - 对于`generic`类型的`Secret`资源, base64的编码和解码都是手动完成的...有什么意义???
    - 没有讲`tls`类型的`Secret`的使用方法.
2. [Kubernetes-token（四）](https://www.jianshu.com/p/1c188189678c)

`secret`类型目前有3种: 

1. generic(通用类型, 任意文件/目录/键值对) type为`Opaque`
2. tls(密钥配置) type为`kubernetes.io/tls`, 同时包括`crt`与`key`.
3. docker-registry: 私有仓库的认证信息, 使用场景比较固定.

上面说的是面向用户的接口, 我们可以通过`k create secret`创建以上3种类型的`Secret`资源, 但是`Secret`其实还有其他类型.

每个`SA`资源, 都绑定了一个`Secret`对象, 而这些`Secret`对象的类型是`kubernetes.io/service-account-secret`.

还有像`bootstrap.kubernetes.io/token`等, 应该无法直接通过`k create secret`创建, 可以使用`kubeadm token list`试试.
