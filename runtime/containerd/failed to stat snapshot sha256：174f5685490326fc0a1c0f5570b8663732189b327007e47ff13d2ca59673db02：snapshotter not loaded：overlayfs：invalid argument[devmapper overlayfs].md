参考文章

1. [container的构建镜像失败：snapshotter not loaded: overlayfs: invalid argument](https://blog.csdn.net/jieshibendan/article/details/122574854)


```
# ctr images pull registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7
registry.cn-hangzhou.aliyuncs.com/generals-space/centos:7:                        resolved       |++++++++++++++++++++++++++++++++++++++|
layer-sha256:2d473b07cdd5f0912cd6f1a703352c82b512407db6b05b43f2553732b55df3bc:    done           |++++++++++++++++++++++++++++++++++++++|
elapsed: 15.4s                                                                    total:  212.4  (13.8 MiB/s)
unpacking linux/amd64 sha256:81b2d7d21f93b1d07e4da5cad7968aa5856aa026b02351bdf3bf61b013de712f...
ctr: failed to stat snapshot sha256:174f5685490326fc0a1c0f5570b8663732189b327007e47ff13d2ca59673db02: snapshotter not loaded: overlayfs: invalid argument
```

