# prometheus-relabel操作

参考文章

1. [Prometheus学习系列（十三）之配置解析](https://www.jianshu.com/p/fb5c82de935d)
    - `<relabel_config>`: replace, keep, drop, hashmod, labelmap, labeldrop, labelkeep
2. [Prometheus中relabel_configs的使用](https://www.li-rui.top/2019/04/16/monitor/Prometheus%E4%B8%ADrelabel_configs%E7%9A%84%E4%BD%BF%E7%94%A8/)
    - 对`relabel_configs`及其子字段的概念解释得很清晰

## 1. 使用场景

prmetheus从服务暴露出的`/metrics`接口中获得的进程状态信息, 通常会包含多种维度(或者说字段). labels这个概念, 表示不同字段的信息. 因为我们总希望来源信息尽可能丰富, 后续可以通过一些手段将不需要的字段过滤, 只显示我们需要的部分. 

这个过滤的过程, 被称为`relabel`. 

需要注意的是, relabel操作不仅可以修改字段的键, 也可以修改对应的值, 十分灵活且强大.

如下图

![](https://gitee.com/generals-space/gitimg/raw/master/CD13204C0C8EBC0C93015F9F86851AFE.png)

黑色气泡框中是prometheus从apiserver获取的所有指标种类, 在经过relabel后, 表格中就只会显示其中两个字段.

```yaml
    scrape_configs:
    ## 这个名称为表格头部蓝色字字的`kubernetes-apiservers (1/1 up)`部分的名称.
    - job_name: 'kubernetes-apiservers'
      kubernetes_sd_configs:
        - role: endpoints
      # Keep only the default/kubernetes service endpoints for the https port. This
      # will add targets for each API server which Kubernetes adds an endpoint to
      # the default/kubernetes service.
      relabel_configs:
      - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
        action: keep
        regex: default;kubernetes;https
```

## 2. 使用方法

### 2.1 replace替换

`replace`操作需要有`source_labels`和`target_label`作为来源和目标匹配, `regex`只是匹配`source_labels`的手段, 其中的分组引用可以用在`replacement`字段中.

```yaml
relabel_configs:
- source_labels: [__meta_kubernetes_node_name]
  regex: (.+)
  target_label: __metrics_path__
  replacement: /api/v1/nodes/${1}/proxy/metrics
```

当然, `replace`操作可以不匹配`source_labels`, 而直接指定`target_label`和`replacement`完成替换.

```yaml
relabel_configs:
- target_label: __address__
  replacement: kubernetes.default.svc:443
```

> 默认操作为`replace`(不指定`action`时, 即为`replace`).

### 2.2 keep/drop

与`replace`操作相比, `keep`, `drop`只需要`source_labels`字段, 无需`target_label`.

- keep: 删除`source_labels`中列表举出的`label`值与`regex`表达式不匹配的job结果; 
- drop: 删除`source_labels`中列表举出的`label`值与`regex`表达式匹配的job结果. 

可以作为过滤器使用, 过滤一下同命名空间下的同类型资源(监控的5种模式, ep, node那种).

### 2.3 labelmap/labelkeep/labeldrop

labelmap, labeldrop, labelkeep对于数据的采集与展示是没有影响的, ta们只影响存储到时间序列中的维度(是否存在, 或是命名).

在prometheus的WebUI中, Status -> Target面板下, 各Job的表格展示中, `Labels`列表示的就是`labelmap`, `labeldrop`与`labelkeep`处理后的结果, 鼠标悬停在上面可以看到`Before Relabeling`时的数据.

------

应用场景, 看起来像是将内置的 label 字段简化, 移除`__meta`前缀.

```yaml
  ## labelmap 这个行为比较特殊
  - separator: ;
    regex: __meta_kubernetes_service_label_(.+)
    replacement: $1
    action: labelmap
```
