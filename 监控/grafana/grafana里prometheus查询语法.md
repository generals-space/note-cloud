# grafana里prometheus查询语法

参考文章

1. [grafana里prometheus查询语法](https://www.cnblogs.com/jonnyan/p/10410614.html)
2. [Grafana文档（Prometheus）](https://segmentfault.com/a/1190000016237454)

prometheus 中查询指定指标的语法大致为`metric{label=value}`, 但是 grafana 里不是这样的, 貌似是有一套独立的语法, 参考文章1中给出了几个示例.

## label_values(metric, label)

查询 prometheus `metric`指标项中, 拥有`label`字段的所有结果, 且只查询`label`字段的值, 不包含结果中的其他字段. 如果对应到 prmetheus 自身的查询语句, 应该是

```
metric{label!=""}
```

但是需要注意, grafana 的 `label_values(metric,label)`只查询`label`值, 结果是一个`lebel`值的列表. 

![](https://gitee.com/generals-space/gitimg/raw/master/02175710666c5ddeee1b8808b76bee5f.png)

而 prometheus 中的`metric{label!=""}`查询出的则是`metric`指标的结果列表.

![](https://gitee.com/generals-space/gitimg/raw/master/c56aaf8a44b8039b755a16882f909245.png)

我找了找, 没有找到在 prometheus 中与`label_values`相似功能的函数. 大概是因为 prometheus 本身只关注指标结果, `label`字段只是作为过滤手段, 所以没有考虑只显示`label`值的情况.

## query_result(query)

这是除了上面的`label_values()`方法后我觉得最有用的 grafana 表达式了, `query_result()`函数中, 可以直接写 prometheus 的查询语句.
