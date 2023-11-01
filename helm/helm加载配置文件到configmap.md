# helm加载配置文件到configmap

参考文章

1. [Kubernetes - How to define ConfigMap built using a file in a yaml?](https://stackoverflow.com/questions/53429486/kubernetes-how-to-define-configmap-built-using-a-file-in-a-yaml)
    - `.File.Get`获取静态文件内容构建`ConfigMap`对象
2. [Feature request: .env file parsing](https://github.com/helm/helm/issues/4796)
    - 希望helm能够拥有加载env格式的配置文件的方法.

## 

但是helm目前貌似无法加载env文件.

以如下env文件为例

```
POSTGRES_USER=wordpress
POSTGRES_PASSWORD=123456
```

使用 configmap 的`--from-env-file`选项可以将其加载为键值对的模式

```
kubectl create configmap my-env --from-env-file=./env
```

```yaml
data:
  POSTGRES_USER: "wordpress" 
  POSTGRES_PASSWORD: "123456" 
```

然后deploy资源可以使用如下字段引入

```yaml
      containers:
      - name: harbor-db-app
        image: goharbor/harbor-db:v1.8.2
        envFrom:
        - configMapRef:
            name: harbor-db-env
```

但是helm并没有提供将`key=val`的格式转换成`key: val`格式的方法, 所以通过`envFrom`加载环境配置文件是不切实际的...

