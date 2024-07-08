参考文章

1. [《Kubernetes设计与实现》第九章：准入控制器 9.1.1 概述](https://renhongcai.gitbook.io/kubernetes/di-jiu-zhang-zhun-ru-kong-zhi-qi)
2. [k8s的Mutating webhook](https://www.jianshu.com/p/a76e8f7d13b7)
    - `--enable-admission-plugins`
    - `--disable-admission-plugins`

每个针对Kubernetes资源对象的操作请求，都要经过kube-apiserver的层层审查才会被放行。对于读操作而言，需要经过认证（是否是合法用户）和鉴权（是否拥有权限）；而对于写操作而言，除了要经过认证和鉴权外，还要检查请求是合乎要求，只有顺利通过这些审查才会被持久化到etcd存储中。

kube-apiserver支持配置多个准入控制器，准入控制器分为修改型（Mutating）控制器和校验型（Validating）控制器。修改型控制器会自动根据指定的策略对请求进行修改，而校验型控制器则只是单纯地检查请求是否合乎要求，充当“看门狗”的角色。

多个准入控制器以插件（Webhook Plugin）的形式被组织起来，kube-apiserver在审查请求时会先把请求交给修改型控制器对请求进行必要的修改，然后再将请求交给校验型控制器进行审查。下图展示了请求的审查完整路径，以及准入控制器所在的位置：

![](https://kubernetes.io/images/blog/2018-06-05-11-ways-not-to-get-hacked/admission-controllers.png)

API请求到达kube-apiserver后会先进行认证（Authentication）和鉴权（Authorization），然后把请求交给修改型准入控制器进行必要的修改（多个修改型准入控制器串行执行），当所有修改型准入控制器执行完毕后，再使用OpenAPI 校验功能进行初步的语法校验，接着再把请求交给校验型准入控制器进行语法或语义的校验（多个修改型准入控制器并行执行），最后再写入etcd。上面中的任何一个审查环节、任何一个准入控制器返回失败，都会造成请求被拒绝。

