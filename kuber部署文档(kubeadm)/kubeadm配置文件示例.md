# kubeadm配置文件示例

参考文章

1. [kubeadm初始化配置文件范例](https://www.guojingyi.cn/912.html)
2. [官方文档 kubeadm config](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-config/)

```
kubeadm config print init-defaults
kubeadm config print init-defaults –component-configs  [KubeProxyConfiguration KubeletConfiguration]
kubeadm config print join-defaults
kubeadm config print join-defaults –component-configs  [KubeProxyConfiguration KubeletConfiguration]
```
