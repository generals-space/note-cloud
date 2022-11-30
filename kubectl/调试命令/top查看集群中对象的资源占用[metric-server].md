# top查看集群中对象的资源占用

参考文章

1. [从kubectl top看K8S监控](https://www.jianshu.com/p/64230e3b6e6c)
    - 有深度, 值得收藏
    - `metric-server`相比于`heapster`的优势
    - `kube-aggregator`在`metric-server`暴露api过程中的作用: `/apis/metrics.k8s.io`
    - `metric-server`只能提供`cpu`和内存指标, 一般来说作为`HPA`的依据已经足够, 但是如果想根据`qps`或是`5xx`错误数进行判断的话, 需要自定义指标.
    - 监控数据的来源: top子命令 -> metric-server -> kubelet -> cAdvisor -> runc/libcontainer -> cgroup
    - `top pod`内存使用量的计算方法(`top pod`值不包含`pause`容器)
    - `top node`的内存计算逻辑
    - `kubectl top pod`和`docker stats`得到的值不一致的原因.
2. [container-monitor](https://yasongxu.gitbook.io/container-monitor/)
    - 容器监控方案.

kube 版本: 1.16.2 单节点集群, 宿主机配置 4C 8G.

`kubectl top`是基础命令, 但是需要部署配套的组件才能获取到监控值.

- 1.8-: [heapter](https://github.com/kubernetes-retired/heapster/blob/master/deploy/kube-config/standalone/heapster-controller.yaml)
- 1.8+: [metric-server](https://github.com/kubernetes-sigs/metrics-server#deployment)

如果不部署监控服务, 执行`top`命令会失败.

```console
$ k top pod
Error from server (NotFound): the server could not find the requested resource (get services http:heapster:)
```

## 关于 metric-server

metric-server 从 kubelet 收集资源信息, 并通过 apiserver 暴露出来. ta 本来是配合 HPA 横向扩容(貌似也有 VPA 纵向扩容)使用的, 监测到资源占用达到临界就自动扩容.

同时ta的信息也可以通过`kubectl top`访问, 方便调试.

部署方法就不在这里详细说明了, 见另一篇文章(`metrics-server`监控专题), 还真有几个小坑.

## 

到 1.16.2, `top`子命令只支持2种资源: `node`和`pod`, 每种资源3个指标: `CPU`, 内存和硬盘.

```console
$ k top node
NAME            CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
k8s-master-01   394m         9%     3191Mi          41%

$ k top pod
NAME                                    CPU(cores)   MEMORY(bytes)
coredns-67c766df46-jln7w                8m           7Mi
coredns-67c766df46-rrw2w                8m           7Mi
etcd-k8s-master-01                      37m          77Mi
kube-apiserver-k8s-master-01            98m          242Mi
kube-controller-manager-k8s-master-01   55m          38Mi
kube-flannel-ds-amd64-r7dft             2m           10Mi
kube-proxy-8x6ch                        9m           14Mi
kube-scheduler-k8s-master-01            4m           13Mi
metrics-server-7cb9f8f89b-bhzcw         2m           12Mi

$ free -m
              total        used        free      shared  buff/cache   available
Mem:           7802        1367         334         341        6100        5794
Swap:             0           0           0
```

仔细算一下, 你会得到如下几个结论:

1. node打印结果和用`free`的不一样, 其实与`top`的结果也不一样.
2. node的CPU和内存使用并不等于所有 Pod 的总和, 而且差的不少...
3. pod的内存值是其实际使用量, 也是做limit限制时判断oom的依据. 
4. pod的使用量等于其所有 container 的总和, 不包括 pause 容器, 值等于`cadvisr`中的`container_memory_working_set_bytes`指标

------

`k top node|pod`的数据有延迟, 大概30s.
