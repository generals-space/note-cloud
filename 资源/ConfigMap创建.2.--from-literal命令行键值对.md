# ConfigMap创建.2.--from-literal命令行键值对

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)
2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

## 1. 从命令行指定键值对

```console
kubectl create configmap myconfig --from-literal=name=general --from-literal=age=21
```

查看一下

```yaml
## kubectl get configmap myconfig -o yaml
apiVersion: v1
data:
  age: "21"
  name: general
kind: ConfigMap
metadata:
  name: myconfig
```

## 2. 混合使用

`--from-literal`, `--from-file`可以组合使用, 还可以多次使用.

```console
kubectl create configmap myconfig --from-literal=sex=male --from-file=config_dir/config.cfg --from-env-file=config_dir/config2.cfg
```

```yaml
## kubectl get configmap myconfig -o yaml
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

不过`--from-env-file`不能和前两个同时使用, 会报错

```log
error: from-env-file cannot be combined with from-file or from-literal
```

而且也只能加载一次(多次加载不会报错, 但只有最后一个文件会生效).
