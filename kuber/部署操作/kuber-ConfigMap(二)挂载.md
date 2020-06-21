# kuber-ConfigMap认识

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)

2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

我们把配置字段读取到configMap对象中, 那怎么使用呢?

在pod的配置文件中, 可以用configMap对象中的字段作为环境变量(env).

我们还能把配置文件读进去, 所以我们还可以把configMap当作目录一样挂载. 但是通过configMap挂载的目录都是只读的, 如果需要写权限, 请使用`hostPath`.

## 1. Pod配置文件中引用

1. 环境变量的引用

```yml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
    - name: test-container
      image: nginx
      env:
        # 为Pod定义的环境变量, 它的值将使用configMap中name字段的值
        - name: NAME
          valueFrom:
            configMapKeyRef:
              # configMap的名称
              name: myconfig1
              # configMap中的name字段
              key: name
```

2. 直接把configMap对象作为环境变量配置文件

```yml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
    - name: test-container
      image: nginx
      envFrom:
      - configMapRef:
          name: myconfig1
```

按照参考文章2所说, 这种方式要求目标config map对象必须为键值类型. 即, 要么用`--from-literal=name=general`手动赋值, 要么用`--from-env-file=config.cfg`读入键值文件.

## 2. 挂载configMap为volume卷

如果把configMap作为卷挂载到容器中的某个目录下, 则**该configMap的键名都会成为文件名, 而它们的值则都会成为对应文件的内容**.

以前面`myconfig6`对象为例. 创建`pod_configmap.yml`文件

```yml
---
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  volumes:
  - name: myconfig
    configMap:
      ## 挂载myconfig6对象
      name: myconfig6

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

```
$ kubectl create -f pod_configmap.yml 
pod "test-pod" created
```

进去看看.

```
$ kubectl exec -it test-pod -- /bin/bash
[root@test-pod /]# cd /root/config
[root@test-pod config]# ls
config1.cfg  sex
[root@test-pod config]# ls -al
total 12
drwxrwxrwx 3 root root 4096 Apr 20 04:13 .
dr-xr-x--- 1 root root 4096 Apr 20 04:13 ..
drwxr-xr-x 2 root root 4096 Apr 20 04:13 ..2018_04_20_04_13_27.897204004
lrwxrwxrwx 1 root root   31 Apr 20 04:13 ..data -> ..2018_04_20_04_13_27.897204004
lrwxrwxrwx 1 root root   18 Apr 20 04:13 config1.cfg -> ..data/config1.cfg
lrwxrwxrwx 1 root root   10 Apr 20 04:13 sex -> ..data/sex
[root@test-pod config]# cat sex
male
[root@test-pod config]# cat config1.cfg 
name=general
age=21
[root@test-pod config]# touch test
touch: cannot touch 'test': Read-only file system
```

挂载后`config1.cfg`和`sex`都是文件了...

------

有时我们希望映射到容器中的文件名改一下, 比如`config1.cfg` -> `profile`, 可以把pod配置写成如下这种.

```yml
---
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  volumes:
  - name: myconfig
    configMap:
      ## 挂载myconfig6对象
      name: myconfig6
      items:
      - key: config1.cfg
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
