# cannot assign requested address

创建`ClusterIP`类型的Service时报如下错误.

```
$ kubectl create -f ./svc.yaml
Get https://[::1]:6443/api/v1/namespaces/oa/resourcequotas: dial tcp [::1]:6443: connect: cannot assign requested address
```

但是NodePort类型就不会.

目前没有找到解决办法, 留空吧. ???
