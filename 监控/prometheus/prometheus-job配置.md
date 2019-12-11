# prometheus-job配置

参考文章

1. [官方文档 kubernetes_sd_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#kubernetes_sd_config)

prometheus webUI -> Status -> Target 页面显示了配置文件中的抓取任务. 这些job需要在配置文件中声明才可以展示, 同时, 只有声明过的job才可以在 Graph 页通过

## 静态任务

以prometheus本身服务为例, 任务的配置格式如下

```yaml
scrape_configs:
  - job_name: prometheus
    static_configs:
      ## target是检测的目标地址, prometheus可以连接自己本地的9090端口
      - targets:
        - localhost:9090
```

对应的web界面如下

![](https://gitee.com/generals-space/gitimg/raw/master/B4B4C41C791DAC5F9C5E88F104DD5ED7.png)

另外, 只有在`scrape_configs`下声明了job, 才能在 Graph 页面显示属于ta的指标. 如果只配置了prometheus这一个服务, 则只能显示与ta相同的指标.

![](https://gitee.com/generals-space/gitimg/raw/master/CD3743C127C8B155DD4AEF98EC76D373.png)

之所以说这样的配置叫作"静态任务", 是因为通过`static_configs`指定的`target`都是常规的声明方法, 即指定`IP:Port`. 

当kuber集群中项目变多, 难道也要用这样的方法, 指定Service名称以监控程序的运行状态吗? 

当然不需要, 因为prometheus提供了`kubernetes_sd_configs`字段来满足这种需求, ta与`static_configs`是平级的关系.

## 动态任务

以apiserver服务为例

```yaml
scrape_configs:
  - job_name: kubernetes-apiservers
    kubernetes_sd_configs:
      - role: endpoints
        namespaces:
          names: 
            - default
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs: []
  ## job kubernetes-apiservers end 
```

最值得重点讲解的就是`kubernetes_sd_configs`(其中`sd`为`service discovery`, 即服务发现), 按照官方文档所说(参考文章1), 此字段提供了5种监控模式: `node`, `service`, `pod`, `endpoints`和`ingress`, 通过`role`字段指定(看数据类型可知, 可以指定多个`role`), prometheus会获取kuber集群中相应资源的信息.

上面的job中指定的是`endpoints`资源, 且通过`namespaces`进行了限定, 我们先看看`default`空间下的ep资源.

```
$ k get ep
NAME         ENDPOINTS            AGE
kubernetes   192.168.0.101:6443   10d
```

由于结果只有这一个, 所以相应的, 出现在web界面中该job下的结果也只有一个.

![](https://gitee.com/generals-space/gitimg/raw/master/2179BC55A47E7A4EAC56025EC15468CF.png)

上面的`kubernetes_sd_configs`配置中, 通过`namespaces`过滤了多余的ep资源. 但这其实是不准确的, 因为default命名空间下如果存在其他的ep资源, 也会出现在job的结果表格中. 我们可以使用ta进行初级的过滤, 但不能作为最终的配置. 如下图.

![](https://gitee.com/generals-space/gitimg/raw/master/71D3422C7B57B74B9BA3503E9CAA5ED1.png)

为了完成精确的过滤, 就必须使用另一个字段`relabel_configs`. 

### `relabel_configs` 实现过滤

```yaml
scrape_configs:
  - job_name: kubernetes-apiservers
    kubernetes_sd_configs:
      - role: endpoints
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
        action: keep
        regex: default;kubernetes;https
  ## job kubernetes-apiservers end ...
```

上面的配置可达到我们的最终目的, 其重点就在于`relabel_configs`. `action`指定为`keep`, 就是将`source_labels`字段中列举出的各`label`的值(可以在webUI中查看, 黑色气泡框中的`Before Relabeling`中的内容即是), 与`regex`中声明的值不相同的`endpoints`资源移除, 只保留满足`regex`规则的行.

即, 服务发现的`endpoints`资源中, `__meta_kubernetes_namespace`标签的值要为`default`, `__meta_kubernetes_service_name`标签的值要为`kubernetes`, `__meta_kubernetes_endpoint_port_name`标签的值要为`https`, 否则就不能出现在此job结果中.

------

关于`relabel_configs`还有许多其他的内容, 可以参见另一篇文章, 这里不再介绍.
