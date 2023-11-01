# helm install报错Error：UPGRADE FAILED：parse error in "" template：bad character U+002D '-'

参考文章

1. [Accessing values of the subchart with dash in the name](https://github.com/helm/helm/issues/2192)

## 问题描述

自定义chart中, template 目录下某个文件, 假如文件名为`deployment.yaml`, 如下

```yaml
key-01: {{ .Values.key-01 }}
```

`values.yaml`的文件内容如下

```yaml
key-01: value01
```

> 注意: `key-01`在yaml的第一层级, 没有父级块, template 中才需要使用`.Values.key-01`来引用;

但在是 helm install/upgrade 时会报错

```log
Error: UPGRADE FAILED: parse error in "deployment.yaml": template: deployment.yaml:217: bad character U+002D '-'
```

问题出在`{{}}`中的`-`符号, `-`在yaml中是有特殊含义的, 所以不能这么用.

## 解决方案

template 引用这种带有`-`符号, 需要使用另一种方法

```yaml
key-01: {{ index .Values "key-01" }}
```
