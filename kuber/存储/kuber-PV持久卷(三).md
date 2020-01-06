# kuber-PV持久卷(二)-挂载NFS

参考文章

1. [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

```console
$ k get pv
NAME       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                          STORAGECLASS   REASON   AGE
pv-mysql   16Gi       RWO            Retain           Bound    default/data-mariadb-slave-0                           2m6s
```

status可能的状态: Available(初始创建) -> Bound(被某个pvc获得) -> Released(被pvc释放)

`accessModes`是pv/pvc的配置中的必填字段, 双方必须是一致的才能绑定(比如pv如果是`ReadWriteMany`, 而pvc是`ReadWriteOnce`就不行).

pv和pvc需要通过`accessModes`和`storageClassName`两个字段进行匹配, 同时需要pv的大小符合pvc的要求就可以绑定.

创建pv时实际的存储可以不事先存在, 比如, 如果是hostPath, 那么目标目录不存在也能创建pv.

Reclaim Policy(资源回收策略):

- Retain(保留, 持有): 被绑定的pvc释放后pv无法被其他pvc重新绑定(解绑的pv处于`Released`状态), 只能手动回收(即此pv无法再被使用). 这是出于对存储在pv中数据的安全考虑, 需要管理员手动回收此pv.
- Delete: 采用此策略的pv将随pvc的移除而移除.
- Recycle (1.16被废弃)

## `accessModes`实验

pv: 5G, Retain, ReadWriteOnce
pvc: 1G, ReadWriteOnce

第一个pvc可以成功绑定, 但是第二个pvc会一直处于pending状态(将第一个pvc删除再重新创建也无法再次绑定).

实验用配置.

```yaml
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