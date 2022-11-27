# kuber-ConfigMap认识

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)

2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

## 1. 整体认识

`configMap`与`pod`, `service`等概念平级.

可以使用`configMap`对象将命令行中指定的键值对, 或配置文件里的内容, 或是目录下的文件存储在kubernetes集群中, 作为配置字典(...还是键值对).

pod容器可以引用`configMap`中的键值对作为环境变量, 启动命令, 甚至挂载卷的一部分来使用.

参考文章1, 2中介绍了它的详细使用方法, 不过写得太走心, 很随性...

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
