# prometheus-网络探测blackbox-exporter关于probe_success与up关系的研究.md

参考文章

1. [Difference between probe_success and up?](https://stackoverflow.com/questions/51984837/difference-between-probe-success-and-up)
    - 有人提问probe_success探测结果与prometheus抓取的状态结果的关系
2. [Prometheus官方文档 Failed scrapes](https://prometheus.io/docs/instrumenting/writing_exporters/#failed-scrapes)
3. [Checking for HTTP 200s with the Blackbox Exporter](https://www.robustperception.io/checking-for-http-200s-with-the-blackbox-exporter)
    - 综合使用`probe_success`与`up`进行告警配置的示例.
4. [blackbox-exporter官方issue: incorrect state](https://github.com/prometheus/blackbox_exporter/issues/152)
    - 使用blackbox_exporter在命令行发送监测请求的示例


按照上篇文章所写的方式部署black-exporter, 用以监控集群中所有service的健康状态(主要是业务层的service, 用端口指代pod中的服务是否正常没毛病吧?), 但是实验时却发生了很奇怪的事情.

我首先创建了一个名为`general-test`的命名空间, 然后部署了两个服务(同时包含deploy与svc), 分别以后缀01和02结尾. 基本配置如下

```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-svc-01
  labels:
    app: nginx-svc-01
  namespace: general-test
spec:
  ports:
    - port: 80
      name: http
      targetPort: 80
  selector:
    ## 注意: service 的 selector 需要指定的是
    ## Deployment -> spec -> template -> labels,
    ## 而不是 Deployment -> metadata -> lables.
    app: nginx-pod-01
## Service end ...

---
apiVersion: apps/v1
kind: Deployment
metadata:
  ## deploy 生成的 pod 的名称也是 nginx-deploy-xxx
  name: nginx-deploy-01
  labels:
    app: nginx-deploy-01
  namespace: general-test
spec:
  replicas: 1
  selector:
    matchLabels:
      ## 这里的 label 是与下面的 template -> metadata -> label 匹配的,
      ## 表示一种管理关系
      app: nginx-pod-01
  template:
    metadata:
      labels:
        app: nginx-pod-01
    spec:
      containers:
      - name: nginx
        image: nginx
        imagePullPolicy: IfNotPresent
## Deployment end ...
```

`general-test`命名空间下的资源如下

```console
$ k get svc -n general-test
NAME           TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
nginx-svc-01   ClusterIP   10.100.180.6     <none>        80/TCP    6m45s
nginx-svc-02   ClusterIP   10.107.230.240   <none>        80/TCP    6m25s
$ k get deploy -n general-test
NAME              READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deploy-01   1/1     1            1           7m17s
nginx-deploy-02   1/1     1            1           6m58s
```

对应的, prometheus的job配置如下(基本上参考prometheus-book的, 只加了对命名空间的限制, 缩小过滤范围, 且这里没有对状态码有其他配置, 因为本文的关注点不在那里)

```yaml
scrape_configs:
  - job_name: kubernetes-services
  ## 使用自动服务发现监控service需要部署 blackbox-exporter 服务.
  ## 见: https://yunlzheng.gitbook.io/prometheus-book/part-iii-prometheus-shi-zhan/readmd/use-prometheus-monitor-kubernetes#dui-ingress-he-service-jin-hang-wang-luo-tan-ce
    metrics_path: /probe
    params:
      module: [http_2xx]
    kubernetes_sd_configs:
      - role: service
        namespaces:
          names:
            - general-test
    relabel_configs:
      ## `__address__`原为服务发现时获取到的service对象地址,
      ## 这里将赋值给获取监控数据的请求参数`__param_target`.
      - source_labels: [__address__]
        target_label: __param_target
      - target_label: __address__
        replacement: blackbox-exporter.monitoring.svc.cluster.local:9115
      ## `__param_target`指定了`blackbox_exporter`要探测的目标,
      ## 这里把发现到的service对象中的instance属性值作为目标, 以实现动态检测.
      - source_labels: [__param_target]
        target_label: instance
      - action: labelmap
        regex: __meta_kubernetes_service_label_(.+)
      - source_labels: [__meta_kubernetes_namespace]
        target_label: kubernetes_namespace
      - source_labels: [__meta_kubernetes_service_name]
        target_label: kubernetes_name
  ## job kubernetes-services end ...
```

prometheus可以成功发现这两个服务.

![](https://gitee.com/generals-space/gitimg/raw/master/0FDA910BC573B08E7D0B3A9AC2025612.jpg)

好像可以了?

但是当我试着删掉其中一个deploy - `nginx-deploy-01`, 模拟服务down掉的情况时, 问题来了, prometheus中的抓取结果没有任何变化.

理论上后端的pod服务不存在后, service将会因为端口不通检测失败才对啊, 这样的结果还怎么做监控?

我在网上查了很久, 为blackbox-exporter调整日志层级, 又找到了参考文章4, 在一个调试用的容器中手动触发了探测行为, 得到如下输出

```
$ curl 'blackbox-exporter.monitoring.svc.cluster.local:9115/probe?target=nginx-svc-01.general-test.svc.cluster.local&module=http_2xx&debug=true'
ts=2019-12-12T15:34:02.186635416Z caller=client.go:250 module=http_2xx target=nginx-svc-01.general-test.svc.cluster.local level=info msg="Making HTTP request" url=http://10.100.180.6 host=nginx-svc-01.general-test.svc.cluster.local
ts=2019-12-12T15:34:03.226094362Z caller=main.go:119 module=http_2xx target=nginx-svc-01.general-test.svc.cluster.local level=error msg="Error for HTTP request" err="Get http://10.100.180.6: dial tcp 10.100.180.6:80: connect: connection refused"
ts=2019-12-12T15:34:03.226275415Z caller=main.go:304 module=http_2xx target=nginx-svc-01.general-test.svc.cluster.local level=error msg="Probe failed" duration_seconds=1.044054542

Metrics that would have been returned:
...省略
probe_http_content_length 0
...省略
probe_http_status_code 0
...省略
probe_success 0
```

`probe_success`证明探测失败, 这就说明**探测失败只能表示探测行为本身的成功与否, 不能表示服务本身是否正常...**

于是我又查了查, 发现有一部分人是直接把`probe_success`与`up`两种属性分别对待的. 在prometheus的查询界面, 可以得到如下结果.

![](https://gitee.com/generals-space/gitimg/raw/master/2E8DF9095057236BE80878A5C101125D.png)

可以看到两个服务的`probe_success`的值是不同的.

![](https://gitee.com/generals-space/gitimg/raw/master/47CB11DB165BA58249B4D0EBAFF32D46.jpg)

另外, `probe_success`也有时间序列记录.

------

参考文章3中有人给出了综合使用`probe_success`与`up`进行告警配置的示例, 这里可以借鉴一下.

```yaml
groups:
- name: example
  rules:
   - alert: ProbeFailing
     expr: up{job="blackbox"} == 0 or probe_success{job="blackbox"} == 0
     for: 10m
```

在`up`本身为0, 或是连续探测失败10m后发出告警.

...只是总觉得这种处理机制怪怪的.
