# FlexVolume基本使用-NFS

参考文章

1. [FlexVolume](https://feisky.gitbooks.io/kubernetes/plugins/flex-volume.html)
    - 在 v1.7 版本，部署新的 FlevVolume 插件后需要重启 kubelet 和 kube-controller-manager
    - v1.8 开始不需要重启它们了
    - 对于新的存储插件, 推荐基于 CSI 构建
2. [k8s与存储--flexvolume解读](https://segmentfault.com/a/1190000020320771)
3. [官方文档 kubernetes/examples](https://github.com/kubernetes/examples/tree/master/staging/volumes/flexvolume)
    - 写得真烂, 还不得不看...

`FlexVolume`插件与`CSI`插件一样, 需要一个类似于`provisioner`的程序, 这个程序需要开发者编写.

参考文章3是kubernetes提供的官方示例, 对于`FlexVolume`提供了`lvm`/`nfs`两种存储类型, 这里我们使用`nfs`为例.

首先在集群所有节点上安装`jq`命令, 该`provisioner`工具(即下面讲到的`nfs`)会用到.

```
yum install -y jq
```

然后`provisioner`的镜像需要自行构建, 参考文章3中`deploy`目录提供了基本步骤(呸).

我将构建相关的脚本拷贝了出来, 并做了一些修改以便能够正常使用, 在当前文件同级的`flexvolume-nfs`目录下. 执行如下命令进行构建.

```
cd flexvolume-nfs
docker build --no-cache=true -f dockerfile -t flex-nfs .
```

该目录中`deploy.sh`为`entrypoint`入口程序, 不过ta好像就是把`nfs`工具拷贝到节点相应的目录`/usr/libexec/kubernetes/kubelet-plugins/volume/exec/`下, 并命名为`k8s~nfs/nfs`, 全路径为`/usr/libexec/kubernetes/kubelet-plugins/volume/exec/k8s~nfs/nfs`.

其中`nfs`为实现了`FlexVolume`接口的执行程序(实际上就是一个`bash`脚本), 每创建一个Pod, `kubelet`就会使用`fork/exec`调用ta, 实现挂载/卸载的操作.

然后部署`flexvolume-nfs-ds.yaml`, 创建一个`DeamonSet`资源将这个镜像部署到所有节点, `nfs`工具也将拷贝到所有节点上.

现在可以创建测试用的 Pod 资源了, 见`flexvolume-nfs-pod.yaml`. Pod 启动后会自动挂载目标 nfs 服务器上指定的共享目录.

------

OK, 现在来看一看`FlexVolume`与`CSI`的区别, 同样以NFS共享存储为例.

`CSI`的基本逻辑是, nfs 服务端创建一个大的共享目录, 之后通过 PVC/SC 的绑定, Pod 可以通过指定 PV 字段, 在那个大目录下创建各自的子目录.

而上面我们创建的`FlexVolume`, 每创建一个 Pod, 都需要事先为其在 nfs 服务端创建对应的目录, 一一对应, 无法实现自动化, 不够智能.

简而言之, `CSI`不只实现了挂载, 还实现了创建的操作, 更有优势.

## FAQ

在做实验期间, 测试用的`Pod`总是无法正常启动, 一直处于`ContainerCreating`状态, `describe`一下, 有如下输出.

```
Events:
  Type     Reason       Age        From                    Message
  ----     ------       ----       ----                    -------
  Normal   Scheduled    <unknown>  default-scheduler       Successfully assigned default/nginx-nfs to k8s-worker-02
  Warning  FailedMount  15s        kubelet, k8s-worker-02  Unable to attach or mount volumes: unmounted volumes=[default-token-sccdd test], unattached volumes=[default-token-sccdd test]: failed to get Plugin from volumeSpec for volume "test" err=no volume plugin matched
```

实际上`provisioner`的`DaemonSet`Pod并没有报错, 而在测试Pod所在的节点上, 查看`/var/log/message`有报错.

```
May 28 13:59:36 k8s-worker-01 kubelet: E0528 13:59:36.729872     731 driver-call.go:267] Failed to unmarshal output for command: init, output: "", error: unexpected end of JSON input
May 28 13:59:36 k8s-worker-01 kubelet: W0528 13:59:36.729912     731 driver-call.go:150] FlexVolume: driver call failed: executable: /usr/libexec/kubernetes/kubelet-plugins/volume/exec/k8s~dummy/dummy, args: [init], error: fork/exec /usr/libexec/kubernetes/kubelet-plugins/volume/exec/k8s~dummy/dummy: no such file or directory, output: ""
May 28 13:59:36 k8s-worker-01 kubelet: E0528 13:59:36.729935     731 plugins.go:766] Error dynamically probing plugins: Error creating Flexvolume plugin from directory k8s~dummy, skipping. Error: unexpected end of JSON input
```

因为官方文档给出的示例是`dummy`存储, 所以会有`dummy`字样出现. 

可能是因为缓存的缘故, 将测试Pod和provisioner的`DaemonSet`全部删除, `/usr/libexec/kubernetes/kubelet-plugins/volume/exec`多余的插件也删除, 重新部署就好了.
