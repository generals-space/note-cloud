# ConfigMap挂载.2.env环境变量[envFrom]

参考文章

1. [Configuring Redis using a ConfigMap](https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/)
2. [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume)

与挂载为 volume 类型, 挂载ConfigMap为环境变量时, 也有2种方式:

1. 将整个ConfigMap挂载到env;
2. 挂载ConfigMap中指定key到某个环境变量;

以如下ConfigMap为例.

```yaml
apiVersion: v1
data:
  name: general
  age: "21"
  content: |-
    hello world
    hello kitty
    hello kugou
kind: ConfigMap
metadata:
  name: myconfig
```

## 1. 把configMap对象作为环境变量配置文件

```yaml
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
        name: myconfig
```

## 2. 挂载部分环境变量

```yaml
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
          name: myconfig
          # configMap中的name字段
          key: name
    - name: CONTENT
      valueFrom:
        configMapKeyRef:
          # configMap的名称
          name: myconfig
          # configMap中的name字段
          key: name
```

在挂载为环境变量时, `CONTENT`的内容会变成单行, 如下

```
$ echo $CONTENT 
hello world hello kitty hello kugou
```
