## ConfigMap创建方法.1.--from-file普通文件

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)
2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

基本语法

```
kubectl create configmap configMap名称 数据源
```

其中, 数据源就是我们上面提到的, 命令行传入的键值对, 或配置文件的内容, 或是一个目录下的文件.

## 1. 从文件读取

**文件的格式没有限制**, 假设文件名为`config.cfg`, 其内容如下.

```ini
name=general
age=21
```

加载方式为

```
kubectl create configmap myconfig --from-file=config.cfg
```

用这两种方式创建的`configMap`不太一样. 比较一下

```yaml
## kubectl get configmap myconfig -o yaml
apiVersion: v1
data:
  config.cfg: |
    name=general
    age=21
kind: ConfigMap
metadata:
  name: myconfig
```

data块中存储的其实也是键值对, 键为`config.cfg`, 为读入的文件名, 值则为文件内容, 这个值其实相当于是一个长字符串了.

### 1.1 自定义读入的键名

上面使用`--from-file`选项创建的configMap对象, `config.cfg`这个键的名称是可以自定义的, 方法如下

```
kubectl create configmap myconfig --from-file=mykey=config.cfg
```

查看一下

```yaml
## kubectl get configmap myconfig -o yaml
apiVersion: v1
data:
  mykey: |
    name=general
    age=21
kind: ConfigMap
metadata:
  name: myconfig
```

其中的键名不再是文件名`config.cfg`而是我们自定义的键名`mykey`.

## 1.2 从目录创建

`--from-file`可以指定一个目录, 这样得到的configMap将会是以各个文件名为键名, 文件内容为其对应的值. 

假设存在如下目录

```log
$ ls config_dir/
config1.cfg	config2.cfg
## 两个文件的内容一样
$ cat config_dir/*
name=general
age=21
name=general
age=21
```

```yaml
## kubectl create configmap myconfig --from-file=config_dir
## kubectl get configmap myconfig -o yaml
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
  name: myconfig
```

正如所料.

注意: **不支持子目录**, 目标目录下的子目录是不会读取的, 所以 data 下只会存在单层的键值对.
