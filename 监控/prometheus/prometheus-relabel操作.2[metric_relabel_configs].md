# prometheus-relabel操作.2[metric_relabel_configs]

`metric_relabel_configs`与`scrape_configs`同级, ta的配置块规则与后者相同, 都是对 label 标签进行增加/过滤/替换/拷贝的功能.

假设当 prometheus 中采集的指标中, 存在 {pod="pod名称", container="container名称"} 字段, 但是某些项目是根据该指标的"pod_name"和"container_name"进行处理的.

此时为了保证兼容性, 可以将"pod"标签复制一份为"pod_name", 把"container"标签复制成"container_name".

```yaml
  scrape_configs:
  ## ...省略
  metric_relabel_configs:
  - source_labels: [pod]
    separator: ;
    regex: (.+)
    target_label: pod_name
    replacement: $1
    action: replace
  - source_labels: [container]
    separator: ;
    regex: (.+)
    target_label: container_name
    replacement: $1
    action: replace
```
