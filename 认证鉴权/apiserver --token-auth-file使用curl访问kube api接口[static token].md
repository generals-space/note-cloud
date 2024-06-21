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
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx,cluster-admin,1
79e5a71195692f0fd989275b3ddfb4a4,cluster-admin,1
```

1. token字符串: 可为自定义的字符串, 长度不限, 没有固定格式
2. ~~Role/ClusterRole 对象, 需要自行绑定权限~~
3. id: 目前不清楚有何作用, 本来 随机, 两行token可以指定同一个id....

## 配置apiserver与验证

将上述内容写入`/etc/kubernetes/pki/token_auth_file`, 然后开启apiserver的静态token认证方式, 在`/etc/kubernetes/manifests/kube-apiserver.yaml`中的`command`添加如下选项

```
--token-auth-file=/etc/kubernetes/pki/token_auth_file
```

重启apiserver生效.

上面两个token都可以使用.

```bash
curl -k -H 'Content-Type: application/json' -H 'Authorization: Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' 'https://127.0.0.1:6443/api/v1/namespaces/kube-system/pods'
curl -k -H 'Content-Type: application/json' -H 'Authorization: Bearer 79e5a71195692f0fd989275b3ddfb4a4' 'https://127.0.0.1:6443/api/v1/namespaces/kube-system/pods'
```

## 补充

关于static token的第2列, 有几点需要注意, 该字段不是 Role/ClusterRole 类型(加了没用), 而是确确实实的"User"对象(也不是"Group"对象).

```bash
$ curl -k -H 'Content-Type: application/json' -H 'Authorization: Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' 'https://127.0.0.1:6443/api/v1/namespaces/kube-system/pods/kube-apisrer-k8s-master-01'
{
  "message": "pods \"kube-apiserer-k8s-master-01\" is forbidden: User \"cluster-admin\" cannot get resource \"pods\" in API group \"\" in the namespace \"kube-system\"",
  "reason": "Forbidden",
}
$ curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:6443/api/v1/pods'
{
  "message": "pods is forbidden: User \"cluster-admin\" cannot list resource \"pods\" in API group \"\" at the cluster scope",
  "reason": "Forbidden",
}
```

当一个用户, 使用成对的static token和user信息(以cluster-admin为例), 访问一个资源, 通过了认证阶段后, 将进入了[kubernetes-1.16.0](https://github.com/kubernetes/kubernetes/blob/v1.16.0/pkg/registry/rbac/validation/rule.go#L186)中的鉴权逻辑.

```go
		for _, clusterRoleBinding := range clusterRoleBindings {
			klog.Infof("cluster role binding name: %s", clusterRoleBinding.Name)
			subjectIndex, applies := appliesTo(user, clusterRoleBinding.Subjects, "") 
			if !applies {
				continue
			}
            // ...省略
        }
```

然后遍历系统中所有的 rolebinding/clusterrolebinding, 调用`appliesTo()`函数, 比较每个 binding 中的 Subject 列表与 user 变量是否相符.

> 这里的 user 就是通过认证阶段后, 得到的用户/组名(本例中为"cluster-admin"), 目标资源类型, 操作类型等信息.

```go
// appliesTo returns whether any of the bindingSubjects applies to the specified subject,
// and if true, the index of the first subject that applies
func appliesTo(user user.Info, bindingSubjects []rbacv1.Subject, namespace string) (int, bool) {
	for i, bindingSubject := range bindingSubjects {
		if appliesToUser(user, bindingSubject, namespace) {
			return i, true
		}
	}
	return 0, false
}

func appliesToUser(user user.Info, subject rbacv1.Subject, namespace string) bool {
	switch subject.Kind {
	case rbacv1.UserKind:  return user.GetName() == subject.Name
	case rbacv1.GroupKind: return has(user.GetGroups(), subject.Name)
	case rbacv1.ServiceAccountKind:
        // ...省略
	default:
		return false
	}
}
```

注意, 上述代码中`case rbacv1.UserKind:`要求, `cluster-admin`必需是一个User对象才可以, 否则就通不过鉴权.
但实际上, kube系统中是没有 User/Group 类型的资源的, 那么如何定义一个 User/Group?

正确的做法是, 直接在`cluster-admin`的`ClusterRoleBinding`资源的Subject中添加一个Kind为User的对象就可以了.

```yaml
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:masters
- apiGroup: rbac.authorization.k8s.io
  kind: User                            ## 这里可以直接定义为 User 类型
  name: cluster-admin
```

