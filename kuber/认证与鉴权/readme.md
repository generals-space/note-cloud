参考文章

1. [官方文档 - Accessing Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/)
    - 通过`kubectl proxy`映射`apiserver`端口, 然后使用`curl`工具通过rest方式访问.
    - 通过查找`secret`资源, 在`curl`请求头中附加`Authorization: Bearer $TOKEN`访问`apiserver`.
    - 通过`client-go`客户端编写程序, 从pod内部访问`apiserver`.
    - 非常实用
2. [Kubernetes RBAC 详解](https://www.qikqiak.com/post/use-rbac-in-k8s/)
    - 没有创建`SA`, 直接使用`openssl`创建密钥对, 在kubectl配置对应的`user`块, 然后`Binding`对象指定的主体为`user`块中的名称即可使用.
3. [官方文档 Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
    - 给出了`Role/ClusterRole`和`RoleBinding/ClusterRoleBinding`的声明示例.
4. [官方文档 kubernetes-api](https://kubernetes.io/docs/reference/kubernetes-api/)
    - 给出了不同版本kuber的API手册链接
    - 手册中包含各种资源的`get`/`list`/`watch`等操作的使用示例, 选项参数, 及对应的restful请求示例, 值得参考.


## 1. 概述

kuber中RBAC的运行原理

1. 首先明确kuber存在各种资源: Node, Pod, Service等, 对于这些资源也分别有各自的操作: create, get, list, watch等;
2. kuber中存在2种可以拥有这些权限的资源: Role/ClusterRole, 管理员可以对这两种实例对象分配指定类型资源的不同的操作权限;
3. 但Role/ClusterRole并不是执行操作的主体, 目前的执行主体只有SA(貌似还有`User`和`Group`, 不过不常用);
4. SA资源并不能直接设置Role/ClusterRole角色, 必须分别通过RoleBinding/ClusterRoleBinding才可以实现绑定功能;
5. 然后SA可以被赋予给Pod(或是deploy, ds, sts等), 这样, pod中运行的程序, 就可以实现受限地访问集群.

> 思路大致理清了, 但总觉得`RoleBinding`/`ClusterRoleBinding`有此多余, 在声明User/SA资源实例的时候不能把`Binding`的功能集成进去吗? 

另外, 参考文章2中有提到`User`类型的`Binding`主体, 与SA同级. 不过kuber中并没有`User`类型的资源, `Binding`中的`User`对象只能用于kubectl配置文件中的`user`字段.

------

在参考文章2中, 由于`Role`和`RoleBinding`资源都是建立在`kube-system`空间下, 所以与之绑定的`User`只拥有`kube-system`的权限, 而没有`default`空间的权限. 如果需要角色获取多个ns的权限, 可以使用`ClusterRole`和`ClusterRoleBinding`.

在`Role/ClusterRole`的声明中, 权限规则的定义不仅需要指明资源类型, 还需要指定其所属的`apiGroup`. 如下示例

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  # "namespace" omitted since ClusterRoles are not namespaced
  name: secret-reader
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
```

当前版本支持的资源及对应的apiGroup可以直接通过`k api-resources`命令查看. 也可以查阅参考文章4中给出的API手册, 更加完善且全面.
