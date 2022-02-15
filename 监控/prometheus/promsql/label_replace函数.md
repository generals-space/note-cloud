# label_replace函数

参考文章

1. [Prometheus label_replace()的使用](https://blog.csdn.net/qq_17303159/article/details/111045149)
2. [官方文档 label_replace()](https://prometheus.io/docs/prometheus/latest/querying/functions/#label_replace)

函数原型

```
label_replace(v instant-vector, dst_label string, replacement string, src_label string, regex string)
```

将目标v指标中的`src_label`字段, 更名为`dst_label`字段, 类似于 prometheus job 中的 relabel 方法, 更通俗点说, 就是**给一个 map 对象中的某个 key 改个名字**.

- `v(instant-vector)`: 即一个指标记录(可以看成是一个结构体), 其中包含多个字段, 常见的有: job="xxx", namespaces="yyy";
- `regex`: 也许`src_label`字段所表示的value并不是全部都有用的, 那么可以使用此参数从`src_label`中截取一部分出来;
- `replacement`: 将`regex`从`src_label`字段值取出的内容, 赋为`dst_label`的值;

假设prometheus采集到的数据如下

```
kube_pod_container_resource_limits_cpu_cores{job="k8s-docker", namespace="zjjpt-logstash", node="hua-dlzx1-d1505-gyt", pod="applog-0"} 10
```

使用如下语句进行处理

```
label_replace(kube_pod_container_resource_limits_cpu_cores, "pod_name", "$1", "pod", "(.*)")
```

可以得到如下结果

```
kube_pod_container_resource_limits_cpu_cores{job="k8s-docker", namespace="zjjpt-logstash", node="hua-dlzx1-d1505-gyt", pod="applog-0", pod_name="applog-0"} 10
```

多了一个"pod_name"字段, ta的值和原来的"pod"字段是一致的.

## 应用场景

label重命名, 应用于2个指标相除的场景, 要求两条记录拥有相同的字段, 如下

```
sum by(pod_name, namespace) (rate(container_cpu_usage_seconds_total[2m])) 
/ 
sum by(pod_name, namespace) 
(
    label_replace(
        kube_pod_container_resource_limits_cpu_cores, "pod_name", "$1", "pod", "(.*)"
    )
) 
* 100
```

上述语句求的是各Pod的实时CPU使用率相较于limits值的占比, 而"container_cpu_usage_seconds_total"指标中只有"pod_name", 没有"pod"字段.

