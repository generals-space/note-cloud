参考文章

1. [gvisor的网络栈实现-简介](https://zboya.github.io/post/google_golang_netstack/)
2. [Kubernetes 安全容器技术 kata gvisor](https://www.finclip.com/news/f/38247.html)

gvisor 只是工程的名字, 不过安装包/执行文件应该叫`runsc`, 是与 runc 同级的, 所以启用ta需要在 contaienrd/cri-o 中完成.
