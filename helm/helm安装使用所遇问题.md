
参考文章

1. [helm官方文档 Quickstart Guide](https://helm.sh/docs/intro/quickstart/)
2. [k8s之helm安装mysql](https://blog.csdn.net/eyeofeagle/article/details/102703065)
3. [helm 安装 mysql 相关注意事项及记录](https://blog.csdn.net/gs80140/article/details/93471482)

helm: 3.0.0
kubernetes: 1.16.2

按照helm官方文档安装好helm, 然后安装了一个`stable/mysql`的chart, 但是生成的 pod 总是处于 pending 状态. 

`describe`查了下 pod 的信息, 发现没有绑定 pvc 对象.

```
$ k describe pod mysql
...省略
Events:
  Type     Reason            Age                  From                    Message
  ----     ------            ----                 ----                    -------
  Warning  FailedScheduling  <unknown>            default-scheduler       pod has unbound immediate PersistentVolumeClaims (repeated 2 times)
  Warning  FailedScheduling  <unknown>            default-scheduler       pod has unbound immediate PersistentVolumeClaims (repeated 2 times)
```

然后我又查了下 pvc 对象

```
$ k get pvc
NAME    STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
mysql   Pending                                      nfs            34s
```

而且没有pv对象.

我以为是我操作哪里出了问题, 毕竟 helm 说是类似 yum/apt 之类的包管理工具, 但我也没听过 yum 安装东西要先建个目录什么的(依赖什么的不算).

查到参考文章2和3, 发现 pv 对象是要自己手动建的...

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv
spec:
  capacity:
    storage: 16Gi
  accessModes:
   - ReadWriteOnce
  hostPath:
    path: /tmp/helm-pv
    type: Directory

```

`apply`一下, pvc 会自动绑定, pod 也能正常启动了.
