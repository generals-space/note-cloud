# RBAC权限管理

参考文章

1. [浅析 kubernetes 的认证与鉴权机制](https://blog.tianfeiyu.com/2019/08/18/k8s_auth_rbac/)
    - 安全措施的三个步骤: 认证、授权、准入控制, 对ta们各自的作用划分地很清晰
    - 对Service Account资源的作用的介绍非常易理解
2. [Kubernetes API访问控制](https://kubernetes.io/zh/docs/reference/access-authn-authz/controlling-access/)
    - 官方文档

对于一个已经存在绑定关系的`Pod`, `ServiceAccount`, `Role(ClusterRole)`和`RoleBinding(ClusterRoleBinding)`, 更新`Role(ClusterRole)`中的资源权限配置, 是不需要重启`Pod`与其中的进程的, 直接就可以生效.
