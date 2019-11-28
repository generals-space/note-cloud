# helm模板语义总结

参考文章 

1. [Chart Development Tips and Tricks](https://helm.sh/docs/howto/charts_tips_and_tricks/)

2. [Built-in Objects](https://helm.sh/docs/topics/chart_template_guide/builtin_objects/)
    - helm chart内置对象, 如`.Release.Name`, `.Release.Namespace`
    - `.Release`表示helm创建的release实例
    - `.Chart`表示`Chart.yaml`文件中的内容空间.
    - `.Values`同理, 表示`values.yaml`内容空间.

3. [helm--chart模板文件简单语法使用](https://www.cnblogs.com/DaweiJ/articles/8779256.html)
    - 用户级试用总结, 更易上手.

go模板中有`template`指令, 作用类似于c语言中的`include`, 可以把其他模板文件中的内容加载到当前文件. 但ta不能在管道操作中使用, 比如 `{{ template "mysql.fullname" . }}` 可以理解, 但是`{{ template "mysql.fullname" | uppper }}` 就不行了.

于是helm扩展`template`指令为`include`指令, 使其可以实现在管道中使用. 见参考文章1.

------

go template注释

```
{{/*
...
*/}}
```

------

`{{ template "mysql.fullname" . }}`: `mysql.fullname`表示使用`helm install xxx stable/mysql`中的名称xxx.

`{{ .Release.Name }}`: `.Release.Name`与fullname一样.

------

`toYaml`函数

`{{ toYaml .Values.resources | indent 10 }}`将目标内容仍然以yaml格式导入, 并修改其缩进空间, 与include/template应该还是有所不同的. 

注意, 只是导入`.Values.resources`下的内容, 而不包括resource字段本身. 如模板中是这么写的

```yaml
        resources:
{{ toYaml .Values.metrics.resources | indent 10 }}
```

而values.yaml文件则是这样的

```yaml
resources:
  requests:
    memory: 256Mi
    cpu: 100m
```

看到了?

------

`$Chart/templates/_helpers.tpl`文件中定义了release实例的名称生成规则, 在指定了`--generate-name`选项时将自动生成随机名称. 此模板采用的是`mustache`格式(文件第一行就可以看出).
