# kubernetes中configMap的认识

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)

2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

## 1. 整体认识

`configMap`与`pod`, `service`等概念平级.

可以使用`configMap`对象将命令行中指定的键值对, 或配置文件里的内容, 或是目录下的文件存储在kubernetes集群中, 作为配置字典(...还是键值对).

pod容器可以引用`configMap`中的键值对作为环境变量, 启动命令, 甚至挂载卷的一部分来使用.

参考文章1, 2中介绍了它的详细使用方法, 不过写得太走心, 很随性...

总之, `configMap`对象的最终使用形态就是键值对!

## 2. configMap的创建方法

基本语法

```
kubectl create configmap configMap名称 数据源
```

其中, 数据源就是我们上面提到的, 命令行传入的键值对, 或配置文件的内容, 或是一个目录下的文件.

### 2.1. 从命令行指定键值对

```
$ kubectl create configmap myconfig1 --from-literal=name=general --from-literal=age=21
configmap "myconfig1" created
```

查看一下

```
$ kubectl get configmap 
NAME        DATA      AGE
myconfig1   2         3s
```

### 2.2 从文件读取

文件的格式没有限制, 假设文件名为`config.cfg`, 其内容如下.

```ini
name=general
age=21
```

加载方式为

```
$ kubectl create configmap myconfig2 --from-file=config.cfg
```

用这两种方式创建的`configMap`不太一样. 比较一下

```yml
$ kubectl get configmap myconfig1 -o yaml
apiVersion: v1
data:
  age: "21"
  name: general
kind: ConfigMap
metadata:
  name: myconfig1

$ kubectl get configmap myconfig2 -o yaml
apiVersion: v1
data:
  config.cfg: |
    name=general
    age=21
kind: ConfigMap
metadata:
  name: myconfig2
```

可以看到两个configMap对象的data结构不同, 前者有两个键: `age`和`name`, 后者只有一个键`config.cfg`, 为读入的文件名, 其值即为文件内容, 但这个值其实相当于是一个长字符串了.

### 2.2.1 `--from-env-file`将文件内容直接读取为键值对

也许你会想从文件中创建和直接用键值对创建相同的configMap, 不想以文件名为键, 文件内容为值, 而是直接把文件中的键值对读取出来, 你可以使用`--from-env-file`这个选项.

但是它要求读入的文件必须符合`key=val`这种格式, 而且可以使用`#`作为注释, 读取时会忽略注释行.

正好我们上面的`config.cfg`符合这个要求, 来试试

```
$ kubectl create configmap myconfig3 --from-env-file=config.cfg
```

```
$ kubectl get configmap myconfig3 -o yaml
apiVersion: v1
data:
  age: "21"
  name: general
kind: ConfigMap
metadata:
  name: myconfig3
```

这样, 就能直接把文件中的`name`和`age`键值对读取出来了. 完美!

### 2.2.2 自定义读入的键名

上面使用`--from-file`选项创建的configMap对象, `config.cfg`这个键的名称是可以自定义的, 方法如下

```
kubectl create configmap myconfig4 --from-file=mykey=config.cfg
```

查看一下

```yml
$ kubectl get configmap myconfig4 -o yaml
apiVersion: v1
data:
  mykey: |
    name=general
    age=21
kind: ConfigMap
metadata:
  name: myconfig4
```

其中的键名不再是文件名`config.cfg`而是我们自定义的键名`mykey`.

### 2.3 从目录创建

`--from-file`可以指定一个目录, 这样得到的configMap将会是以各个文件名为键名, 文件内容为其对应的值. 

假设存在如下目录, 

如下

```
$ ls config_dir/
config1.cfg	config2.cfg
## 两个文件的内容一样, 复制上面的`config.cfg`
$ cat config_dir/*
name=general
age=21
name=general
age=21
```

```yml
$ kubectl create configmap myconfig5 --from-file=config_dir
$ kubectl get configmap myconfig5 -o yaml
apiVersion: v1
data:
  config1.cfg: |
    name=general
    age=21
  config2.cfg: |
    name=general
    age=21
kind: ConfigMap
metadata:
  name: myconfig5
```

正如所料.

> 注意: 不支持子目录, 目标目录下的子目录是不会读取的.

### 2.4 混合使用

`--from-literal`, `--from-file`可以组合使用, 还可以多次使用.

```
$ kubectl create configmap myconfig6 --from-literal=sex=male --from-file=config_dir/config1.cfg --from-env-file=config_dir/config2.cfg
$ kubectl get configmap myconfig6 -o yaml
apiVersion: v1
data:
  config1.cfg: |
    name=general
    age=21
  sex: male
kind: ConfigMap
metadata:
  name: myconfig6
```

不过`--from-env-file`不能和前两个同时使用, 而且也只能加载一次(多次加载不会报错, 但只有最后一个文件会生效)

> error: from-env-file cannot be combined with from-file or from-literal

## 3. configMap的使用方法

我们把配置字段读取到configMap对象中, 那怎么使用呢?

在pod的配置文件中, 可以用configMap对象中的字段作为环境变量(env).

> 官方文档中只写了pod, 我觉得service, deployment等对象的配置也是一样的道理.

我们还能把配置文件读进去, 所以我们还可以把configMap当作目录一样挂载. 但是通过configMap挂载的目录都是只读的, 如果需要写权限, 请使用`hostPath`.

### 3.1 Pod配置文件中引用

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

直接把configMap对象作为环境变量配置文件

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

### 3.2 挂载configMap为volume卷

如果把configMap作为卷挂载到容器中的某个目录下, 则该configMap的键名都会成为文件名, 而它们的值则都会成为对应文件的内容.

以上面`myconfig6`对象为例. 创建`pod_configmap.yml`文件

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

`config1.cfg`和`sex`都是文件了...

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