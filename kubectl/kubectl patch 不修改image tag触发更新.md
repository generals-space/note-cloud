# patch 不修改image tag触发更新

参考文章

1. [k8s使用技巧](https://blog.csdn.net/rariki/article/details/78595830)

在 kube 集群中, 修改了 configmap 变量后, 相关的 Pod 并不会自动重启, 需要我们手动完成.

一般这种场景就是 daemonset, deployment 等, 且 replicas 数值比较大, 一个一个删除重建太麻烦, 是否有一个办法能够触发重启机制呢?

参考文章1对 label 和 env 变量两种情况分别做了讨论, 这里我们也只使用 env 的方式, 通过 patch 子命令完成.

```bash
k patch ds -n kube-system calico-node -p '{"spec": {"template": {"spec": {"containers": [{"name": "calico-node", "env": [{"name":"updateTime","value":"123456"}]}]}}}}'
```


```bash
## $function:   不修改 image tag 触发 daemonset/deployment 的自动更新
## $1:          
function kup() {
    k patch ds -n kube-system calico-node -p '{"spec": {"template": {"spec": {"containers": [{"name": "calico-node", "env": [{"name":"updateTime","value":"123456"}]}]}}}}'
}
```
