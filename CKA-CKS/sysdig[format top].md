参考文章

1. [使用Sysdig监测你的容器](https://www.51cto.com/article/677972.html)
    - sysdig直接从Linux 内核(而不是/proc)收集系统调用和事件，并(自行)执行strace、tcpdump、htop、iftop、lsof和Wireshark所做的工作

`sysdig`, 类似`strace`获取目标进程/容器的系统调用行为, 不过前者输出的信息要更详细.

## 获取容器中的所有事件

```console
controlplane $ crictl ps
CONTAINER        IMAGE            CREATED        STATE      NAME              ATTEMPT    POD ID           POD
c462b68c60170    deb04688c4a35    2 weeks ago    Running    kube-apiserver    2          c6b23d817f742    kube-apiserver-controlplane
```

```log
controlplane $ sysdig container.name=kube-apiserver | head -n5 
15 13:30:15.099243000 0 container:c462b68c6017 (-1) > container json={"container":{"Mounts":[],"cpu_period":100000,"cpu_quota":0,"cpu_shares":51,"cpuset_cpu_count":0,"env":[],"id":"c462b68c6017","image":"registry.k8s.io/kube-apiserver:v1.26.1","imagedigest":"sha256:xxxxxx","imageid":"xxxxxx","imagerepo":"registry.k8s.io/kube-apiserver","imagetag":"v1.26.1","ip":"0.0.0.0","is_pod_sandbox":false,"labels":{},"memory_limit":0,"metadata_deadline":0,"name":"kube-apiserver","port_mappings":[],"privileged":false,"swap_limit":0,"type":7}}

22 13:30:15.130347873 0 kube-apiserver (33463) < nanosleep res=0 
23 13:30:15.130353119 0 kube-apiserver (33463) > switch next=56625(sysdig) pgft_maj=0 pgft_min=1 vm_size=1119196 vm_rss=303352 vm_swap=0 
25 13:30:15.130401884 0 kube-apiserver (33463) < nanosleep res=0 
```

> `container.id`可能无法实现对容器的筛选, 会卡住, 需要使用`container.name`.

```
$ sysdig --help
## ...省略
Output format:

By default, sysdig prints the information for each captured event on a single
 line with the following format:

 %evt.num %evt.outputtime %evt.cpu %proc.name (%thread.tid) %evt.dir %evt.type %evt.info
```

- evt.num: 所谓的事件序号, 是递增的, 但不一定连续;
- evt.time: 事件的时间戳
- evt.cpu: 当前事件发生时, 所在的cpu id;
- proc.name: 当前事件所在的进程名, 类似于 top 中的`COMMAND`列, 或是`ps -ef`的`CMD`列, 没有参数.
    - 如果筛选目标是容器, 则会有`container:${容器id}`
- thread.tid: 线程号, 如果进程只有一个主线程, 则与`proc.pid`相等;
- evt.dir: 事件方向, `>`表示进入一个事件, `<`表示退出事件(应该是分别只事件开始与结束吧?);
- evt.type: 事件类型, 如`open`, `read`, `sleep`等;
- evt.info is the list of event arguments.

## -l 获取所有过滤条件(可展示列名)

```console
$ sysdig -l

Field Class: process
proc.pid        the id of the process generating the event.
proc.exe        the first command line argument (usually the executable name or a custom one).
proc.name       the name (excluding the path) of the executable generating the event.
proc.args       the arguments passed on the command line when starting the process generating the event.

Field Class: evt
evt.num         event number.
evt.time        event timestamp as a time string that includes the nanosecond part.

Field Class: container
container.id    the container id.
container.name  the container name.
container.image the container image name (e.g. sysdig/sysdig:latest for docker).

Field Class: k8s
k8s.pod.name    Kubernetes pod name.
k8s.pod.id      Kubernetes pod id.
k8s.pod.label   Kubernetes pod label. E.g. 'k8s.pod.label.foo'.
```

每个filed, 都可以作为过滤条件, 也可以作为输出结果的列. 如

```
sysdig container.name=kube-apiserver -p'%container.image'
```
