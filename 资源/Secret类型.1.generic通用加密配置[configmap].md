# Secret类型.1.generic通用加密配置[configmap]

参考文章

1. [Kubernetes中的Secret配置](https://www.cnblogs.com/Leslieblog/p/10158429.html)
    - 对于`generic`类型的`Secret`资源, base64的编码和解码都是手动完成的...有什么意义???
    - 没有讲`tls`类型的`Secret`的使用方法.
2. [Kubernetes-token（四）](https://www.jianshu.com/p/1c188189678c)

`generic`类型的`Secret`资源与`ConfigMap`几乎完全一致, 使用方法也大致相同. 

```
kubectl create secret generic db-user-pass --from-file=./username.txt --from-file=./password.txt
```

除了`--from-file`, 同样可以使用`--from-literal`, `--from-env-file`, 具体可见ConfigMap的文章.

不过要注意, 这种方式创建出来的Secret的type字段其实是`Opaque`, 不明白为什么.

```console
$ k get secret
NAME            TYPE      DATA   AGE
db-user-pass    Opaque    1      385d
```

------

唯一区别应该是, `kubectl describe`一个`Secret`资源时, 并不会打印其中的内容, 然后就是存入data块中的value值是经过base64加密的.

不过感觉没什么用, `kubectl get`仍然可以通过`-o yaml`得到其中存储的信息, 解密也很简单, 自欺欺人罢了...
