# service-iptables vs ipvs

参考文章

1. [k8s集群 iptables vs ipvs性能测试](https://www.flftuu.com/2019/11/01/k8s%E9%9B%86%E7%BE%A4-iptables-vs-ipvs%E6%80%A7%E8%83%BD%E6%B5%8B%E8%AF%95/)

iptables 模式下对 node 节点的资源消耗明显大于 ipvs 模式

Iptables 存在的问题:

1. 规则顺序匹配延迟大
2. 访问 service 时需要遍历每条链知道匹配, 时间复杂度 O(N), 当规则数增加时, 匹配时间也增加
3. 规则更新延迟大
4. iptables 规则更新不是增量式的, 每更新一条规则, 都会把全部规则加载刷新一遍
5. 规则数大时, 会出现 kernel lock
6. svc 数增加到 5000 时, 会频繁出现"Another app is currently holding the xtables lock. Stopped waiting after 5s", 导致规则更新延迟变大, `kube-proxy`定期同步时也会因为超时导致 CrashLoopBackOff

iptables 如何实现负载均衡? 把单个 ServiceIP 转发到多个 PodIP ?
