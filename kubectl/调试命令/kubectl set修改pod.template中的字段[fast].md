# kubectl set修改pod.template中的字段[fast]

参考文章

1. [kubectl for Docker Users](https://kubernetes.io/docs/reference/kubectl/docker-cli-to-kubectl/)


```
k set env deploy mydeploy KEY=VALUE
```

`k set`对`deploy`, `sts`, `ds`的修改能力基本是相同的, 不过对 pod 能够修改的字段有限(只有`image`).

