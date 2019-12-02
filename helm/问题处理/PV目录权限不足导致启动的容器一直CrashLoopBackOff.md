# PV目录权限不足导致启动的容器一直CrashLoopBackOff

```
$ helm install redis stable/redis
$ k get pod
NAME             READY   STATUS              RESTARTS   AGE
redis-master-0   0/1     CrashLoopBackOff    9          25m
redis-slave-0    0/1     ContainerCreating   0          2m17s
```

`describe`没有查看到异常, 怀疑了很久, 好在pod有日志打印出来.

```
# k logs -f redis-master-0
 06:41:30.35 INFO  ==> ** Starting Redis **
1:C 02 Dec 2019 06:41:30.358 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
1:C 02 Dec 2019 06:41:30.358 # Redis version=5.0.7, bits=64, commit=00000000, modified=0, pid=1, just started
1:C 02 Dec 2019 06:41:30.358 # Configuration loaded
1:M 02 Dec 2019 06:41:30.359 # Can't open the append-only file: Permission denied
```

用于redis的pv使用本地路径, 没有额外修改权限.

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: redis-master
spec:
  capacity:
    storage: 8Gi
  accessModes:
   - ReadWriteOnce
  hostPath:
    path: /tmp/helm-redis-master
    type: Directory
```

把`/tmp/helm-redis-master`目录权限改成777就可以正常启动了.
