# docker stats查看容器资源占用情况

参考文章

1. [聊聊 Docker 容器的资源管理](https://cloud.tencent.com/developer/article/1563718)
    - 介绍了`docker stats`和`docker top`子命令的使用方法与输出解释.
    - `sha256sum /dev/zero`消耗CPU资源的命令(单线程).
    - `docker update`更新指定容器的资源限制参数.
    - **`--cpus`是在 Docker 1.13 时新增的, 可用于替代原先的`--cpu-period`和`--cpu-quota`**
    - 一般情况下, 推荐直接使用`--cpus`, 而无需单独设置`--cpu-period`和`--cpu-quota`
    - 主要介绍CPU/Memory的相关参数与实验.

如果不加任何限制, docker 容器将共享宿主机的资源: CPU/内存, 没有限制.

```log
# 启动一个容器
$ docker run -d redis
c98c9831ee73e9b71719b404f5ecf3b408de0b69aec0f781e42d815575d28ada
# 查看其所占用资源的情况
$ docker stats --no-stream $(docker ps -ql)
CONTAINER ID        NAME                CPU %    MEM USAGE / LIMIT     MEM %    NET I/O        BLOCK I/O    PIDS
c98c9831ee73        amazing_torvalds    0.08%    2.613MiB / 15.56GiB   0.02%    3.66kB / 0B    0B / 0B      4
```

- `Container ID`:       容器的 ID, 也是一个容器生命周期内不会变更的信息.
- `Name`:               容器的名称, 如果没有手动使用 --name 参数指定, 则 Docker 会随机生成一个, 运行过程中也可以通过命令修改.
- `CPU %`:              容器正在使用的 CPU 资源的百分比, `100%`表示的是单核 cpu 满载. 这里面涉及了比较多细节, 下面会详细说.
- `Mem Usage/Limit`:    当前内存的使用及容器可用的最大内存(这里我使用了一台 16G 的电脑进行测试).
- `Mem %`:              容器正在使用的内存资源的百分比.
- `Net I/O`:            容器通过其网络接口**发送**和**接收**到的数据量(使用`ping`命令可以看到明显变化).
- `Block I/O`:          容器通过块设备**读取**和**写入**的数据量.
- `Pids`:               容器创建的进程或线程数.

> 注意: 进入容器内部然后使用`free`, `top`等命令查询容器资源总量是不正确的, ta们总是能获取宿主机的资源总量, 无法得到自己本身的信息, 因为ta们获取信息的来源是`/proc/{meminfo, $PID}`.

## 参数`--cpuset-cpus`与`--cpus`

`--cpuset-cpus`: 可以设置容器绑定的 CPU 核心序号列表, 指定多个可用逗号分隔. 比较容易理解.

`--cpus`: 可以是小数, 小数点后精度为2位, 值不可超过宿主机CPU核心数, 否则会出现如下错误(如下示例为4核宿主机).

```log
Error response from daemon: Range of CPUs is from 0.01 to 4.00, as there are only 4 CPUs available
```

一台8核服务器的 CPU 使用率最高可达 800%, 当`--cpus`设置为2时, 容器最多可占用200%的CPU使用率. 这是未设置`--cpuset-cpus`时, 假如设置了`--cpuset-cpus`为2, 那么CPU使用率的分配上限就变成了200%. 但这不意味着`--cpus`的值的取值范围就成了`0.01 - 2.00`, `--cpus`的取值范围只受限于宿主机的 CPU 核心数, 与`--cpuset-cpus`的设置值无关(不过容器实际可使用的上限其实就是 200%了).

另外, 一般情况下, 推荐直接使用`--cpus`, 而无需单独设置`--cpu-period`和`--cpu-quota`.

