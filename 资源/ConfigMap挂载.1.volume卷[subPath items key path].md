# ConfigMap挂载.1.volume卷

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)
2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

ConfigMap中的数据(data块), 本质上是**键值对**, 每个键值对对应一个配置文件. 而作为卷挂载时, 有两种方式

1. 将整个ConfigMap挂载到一个目录, 每个键值对就是该目录下的一个配置文件(键名都会成为文件名, 而它们的值则都会成为对应文件的内容);
2. 将ConfigMap各个键值对分别挂载为独立的文件;

假设ConfigMap对象如下

```yaml
apiVersion: v1
data:
  config.cfg: |
    name=general
    age=21
  sex: male
kind: ConfigMap
metadata:
  name: myconfig
```

## 挂载ConfigMap(全部或部分键)到目录

创建`pod_configmap.yml`文件, 并挂载ta

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  volumes:
  - name: myconfig
    configMap:
      ## 挂载myconfig对象
      name: myconfig

  containers:
  - name: test-container
    image: centos
    command: ["tail", "-f", "/etc/profile"]
    volumeMounts:
    - name: myconfig
      ## 在这一示例中, mountPath需要是一个目录, 如果这个目录不存在会自动创建为目录,
      ## configmap 中的键会作为文件写到这个目录下.
      ## 如果这个路径存在但却是个文件, 那么容器启动就会出错. 如下
      ## unknown: Are you trying to mount a directory onto a file (or vice-versa)? 
      ## Check if the specified host path exists and is the expected type
      mountPath: /root/config
```

```console
$ kubectl create -f pod_configmap.yml 
pod "test-pod" created
```

进去看看.

```console
$ kubectl exec -it test-pod -- /bin/bash
[root@test-pod /]# cd /root/config
[root@test-pod config]# ls
config.cfg  sex
[root@test-pod config]# ls -al
total 12
drwxrwxrwx 3 root root 4096 Apr 20 04:13 .
dr-xr-x--- 1 root root 4096 Apr 20 04:13 ..
drwxr-xr-x 2 root root 4096 Apr 20 04:13 ..2018_04_20_04_13_27.897204004
lrwxrwxrwx 1 root root   31 Apr 20 04:13 ..data -> ..2018_04_20_04_13_27.897204004
lrwxrwxrwx 1 root root   18 Apr 20 04:13 config.cfg -> ..data/config.cfg
lrwxrwxrwx 1 root root   10 Apr 20 04:13 sex -> ..data/sex
[root@test-pod config]# cat sex
male
[root@test-pod config]# cat config.cfg 
name=general
age=21
```

挂载后`config.cfg`和`sex`都是文件了...

ConfigMap映射的目录是只读的, 不可再创建其他文件.

```console
[root@test-pod config]# touch test
touch: cannot touch 'test': Read-only file system
```

如果目标目录不存在, 会自动创建.

如果挂载的目标目录是容器内的一个已经存在的目录, 则会将其覆盖, 该目录下原本的文件会丢失.

------

有时我们希望映射到容器中的文件名改一下, 比如`config.cfg` -> `profile`, 可以把pod配置写成如下这种.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  volumes:
  - name: myconfig
    configMap:
      ## 挂载myconfig对象
      name: myconfig
      items:
      - key: config.cfg
        path: profile

  containers:
  - name: test-container
    image: centos
    command: ["tail", "-f", "/etc/profile"]
    volumeMounts:
    - name: myconfig
      mountPath: /root/config
```

...md, 用`items`字段需要显示定义所有键, 没定义的就不映射了. 上面的pod创建出来, `/root/config`下就只有一个`profile`文件.

注意: 这种映射方式仍然会覆盖容器内的`/root/config`目录.

## 挂载ConfigMap中的指定键到文件(subPath)

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  volumes:
  - name: myconfig
    configMap:
      ## 挂载myconfig对象
      name: myconfig

  containers:
  - name: test-container
    image: centos
    command: ["tail", "-f", "/etc/profile"]
    volumeMounts:
    - name: myconfig
      mountPath: /root/config1/myname
      subPath: name
    - name: myconfig
      mountPath: /root/config2/myage
      subPath: age
```

这种方法会分别将ConfigMap中的`name`与`age`键映射到`/root/config1/myname`和`/root/config2/myage`文件, 互不干扰.

而且这种方式不会覆盖原目录, 如果`mountPath`是`/etc/myname`, 那么在将`name`文件映射到`/etc/myname`的同时, `/etc`目录下的文件也不会被清空.

> 本节的`subPath`可以与上节的`items`共同使用, `subPath`的值应该与`items[].path`一致即可...不过没多大意义就是了.

