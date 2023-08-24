参考文章

1. [helm的安装、使用以及自定义Chart](https://blog.csdn.net/u010606397/article/details/112062312)

如下命令可以创建一个包含示例模板目录.

```
helm create xxx-chart
```

### 安装本地 chart

使用`helm fetch`下来chart工程并解压, 查看`values.yaml`哪些可以修改的字段, 写到`myval.yaml`文件中, 然后使用如下命令安装

```
helm install mychart -f myval.yaml ./xxx-chart
```

也可以先到 chart 目录, 再安装

```
cd xxx-chart
helm install mychart -f values.yaml ./
```

## 移除

```
$ helm uninstall mychart
release "mychart" uninstalled
```
