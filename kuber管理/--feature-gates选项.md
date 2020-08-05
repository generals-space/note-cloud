参考文章

1. [(二进制安装)k8s1.9 证书过期及开启自动续期方案.md](https://blog.csdn.net/feifei3851/article/details/88390425)

```
--feature-gates=RotateKubeletClientCertificate=true,RotateKubeletServerCertificate=true
```

`--feature-gates`多个参数用逗号`,`分隔.
