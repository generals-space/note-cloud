helm2分为客户端helm与服务端tiller.

如果服务端已经部署, 可以单纯安装客户端.

```
helm init --client-only
```

查找charts, 不像helm3那样需要指定`hub`或`repo`, 直接就是全局搜索.

```
helm search prometheus
```

安装charts

```
helm install bitnami/kubeapps --name kubeapps --version 1.2.0
```
