参考文章

1. [helm的安装、使用以及自定义Chart](https://blog.csdn.net/u010606397/article/details/112062312)

如下命令可以创建一个包含示例模板目录.

```
helm create xxx-chart
```

### 安装本地 chart

使用`helm fetch`下来chart工程并解压, 查看`values.yaml`哪些可以修改的字段, 写到`myval.yaml`文件中, 然后使用如下命令安装

```
helm install kubeapps -f myval.yaml ./xxx-chart
```
