参考文章


1. [切换容器的根文件系统](https://osh-2020.github.io/lab-4/pivot_root/)
    - `pivot_root new_root put_old`的使用方法, 包括`umount put_old`的方法.
2. [pivot_root](https://zhuanlan.zhihu.com/p/101096040)
    - `pivot_root . .`可以避免创建临时目录.
    - [rootfs: make pivot_root not use a temporary directory](https://github.com/opencontainers/runc/commit/f8e6b5af5e120ab7599885bd13a932d970ccc748)
