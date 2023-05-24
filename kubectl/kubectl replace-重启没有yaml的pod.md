# 重启没有yaml的pod

参考文章

1. [k8s 重启pod](https://www.jianshu.com/p/baa6b11062de)

在没有pod 的yaml文件时，强制重启某个pod

```
kubectl get pod Pod名 -n 命名空间 -o yaml | kubectl replace --force -f -
```

简单来说就是先得到目标pod的yaml文件, 然后使用`replace`子命令将正在运行的pod替换.

有一个pod的状态变成了`CrashLoopBackOff`, 想用这条命令重启.

```
# k get pod coredns-6967fb4995-w4pgk -n kube-system -o yaml | k replace -f -
pod/coredns-6967fb4995-w4pgk replaced
```

后半段的replace命令`--force`选项貌似是必须的, 不加不会删除crash的pod.

```
# k get pod coredns-6967fb4995-w4pgk -n kube-system -o yaml | k replace --force -f -
pod "coredns-6967fb4995-w4pgk" deleted
pod/coredns-6967fb4995-w4pgk replaced
```
