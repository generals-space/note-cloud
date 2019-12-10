# RBAC权限管理

参考文章

1. [Kubernetes-基于RBAC的授权](https://www.kubernetes.org.cn/4062.html)
    - 介绍RBAC的概念与原理
2. [Creating sample user](https://github.com/kubernetes/dashboard/wiki/Creating-sample-user)
    - 实例讲解创建ServiceAccount对象及Binding对象, 以及如何得到该对象的访问密钥.
    - **注意: 示例中`cluster-admin`角色是系统预先创建的`ClusterRole`对象, 所以省略了角色创建的过程**
