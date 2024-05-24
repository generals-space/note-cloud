# kubectl scale横向扩缩容[replicas fast]

```
k scale --replicas=4 deploy mydeploy
```

支持该操作的资源只有`deploy`与`sts`.

## --current-replicas

除了可以指定资源对象的目标节点数外, 还可以预选进行筛选, 只扩缩符合条件的对象.

```log
$ k scale deploy esrally --replicas=1 --current-replicas=3
error: Expected replicas to be 3, was 2
```

只进行3->1的变换, 如果目标资源不是3副本, 就终止.
