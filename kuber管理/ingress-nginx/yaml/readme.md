当前目录是 0.25.0 版本的部署文件, 我把ta们拆分出来了, 官方的yaml配置只包括ns, configmap, rbac和deploy 4个部分.

但是毕竟是nginx, 对外总要用一个访问入口, 所以至少应该包含一个service吧? 不然外部请求无法转发进来.

`05.svc.yaml`是我自己添加的针对ingress controller的 service 部署文件. 这样, 在部署 ingress 资源的时候, 前端负载均衡器就可以将请求转发给 worker 节点上的 30080/30443 端口, 实现 service 对外提供服务的目的.

但是这样有一个问题, 就是在扩容 worker 节点的时候, 是不是也要把新增的 worker 添加到前端负载均衡器的配置中呢? 另外, deploy 的 replica 值怎么设置呢? 只在单个节点上运行压力应该会很大吧?

铜板街的做法是, 修改 deploy 为`hostNetwork: true`, 将 controller 这个pod部署在3台 master节点上, 然后前端SLB只要把请求请求转发到这3台就可以了. 我倒是觉得不失为一种优雅的解决方案, `04.deploy-tbj.yaml`就是这种方案的配置.
