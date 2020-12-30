# kubectl版本不匹配导致的the server could not find the requested resource

参考文章

1. [Error: the server could not find the requested resource (post replicationcontrollers)](https://github.com/kubernetes/kubernetes/issues/9945)
2. [Determine what resource was not found from “Error from server (NotFound): the server could not find the requested resource”](https://stackoverflow.com/questions/51180147/determine-what-resource-was-not-found-from-error-from-server-notfound-the-se)

kubectl client: 1.9.2
kubectl server: 1.13.2

使用`kubectl edit`命令修改一个资源, 完成后按`:wq`保存退出时, 提示如下错误

```
Error from server (NotFound): the server could not find the requested resource
```

目标资源也没修改成功.

这其实是因为 kubectl 客户端与目标 kuber 集群的版本不匹配导致的, 更新一下本地的 kubectl 版本即可.

------

另外, 使用低版本的 kubectl 执行`kubectl get node -o wide`时, 可能无法显示 node 节点的 IP, 因为没有`INTERNAL-IP`这一列.
