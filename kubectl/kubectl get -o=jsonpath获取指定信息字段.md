# kubectl get -o=jsonpath获取指定信息字段

参考文章

1. [kubectl技巧之通过jsonpath截取属性](https://www.cnblogs.com/tylerzhou/p/11049050.html)
    - 很深入, 也很详细
    - 表格里的`‘`和`’`其实应该是单引号, 见参考文章4, 有原版的表格.
2. [kubectl技巧之通过go-template截取属性](https://www.cnblogs.com/tylerzhou/archive/2004/01/13/11047456.html)
    - 与参考文章1为同一系列, 值得收藏
3. [How do I extract multiple values from kubectl with jsonpath](https://stackoverflow.com/questions/46229072/how-do-i-extract-multiple-values-from-kubectl-with-jsonpath)
    - 采纳答案中, `kubectl get pods`没有指明pod名称, 所以会获取所有pod, 所以`jsonpath`语句中写了`{.items[*]}`
4. [kuber官方文档 JSONPath Support](https://v1-16.docs.kubernetes.io/docs/reference/kubectl/jsonpath/)

- kubectl: 1.16.2
- kuber server: 1.17.3

## 简单概念

`jsonpath="{.spec}"`这种是最简单的方法, 不多说.

其中`.spec`就是取对象的成员属性的方法, 同样的, 还可以使用`["key"]`, 比如`jsonpath="{['spec']}"`, 结果是一样的(很像js).

> 注意: `jsonpath`后面的信息最好使用双引号包裹, 内层使用单引号, 不要反过来用.

特殊符号`@`表示当前对象本身, 但是除了参考文章1中提到的`jsonpath="{@}"`外, 就只有在过滤器中使用了, 我暂没有找到其他的使用方法.

还有一个就是双点号`..`, 比如你可以确定整个对象中某个key是唯一的, 比如`qosClass`, 可以不用管ta在哪个子对象中, 层级有多深, 直接使用`jsonpath="{..qosClass}"`就可以将ta的值直接取出. 当然其实这个key也不一定要唯一, `jsonpath="{..name}"`就可能把`metadata`中的`name`名称, `containers`中的`name`名称, 以及`env`中的`name`名称, 都取出来, 以空格分隔.

## 示例

1. 获取名为`SSH_PASSWORD`的环境变量的值

```
k get pod -o=jsonpath='{.spec.containers[0].env[?(@.name=="SSH_PASSWORD")].value}'
```

2. 获取多个字段, 比如名为`kibanaPass`和`SSH_PASSWORD`这两个环境变量的值, 可以用多个`{}`实现

```
k get pod -o=jsonpath='{.spec.containers[0].env[?(@.name=="kibanaPass")].value}{" "}{.spec.containers[0].env[?(@.name=="SSH_PASSWORD")].value}{"\n"}'
```

> `?()`语法可以作为过滤器, 这里的`@`符号表示当前对象本身, 见参考文章1.

> `{" "}`是将两个字段的结果用空格分隔, 不然会连在一起, 同样`{"\n"}`有一个空行.

3. 如果是同一个对象的多个字段, 可以直接通过逗号`,`进行分隔, 以指定多个key.

```
k get pod -o=jsonpath="{.spec.containers[0]['name', 'image']}"
```

`jsonpath`必须要使用双引号包含, 内部key的名称则用单引号包裹, 如果是反过来`jsonpath='{.spec.containers[0]["name", "image"]}'`, 则会报如下错误

```
error: error parsing jsonpath {.spec.containers[0]["name", "image"]}, invalid array index "name"
```

