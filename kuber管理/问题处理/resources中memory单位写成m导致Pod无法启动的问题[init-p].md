# resources中memory单位写成m导致Pod无法启动的问题[init-p]

参考文章

1. [Kubernetes因限制内存配置引发的错误](https://cloud.tencent.com/developer/article/1411527)
    - CPU-单位m、内存-单位Mi

实际是因为我的`memory`的单位写成了`m`, 其实应该是`Mi`才对...

------

以下都是错误经验...Ծ‸Ծ

```yaml
      - name: centos7
        image: generals/centos7
        resources:
          limits:
            memory: 100m
```

当`limits`下只有`memory`, 没有`cpu`字段时, Pod 将一直处于`ContainerCreating`状态, 使用`describe`查看, 可以得到如下信息

```
Events:
  Type     Reason                  Age                 From                    Message
  ----     ------                  ----                ----                    -------
  Normal   Scheduled               <unknown>           default-scheduler       Successfully assigned default/test-ds-96gqz to k8s-master-01
  Warning  FailedCreatePodSandBox  57s (x3 over 67s)   kubelet, k8s-master-01  Failed create pod sandbox: rpc error: code = Unknown desc = failed to start sandbox container for pod "test-ds-96gqz": Error response from daemon: OCI runtime create failed: container_linux.go:346: starting container process caused "process_linux.go:319: getting the final child's pid from pipe caused \"read init-p: connection reset by peer\"": unknown
```

`requests`就随便了, 没影响. 另外, 单纯设置`limits.cpu`也是可以的...
