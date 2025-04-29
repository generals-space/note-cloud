
参考文章

1. [kubernetes-issue-1：ephemeral-storage引发的pod驱逐问题](https://cloud.tencent.com/developer/article/1456389)
    - 在每个Kubernetes的节点上，kubelet的根目录(默认是/var/lib/kubelet)和日志目录(/var/log)保存在节点的主分区上，这个分区同时也会被Pod的EmptyDir类型的volume、容器日志、镜像的层、容器的可写层所占用。
    - 磁盘容量不够时，kubernetes会清理镜像文件去腾空磁盘，省出资源去优先启动pod。
2. [体验ephemeral-storage特性来对Kubernetes中的应用做存储的限制和隔离](https://developer.aliyun.com/article/594066)
    - 构造了一个因 ephemeral-storage 而被驱逐的 Pod 示例.

## 问题描述

kube: v1.25.1

集群运行得好好的, Pod突然异常了, 还不只一个, 其他命名空间中也有多个Pod变成了`ContainerStatusUnknown`状态.

```log
[root@k8s-master-01 ~]# kwd pod
NAME                                                              READY   STATUS                   RESTARTS        AGE   IP               NODE         
cluster-api-provider-aliyun-controller-manager-5dcfd7c8d4-57qmj   0/2     ContainerStatusUnknown   54 (2d6h ago)   41d   10.254.151.139   k8s-master-01
cluster-api-provider-aliyun-controller-manager-5dcfd7c8d4-d4976   0/2     ImagePullBackOff         0               81m   10.254.151.163   k8s-master-01
```

上面2个Pod属于同一个 deployment, 删掉`ContainerStatusUnknown`状态的Pod后会新建Pod, 但仍然无法正常运行.

describe一下有如下输出.

```log
[root@k8s-master-01 ~]# kde pod cluster-api-provider-aliyun-controller-manager-5dcfd7c8d4-57qmj
Status:           Failed
Reason:           Evicted
Message:          The node was low on resource: ephemeral-storage. Container manager was using 1475604Ki, which exceeds its request of 0. Container kube-rbac-proxy was using 40Ki, which exceeds its request of 0.
```

根据上述`.status.message`中的信息, 宿主机上的磁盘空间不足, 导致Pod被驱逐`Evicted`.

使用`df -h`

```log
[root@k8s-master-01 ~]# df -h
文件系统        容量  已用  可用 已用% 挂载点
devtmpfs        7.7G     0  7.7G    0% /dev
tmpfs           7.7G     0  7.7G    0% /dev/shm
tmpfs           7.7G  3.2M  7.7G    1% /run
/dev/vda1        59G   43G   14G   76% /
## 省略...
```

不过这话说的也不准确, describe node 发现主机是正常调度的, 没有任何异常事件.

```log
[root@k8s-master-01 ~]# kde node k8s-master-01
Taints:             <none>
Unschedulable:      false
Conditions:
  Type            Status  LastHeartbeatTime          LastTransitionTime         Reason                      Message
  ----            ------  -----------------          ------------------         ------                      -------
  MemoryPressure  False   Mon, 27 May 2024 16:23:37  Thu, 11 Apr 2024 16:09:31  KubeletHasSufficientMemory  kubelet has sufficient memory available
  DiskPressure    False   Mon, 27 May 2024 16:23:37  Mon, 27 May 2024 15:16:57  KubeletHasNoDiskPressure    kubelet has no disk pressure
  PIDPressure     False   Mon, 27 May 2024 16:23:37  Thu, 11 Apr 2024 16:09:31  KubeletHasSufficientPID     kubelet has sufficient PID available
  Ready           True    Mon, 27 May 2024 16:23:37  Thu, 11 Apr 2024 16:09:31  KubeletReady                kubelet is posting ready status
Allocated resources:
  (Total limits may be over 100 percent, i.e., overcommitted.)
  Resource           Requests     Limits
  --------           --------     ------
  cpu                1105m (27%)  500m (12%)
  memory             304Mi (1%)   468Mi (3%)
  ephemeral-storage  0 (0%)       0 (0%)
  hugepages-1Gi      0 (0%)       0 (0%)
  hugepages-2Mi      0 (0%)       0 (0%)
```

```log
$ kubectl proxy
Starting to serve on 127.0.0.1:8001
```

```log
[root@k8s-master-01 ~]# curl --silent -X GET http://127.0.0.1:8001/api/v1/nodes/k8s-master-01/proxy/configz | jq . | grep evictionHard -A 5
    "evictionHard": {
      "imagefs.available": "15%",
      "memory.available": "100Mi",
      "nodefs.available": "10%",
      "nodefs.inodesFree": "5%"
    },
```
