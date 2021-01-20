# 容器内设置 sysctl 内核参数[privileged]

参考文章

1. [给容器设置内核参数](https://tencentcloudcontainerteam.github.io/2018/11/19/kernel-parameters-and-container/)
2. [Using sysctls in a Kubernetes Cluster](https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster/)
    - `safe`与`unsafe`的`sysctl`参数
    - `PodSecurityPolicy`资源, 可用于代替`kubelet`启动参数配置可用/禁止的`sysctl`参数???

内核方面做了大量的工作, 把一部分`sysctl`内核参数进行了**namespace化**(namespaced), 即多个容器和主机可以各自独立设置某些内核参数. 例如, 可以通过`net.ipv4.ip_local_port_range`, 在不同容器中设置不同的端口范围. 

运行一个具有`privileged`权限的容器(参考下一节内容), 然后在容器中修改该参数, 然后在宿主机上查询此参数, 看能否看到容器在中所做的修改. 如果看不到, 那就是`namespaced`, 否则不是. 

目前已经namespace化的`sysctl`内核参数：

- `kernel.shm*`
- `kernel.msg*`
- `kernel.sem`
- `fs.mqueue.*`
- `net.*.`

> 注意, 某些参数如`vm.*`并没有namespace化. 比如`vm.max_map_count`, 在主机或者一个容器中设置它, 宿主机及其他所有容器都会受影响(使用最新的值). 

## 1. docker 中设置 sysctl

正常运行的docker容器中, 是不能修改任何sysctl内核参数的. 因为`/proc/sys`是以只读方式挂载到容器里面的. 在容器中执行如下命令

```console
$ mount | grep proc
proc on /proc type proc (ro,nosuid,nodev,noexec,relatime)
```

要给容器设置不一样的sysctl内核参数, 有多种方式. 

### 1.1 `--privileged`

```
docker run --privileged -it ubuntu bash
```

整个`/proc`目录都是以`rw`权限挂载的

```console
$ mount | grep proc
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
```

这样, 在容器中, 可以任意修改sysctl内核参数.

> 注意: 如果修改的是namespaced的参数, 则不会影响host和其他容器. 反之, 则会影响它们. 

如果想在容器中修改主机的`net.ipv4.ip_default_ttl`参数, 则除了`--privileged`, 还需要加上 `--net=host`. 

### 1.2 把`/proc/sys`挂载到容器里面

```
docker run -v /proc/sys:/proc/sys -it ubuntu bash
```

然后写容器内的`/proc/sys`文件.

注意: 这样操作, 效果类似于`--privileged`, 对于`namespaced`的参数, 不会影响host和其他容器. 

### 1.3 `docker run`的`--sysctl`选项

使用`docker`提供的`--sysctl`指定内核参数.

```console
$ docker run -it --sysctl 'net.ipv4.ip_default_ttl=63' ubuntu sysctl net.ipv4.ip_default_ttl
net.ipv4.ip_default_ttl = 63
```

注意:

1. 只有`namespaced`参数才可以, 否则会报错”invalid argument…”
2. 这种方式只是在容器初始化过程中完成内核参数的修改, 容器运行起来以后, `/proc/sys`仍然是以只读方式挂载的, 在容器中不能再次修改`sysctl`内核参数. 

## 2. kuber 中设置 sysctl

### 2.1 通过`sysctls`和`unsafe-sysctls`注解

k8s进一步把`syctl`参数分为`safe`和`unsafe`, 非`namespaced`的参数, 肯定是unsafe. namespaced参数, 也只有一部分被认为是unsafe的. 

`safe`的条件：

1. must not have any influence on any other pod on the node(不能影响当前主机的其他 Pod)
2. must not allow to harm the node’s health(不可防碍宿主机的正常运行)
3. must not allow to gain CPU or memory resources outside of the resource limits of a pod(不可占用超过 resource limit 限制的 cpu/内存资源)

在`pkg/kubelet/sysctl/whitelist.go`中维护了 safe sysctl 参数的名单. 在1.7.8的代码中, 只有三个参数被认为是safe的：

- kernel.shm_rmid_forced,
- net.ipv4.ip_local_port_range,
- net.ipv4.tcp_syncookies

如果要设置一个Pod中的`safe`参数, 通过`security.alpha.kubernetes.io/sysctls`这个annotation来传递给kubelet. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: sysctl-example
  annotations:
    security.alpha.kubernetes.io/sysctls: kernel.shm_rmid_forced=1
spec:
  ## ...
```

如果要设置一个namespaced, 但是unsafe的参数, 要使用另一个`annotation: security.alpha.kubernetes.io/unsafe-sysctls`, 另外还要给kubelet一个特殊的启动参数. 

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: sysctl-example
  annotations:
    security.alpha.kubernetes.io/sysctls: kernel.shm_rmid_forced=1
    security.alpha.kubernetes.io/unsafe-sysctls: net.ipv4.route.min_pmtu=1000,kernel.msgmax=1 2 3
spec:
  ## ...
```

kubelet 增加`--experimental-allowed-unsafe-sysctls`启动参数

```
kubelet --experimental-allowed-unsafe-sysctls 'kernel.msg*,net.ipv4.route.min_pmtu'
```

> 没实验过, 不过看起来像是要事先把所有需要用到的`unsafe`都列举出来才行.

### 2.2 privileged Pod

如果要修改的是非namespaced的参数, 如`vm.*`, 那就没办法使用以上方法. 可以给Pod privileged权限, 然后在容器的初始化脚本或代码中去修改sysctl参数(反正是全局的). 

创建Pod/Deployment/Daemonset等对象时, 给容器的spec指定`securityContext.privileged=true`

```yaml
spec:
  containers:
  - image: nginx:alpine
    securityContext:
      privileged: true

```

这样跟`docker run –privileged`效果一样, 在Pod中`/proc`是以`rw`权限mount的, 可以直接修改相关sysctl内核参数. 
