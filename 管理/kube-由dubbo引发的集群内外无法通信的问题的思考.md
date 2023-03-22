# kuber-由dubbo引发的集群内外无法通信的问题的思考

参考文章

1. [个人开发环境下和k8s集群内svc,pod网络互通](https://blog.csdn.net/sltin/article/details/100044930)
2. [generals-space/cni-terway](https://github.com/generals-space/cni-terway)
3. [k8s中hostname, hosts文件, DNS和代理问题, service和pod的访问问题](https://blog.csdn.net/luanpeng825485697/article/details/84108166)
4. [研究 Dubbo 网卡地址注册时的一点思考](https://dubbo.apache.org/zh-cn/blog/dubbo-network-interfaces.html)
    - 使用`-DDUBBO_IP_TO_REGISTRY`参数可以让服务注册到dubbo时为指定的地址.

国内很多公司都使用Java+Dubbo+ZK作为服务注册与发现的框架, 在将业务迁移到kuber集群的过程中发现, 由于集群内Pod处于独立的子网, 部署在Pod中的服务注册到dubbo中的地址是Pod的IP, 而总有些服务不想部署到集群中, 且部署在kuber节点之外. 这样双方就无法相互通信了.

不只是集群内外, 也有多个集群之间使用同一套dubbo+zk的通信问题, 本文就来探讨一下这种情况的解决方法.

------

其实kuber本身就遵循了微服务架构, 同时提供了微服务化的运行环境. 其中的Service资源就可以实现服务注册和发现的功能. 所以最标准的做法就是, 放弃dubbo, 全面发挥Service的作用...

不过大部分人是不会放弃dubbo的, 所以文章还得继续.

本质上kuber集群中的Pod运行在一个被隔离的内部网络, 所以可以使用vpn来解决这样的问题. 这个思路由参考文章1提出并进行了实践, 只要新增一个node专门部署openvpn服务即可. 我觉得可行, 还没实验过, 先mark一下. <???>

另外就是彻底一点, 使用类似于虚拟机的桥接网络的模型(与docker的桥接网络不同), 使集群中的Pod直接获得宿主机网络的IP地址, 这样可以直接实现Pod与宿主机网络集群外节点的通信, 此时Pod与宿主机处于平级. 

参考文章2中的项目就是实现了桥接网络的CNI插件. 不过这个插件并不完美, 因为这是我在本地编写并运行的. 如果部署在真实的网络中, 存在三层交换机的话, 可以配合dhcp+vlan为集群划分独立的小子网. 比如宿主机网络为`192.168.0.0/16`, 可在三层交换机上配置按照vlan id为集群分配`192.168.171.0/24`的地址.

不过这个插件也只是实现了这样的基础功能, 与flannel一样没什么特色, 远不如calico灵活和强大, 如果你同时还想部署calico, 那么文章还是要继续.

下面的方法要借用`nodePort`这种服务类型.

我们可以考虑让服务在注册的时候将服务地址指定到当前Pod所在的`nodeIP+nodePort`, 这样集群外的服务就可以通过访问`集群中的节点IP:服务的nodePort`来访问到对应的服务.

这样做也有弊端, 首先就是由于集群中的`nodePort`是共用的, 从集群的任意节点都可以访问到. 那么不同的服务就不能使用同一个`nodePort`, 需要预先划分, 同时不同业务线理论上也应划分不同的端口范围以便管理, 这会大大增加运维成本.

当你能够接受这些弊端, 那接下来就是具体实现了.

我们可以在部署文件中通过环境变量传入, 告诉Pod中的服务要注册的地址. 但是`nodePort`还好说, 这是一个固定值, 写在`Service`配置中, 只要保持两者一致即可. 但是Pod处于哪个节点却是不确定的, 那么如何告诉服务Pod当前运行在哪个节点呢?

```yaml
      containers:
      - name: xxx
        env:
        - name: DUBBO_HOST
          valueFrom:
            fieldRef:
              fieldPath: status.hostIP
        - name: DUBBO_PROTOCOL_PORT
          ## 这里的值要和该pod对应的Service资源中的nodePort字段保持一致.
          value: "30010"
```

使用`valueFrom`, 可以动态确定Pod的运行状态, 并将其写入到Pod的环境变量中. 至于有效的`fieldPath`有哪些, 可以查看`k get pod pod名称 -o yaml`, `metadata`, `spec`, `status`下面的字段都可以使用. 

据我所知, Eureka或Consul的java库都是支持读取环境变量的, dubbo也可以. 参考文章4中提到可以使用`-DDUBBO_IP_TO_REGISTRY=$DUBBO_HOST`指定该服务注册到宿主机IP上.

![](https://gitee.com/generals-space/gitimg/raw/master/e410797ef876bf89eaa89f0f89514e48.png)
