参考文章

1. [kube-proxy 源码解析](https://cizixs.com/2017/04/07/kube-proxy-source-code-analysis/)
2. [官方文档 kube-proxy命令行选项的中文翻译](https://kubernetes.io/zh/docs/reference/command-line-tools-reference/kube-proxy/)
3. [kubernetes网络相关总结](http://codemacro.com/2018/04/01/kube-network/)
    - 备份地址: [kubernetes网络相关总结](https://blog.csdn.net/dest_dest/article/details/80695734)
    - kuber 中 proxy 组件使用的 iptables 架构图, 值得一看.

proxy的ipvs模型在1.8版本时由华为提交pr. 

changelog可查[CHANGELOG-1.8](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.8.md#kube-proxy-ipvs-mode) 

pr地址为[Implement IPVS-based in-cluster service load balancing](https://github.com/kubernetes/kubernetes/pull/46580), pr中提到此pr解决了什么问题, 与iptables模型相比有什么优势.

但貌似直到1.17都没有实现GA.
