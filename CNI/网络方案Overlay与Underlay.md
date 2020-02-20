# CNI

参考文章

1. [深入浅出容器网络（一）](https://blog.firemiles.top/2018/05/27/深入浅出容器网络（一）/)
    - CNI和CNM的相关API, 没什么实用价值
2. [深入浅出容器网络（二）](https://blog.firemiles.top/2018/08/06/深入浅出容器网络（二）/)
    - 组网类型: `Underlay L2`, `Overlay L2`和`Overlay L3`
    - `Underlay L2`: `MACVLAN`, `IPVLAN`
    - `Underlay L3`: `flannel`的`host-gw`, `calico`的`BGP`
    - `Overlay L2`: `VxLAN`, `NVGRE`, `GENEVE`
    - `Overlay L3`: 
3. [Kubernetes网络方案的三大类别和六个场景](https://sq.163yun.com/blog/article/223878660638527488)
    - 总纲级别, 值得一读
    - 协议栈层级: L2, L3, L2+L3
    - 穿越形态: Underlay, Overlay
    - 隔离方式: Flat, VLAN, VxLAN/GRE
    - 分析了多个第三方网络插件提供的网络方案的分类, 优缺点及适用场景.

- `Overlay`: `Overlay`在云化场景比较常见. `Overlay`下面是受控的`VPC`网络, 当出现不属于`VPC`管辖范围中的`IP`或者`MAC`, `VPC`将不允许此`IP`/`MAC`穿越. 出现这种情况时, 我们都利用`Overlay`方式来做. 
- `Underlay`: 在一个较好的一个可控的网络场景下, 我们一般利用`Underlay`. 可以这样通俗的理解, 无论下面是裸机还是虚拟机, 只要网络可控, 整个容器的网络便可直接穿过去, 这就是`Underlay`. 

个人的理解是:

云上VPC可通过创建交换机组建子网, 但是不允许用户启用虚拟网卡, 虚拟交换机`bridge`或`macvlan`等这些虚拟网络设备, 也不允许用户私自为网卡添加额外的子网IP. 来自这些设备的流量将会被屏蔽, 可以说云上的网络环境是受限的物理网络, 这也是为了租户的隔离与安全着想. 所以`Underlay`的手段行不通, 只能依靠`Overlay`的形式.

而传统的IDC网络, 或是本地虚拟机, 却没有这样的限制. 

`Underlay`不等于L2, `Overlay`也不等于L3.

`MACVLAN`与`IPVLAN`就是`Underlay`(我觉得`VLAN`也算是, 因为VPC网络不允许带有`vlan tag`的流量), 只在IDC传统网络才可以使用, 且是L2.

以`flannel`的`VxLAN`模型为例, `flannel`会创建`vxlan`设备, 但是流经`vxlan`设备的数据包会被其封装成`UDP`包, 最终通过物理网卡发出. 由于不对物理网络做出修改, 因此可以称为`Overlay`. 而`VxLAN`本身就工作在L2型的`Overlay`, 再加上L3的`iptables`和L2的`ARP`, 可以说`flannel`的`VxLAN`模型是`Overlay + L2 + L3`.
