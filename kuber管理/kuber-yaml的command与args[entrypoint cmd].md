# kuber-yaml的command与args[entrypoint cmd]

参考文章

1. [[kubernetes]spec yaml定义command/args](https://blog.csdn.net/jettery/article/details/86498516)

docker 的`entrypoint`等同于 kuber 的`command`, 且 docker 的`cmd`则等同于 kuber 的`args`.

```yaml
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        command:
        - nginx
        args:
        - --debug
```

如果一个 dockerfile 中声明了`entrypoint`指令, 在 kuber 中部署时, `command`是可以将其覆盖掉的.
