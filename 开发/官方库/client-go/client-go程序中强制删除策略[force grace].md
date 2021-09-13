# client-go程序中强制删除策略[force grace]

参考文章

1. [how to delete pod forcely?](https://github.com/fabric8io/kubernetes-client/issues/1195)

使用`kubectl`强制删除一个 pod 需要两个参数: `--force`, `--grace-period=0`

```
kubectl delete pod pod名称 --grace-period=0 --force
```

但在使用 client-go 提供的`Delete()`方法中, `metav1.DeleteOptions{}`参数中并没有`force`成员.

其实在程序里需要进行如下操作

```go
// force 需要为指针类型
force := int64(0)
err = client.CoreV1().Pods(namespace).Delete(
    podName, &metav1.DeleteOptions{
        GracePeriodSeconds: &force,
    }
)
```
