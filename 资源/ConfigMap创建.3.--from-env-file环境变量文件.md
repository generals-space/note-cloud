# ConfigMap创建.3.--from-env-file环境变量文件

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)
2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

也许你会想从文件中创建和直接用键值对创建相同的configMap, 不想以文件名为键, 文件内容为值, 而是直接把文件中的键值对读取出来, 你可以使用`--from-env-file`这个选项.

但是它要求读入的文件必须符合`key=val`这种格式, 而且可以使用`#`作为注释, 读取时会忽略注释行.

```ini
name=general
age=21
```

正好我们上面的`config.cfg`符合这个要求, 来试试

```
kubectl create configmap myconfig --from-env-file=config.cfg
```

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

`--from-env-file`不能和前两个同时使用, 而且也只能加载一次(多次加载不会报错, 但只有最后一个文件会生效)
