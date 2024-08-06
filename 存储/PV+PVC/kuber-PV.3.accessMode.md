# kuber-PV.3.accessMode与reclaimPolicy

参考文章

1. [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

## `accessModes`实验

- `ReadWriteOnce(RWO)`: 读写权限, 但是只能被单个节点挂载
- `ReadOnlyMany(ROX)`: 只读权限, 可以被多个节点挂载
- `ReadWriteMany(RWX)`: 读写权限, 可以被多个节点挂载

先说结论

`accessMode`的`Once`与`Many`指的不是PV与PVC之间的一对多关系, 而是运行在不同node节点的Pod是否可以挂载同一个PVC.

我最开始设想的是, 一个申请了5G空间的PV, 可以同时被5个申请1G空间的PVC绑定. 做完一系列实验后才意识到这个想法太蠢了, 同一个目录被多个PVC共用根本无法区分, PVC又不会创建子目录.

实际上, 一个PV只能与一个PVC绑定, 反过来也一样.

关于`Once`与`Many`的区别, 首先, 不同的`provisioner`的支持是不同的, 比如`hostPath`就只支持`Once`而不支持`Many`, 而`NFS`都支持.

以原生的`hostPath`为例, `pv -> pvc -> pod`的创建过程中, pv和pvc的创建只是集群内部资源的创建, 这两者不存在于任一节点, 而是在 etcd 里. 只有对应的pod的出现, 才会在Pod被分配的节点上创建`hostPath`指定的目录.

在实验NFS时, 先在集群外部创建NFS服务端(NFS无法挂载本机上的共享目录, 否则会出错), 然后创建`Many`类型的PV与PVC资源. 然后就是创建一个daemonset资源, 关于volume的部分如下

```yaml
      volumes:
        - name: data-vol
          persistentVolumeClaim:
            claimName: pvc1
      containers:
      - name: centos7
        volumeMounts:
        - name: data-vol
          mountPath: /data
```

这样生成的pod全部挂载了同一个pvc对象, 就是说最终大家都挂载了同一个目录.

但是尝试将PV和PVC的`Many`改成`Once`, 重新创建`pv -> pvc -> pod`, 结果和之前没区别, 就是说只要使用NFS, 不管是`Once`还是`Many`, 最终都是`Many`...

当我对`accessMode`有了比较正确的认知后, 我又回头对`hostPath`做了次验证.

```yaml
---
kind: PersistentVolume
apiVersion: v1
metadata:
  name: pv1
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteMany
  hostPath:
    path: /tmp/pv1
    type: Directory
  ## nfs: ############################################# 只有这里不一样.
  ##   path: "/mnt/nfsfold"
  ##   server: 172.16.91.128
```

```yaml
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: pvc1
spec:
  resources:
    requests:
      storage: 1Gi
  accessModes:
    - ReadWriteMany
```

一对一的pv与pvc情况下, 两者都可以创建成功. 但是创建ds时, 生成的pod却都无法启动, 都处于`ContainerCreating`状态. 对ta们`describe`一下, 发现有如下输出

```log
Events:
  Type     Reason       Age               From                    Message
  ----     ------       ----              ----                    -------
  Normal   Scheduled    <unknown>         default-scheduler       Successfully assigned default/local-pv-ds-gn4cv to k8s-master-01
  Warning  FailedMount  2s (x7 over 33s)  kubelet, k8s-master-01  MountVolume.SetUp failed for volume "pv1" : hostPath type check failed: /tmp/pv1 is not a directory
```

...原来只是因为宿主机上的路径不存在, 手动创建后移除出问题的pod重新运行就可以了.

ok, 然后把pv和pvc的`Many`改成`Once`, 再次实验. 与上面没有区别.

就是说, 其实`hostPath`与`NFS`很像, 只不过是完全反过来, 不管是`Once`还是`Many`, 都是`Once`.

而且绑定的pv与pvc, 一旦pvc被移除了, pv就会变成`Release`状态, 无法再次被绑定. 只要pvc不被移除, 就仍可以被多个pod绑定.

> 不过虽然`hostPath`不支持`ReadWriteMany(RWX)`, 但是在部署文件中这么写也没关系, 创建pv和pvc资源时也不会报错. 而`rancher`的`local-path`则不是, 人家明确说明只支持`Once`, 在PVC中写`Many`的话根本无法创建成功.

------

下面是我之前做的实验, 主要当时没搞清楚`Once`和`Many`根本不是pv和pvc的一对一还是一对多, 废了很多时间.

### 1. ReadWriteOnce

- pv: 5G, Retain, ReadWriteOnce
- pvc: 1G, ReadWriteOnce

实验用配置.

```yaml
---
kind: PersistentVolume
apiVersion: v1
metadata:
  name: pv1
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /tmp/pv1
    type: Directory
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: pvc1
spec:
  resources:
    requests:
      storage: 1Gi
  accessModes:
    - ReadWriteOnce
```

使用上面的部署文件创建的第一个pvc可以成功绑定pv对象, 但是再创建第二个pvc(注意将名称由`pvc1`改成`pvc2`)就会一直处于`pending`状态(将第一个pvc删除再重新创建也无法再次绑定).

### 2. ReadWriteMany

- pv: 5G, Retain, ReadWriteMany
- pvc: 1G, ReadWriteMany

部署文件就不再贴出来了, 直接将上面的做一下修改即可.

...嗯, 结果与上面的没有什么不同, 第二个pvc也是`pending`.
