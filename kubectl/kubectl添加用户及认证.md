# kubectl添加用户及认证

参考文章

1. [为Kubernetes集群添加用户](https://zhuanlan.zhihu.com/p/43237959)
    - 介绍了服务账号(ServiceAccount)和普通意义上的用户(User)在概念和应用上的区别
    - kubectl操作kuber集群的3种认证方式: ca.crt+key; 静态token; 静态密码;
    - 实例演示为新用户签发证书, 并添加角色, 然后配置kubectl允许其操作kuber集群的方案
2. [How to Add Users to Kubernetes (kubectl)?](https://stackoverflow.com/questions/42170380/how-to-add-users-to-kubernetes-kubectl)
    - 高票回答实例展示了创建ServiceAccount对象后, 将生成的token配置到kubectl客户端实现授权访问的方法
    - `kubectl config`命令各参数的应用, 与生成的`config`文件内容的对比可以了解一下

