# calico网络策略

参考文章

1. [Kuber官方文档 网络策略](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
    - 网络策略通过网络插件来实现, 所以用户必须使用支持`NetworkPolicy`的网络解决方案(简单地创建资源对象, 而没有控制器来使它生效的话, 是没有任何作用的).
    - 有中文版, 但是关于对示例中策略的解释还是要看英文的...否则完全不明白在讲什么
    - 中文版页面的"网络策略入门指南"链接不存在, 英文版链接正常, 不过目标页面不是同一个了.
    - [Declare Network Policy](https://kubernetes.io/docs/tasks/administer-cluster/declare-network-policy/)
2. [借助 Calico，管窥 Kubernetes 网络策略](https://blog.fleeto.us/post/network-policy-basic-calico/)
    - 以Nginx为例简单介绍NetworkPolicy的基本使用方法
    - 该文文末给出的链接已经不存在了, 可以见参考文章1中的[Declare Network Policy]()链接
3. [calico 网络结合 k8s networkpolicy 实现租户隔离及部分租户下业务隔离](https://blog.csdn.net/qianggezhishen/article/details/80390598)
4. [Calico官方文档 Network policy](https://docs.projectcalico.org/v3.10/reference/resources/networkpolicy)

关于`NetworkPolicy`的声明还是很简单的(如果只像参考文章1中那样的话), 下面给出几个简单示例

**示例1**

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: default
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  ingress:
    - from:
        - podSelector: {}
```

只允许`default`命名空间下的各pod(应该还有svc资源等)相互访问(当然各pod还是能访问外网的, 因为没有限制`egress`), 但是宿主机无法访问这个pod, 其他空间的pod也无法访问.

虽然有`namespaceSelector`用于放行指定标签的ns, 但是貌似没有直接提供通过名称书写的规则. 

这里也只能以`NetworkPolicy`所在空间为"全局"设置, 通过设置为空的`podSelector`对同一命名空间下的pod不做限制来实现这样的功能, 但是无法对其他指定命名空间进行设置.

看来创建ns的时候还要指定好标签才行啊...

**示例2**

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: default
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  ingress:
    - from:
        - podSelector: {}
      ports:
        - protocol: TCP
          port: 80
```

与上一示例相比, 添加了`ports`块, 这里需要注意, ingress是一个`[]`, 但ta的成员是`{}`, `from`和`ports`是其中的两个键, 这者平级. 转换成json看可能更清楚一些.

```json
"ingress": [
    {
        "from": [
            {
                "podSelector": {}
            }
        ],
        "ports": [
            {
                "port": 80,
                "protocol": "TCP"
            }
        ]
    }
],
```

添加`ports`后, default空间下的pod就无法随意互相访问了, 只能互相访问彼此的`80`端口, 其他形式的请求都将被拒绝.
