# 修改imagePullPolicy强制替换Pod镜像[IfNotPresent]

## 场景描述

有一个statefulset, 其中container的image字段为`centos7:v1.0`.

如果我们对`v1.0`版本的镜像, 做了一些修改, 但是又不希望为其额外创建一个版本(如`v1.1`), 而是希望重建构建一个`v1.0`, 然后强制推送到镜像仓库, 将原本的覆盖.

如何使新的`v1.0`镜像生效?

说明: 出于某种原因, 我们将Pod与Node主机做了绑定, Pod删除重建后仍然会被调度到原有节点.

## 1. imagePullPolicy: Always

将sts中`imagePullPolicy`镜像拉取策略修改为`Always`, 此时重启Pod, 不管新调度到的主机上是否有原本的`v1.0`镜像, kuber都会重新拉取最新的`v1.0`镜像.

## 2. imagePullPolicy: IfNotPresent

如果不更改`imagePullPolicy`字段, 保持默认`IfNotPresent`, 然后**手动**将新的`v1.0`拉取到主机上(此时, 旧的`v1.0`镜像tag就会变成`none`了), 重启Pod, 镜像会替换吗?

答案是: 可以.

`IfNotPresent`只是保证当本地镜像Tag存在是, 不会有"Pull"这一动作, 但是如果新镜像的Tag把旧镜像替换了, 如果Pod重启, 还是会使用新镜像启动的.

其实换一种思路可以更好地理解.

sts, deploy, ds中对于`container.image`字段, 单纯指定了镜像Tag, 并没有保存镜像id. 当我们删除一个Pod时, 被sts, deploy, ds各自的controller manager监听到Pod被删除, 则会重新构建一个Pod对象并提交, Kubelet会使用目标镜像重建一个. 

此时Pod controller才会通知kubelet重建容器, kubelet并不会管上一个Pod用的是哪个镜像, 而是直接按照镜像Tag来, 所以重建的容器使用的是用最新Tag的镜像.

------

如果先尝试在主机上使用`docker rmi centos:v1.0`移除旧的镜像, 可能会出现

```log
Error response from daemon: conflict: unable to remove repository reference "centos7:v1.0" (must force) - container 5c8987f39338 is using its referenced image 902dead0a02e
```

但是如果使用`docker rmi -f centos:v1.0`时, 只会将此镜像tag与image id解绑, 镜像本身还是存在的.

如果使用`docker rmi -f 902dead0a02e`强删镜像id时, 则会出现如下错误.

```log
Error response from daemon: conflict: unable to delete 902dead0a02e (cannot be forced) - image is being used by running container bfbf5ba1e22d
```

## 3. Pod中包含多个容器的场景

还有一种场景, 假设statefulset的Pod中, 包含2个container容器A和B. 

当我们只希望强制更新其中B的镜像, 但不能影响A的运行, 如何实现?

第1种方式由于要重启Pod, 所以确定不可行.

解决方法是, 到pod所在主机上, 将容器B需要更新的镜像拉下来(此时旧镜像的tag将变为`none`), 然后用`docker kill`将容器B干掉, kubelet会自动将容器B重新启动, 不影响容器A的运行, 并且, 新启动的容器B将会使用新版本的镜像.

> 容器B的镜像拉取策略不必修改为`Always`.

上述方法实践可行, 但是在应用之前仍需谨慎.
