# kubectl get jsonpath获取特殊字符的字段

在 k8s 各资源的 status 块中, 貌似键名中可以包含一些特殊字段, 比如点号`.`, 斜线`/`等, 直接使用是无法获取到这些字段的值的.

```yaml
status:
  allocatable:
    middleware.currentcpu: 500
    middleware.disk/capacity: 1024000000
```

下面两种

```
k get node node名称 -o=jsonpath="{.status.allocatable.middleware.currentcpu}" 
k get node node名称 -o=jsonpath="{.status.allocatable['middleware.disk/capacity']}" 
```
