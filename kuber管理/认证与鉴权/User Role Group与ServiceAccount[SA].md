参考文章

1. [kubernetes：kubernetes api访问控制、useraccount和serviceaccount的创建和绑定、rbac基于角色的访问授权](https://blog.csdn.net/weixin_43384009/article/details/105980976)
    - 关于`Group`的描述
    - kuber 提供了四个内置的`ClusterRole`供用户直接使用: `cluster-amdin`, `admin`, `edit`, `view`
2. [Kubernetes-一文详解ServiceAccount与RBAC权限控制](https://blog.51cto.com/mageedu/2553145)
    - Kubernetes中账户区分为: User Accounts(用户账户) 和 Service Accounts(服务账户) 两种
    - `UserAccount`是给kubernetes集群外部用户使用的, 例如运维或者集群管理人员, 使用kubectl命令时用的就是UserAccount账户;
        - `UserAccount`是全局性的, 在集群所有namespaces中, 名称具有唯一性, 默认情况下用户为admin;
    - `ServiceAccount`是给运行在Pod的程序使用的身份认证, Pod容器的进程需要访问API Server时用的就是ServiceAccount账户;
       - `ServiceAccount`仅局限它所在的namespace, 每个namespace创建时都会自动创建一个default service account;
       - 创建Pod时, 如果没有指定`ServiceAccount`, Pod则会使用default Service Account.
3. [Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
    - 官方文档
4. [kubernetes安全控制认证与授权(一)](https://blog.csdn.net/yan234280533/article/details/75808048)
    - kuber 给没有通过认证的请求一个特殊的用户名`system:anonymous`和组名`system:unauthenticated`.


`RoleBinding`与`ClusterRoleBinding`可绑定的主体有3类: `Users`, `Groups`与`ServiceAccount`

## Group的含义

kuber 拥有"用户组"(Group)的概念

`ServiceAccount`对应内置"用户"的名字是: `system:serviceaccount:<ServiceAccount名字>`, 而用户组所对应的内置名字是: `system:serviceaccounts:<Namespace名字>`

示例1: 表示mynamespace中的所有ServiceAccount

```yaml
subjects:
- kind: Group
  name: system:serviceaccounts:mynamespace
  apiGroup: rbac.authorization.k8s.io
```

示例2: 表示整个系统中的所有ServiceAccount

```yaml
subjects:
- kind: Group
  name: system:serviceaccounts
  apiGroup: rbac.authorization.k8s.io
```
