# Role与ClusterRole中subresource的配置方式

参考文章

1. [How to refer to all subresources in a Role definition?](https://stackoverflow.com/questions/57872201/how-to-refer-to-all-subresources-in-a-role-definition)

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: temp-role
  namespace: stackoverflow
rules:
- apiGroups: [""]
  resources:
  - pods
  - pods/log
  verbs:
  - get
```
