参考文章

1. [FlexVolume](https://feisky.gitbooks.io/kubernetes/plugins/flex-volume.html)
2. [k8s与存储--flexvolume解读](https://segmentfault.com/a/1190000020320771)

1. `FlexVolume`与`CSI`是两个并列的概念, 作用可视作相同.
2. `CSI`可与`PV/PVC`配合使用, 而`FlexVolume`只能单独使用.
3. 在实际部署过`CSI`和`FlexVolume`后会发现, 以NFS为例
    - `CSI`的基本逻辑是, nfs 服务端创建一个大的共享目录, 之后通过 PVC/SC 的绑定, Pod 可以通过指定 PV 字段, 在那个大目录下创建各自的子目录.
    - `FlexVolume`, 每创建一个 Pod, 都需要事先为其在 nfs 服务端创建对应的目录, 一一对应, 无法实现自动化, 不够智能.

CSI出现后, 建议使用CSI替代.


volumes -> hostPath, 然后`volumeMounts`进行挂载.

volumes -> persistentVolumeClaim, pvc 绑定 pv, pv 指定 hostPath.

除了`hostPath`, `nfs`类型的卷, 我还只见过通过`volumes`直接挂载, 没见过通过 pv/pvc 挂载的.
