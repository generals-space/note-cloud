# devicemapper目录清理

参考文章

1. [docker devicemapper 驱动空间未释放解决方案](https://zhuanlan.zhihu.com/p/566092834)
2. [Delete data in a container devicemapper can not free used space](https://github.com/moby/moby/issues/18867)
3. [解决 Docker 的 DeviceMapper 占用空间过大](https://www.cnblogs.com/gaoyuechen/p/17387133.html)
    - `--storage-opt dm.loopdatasize=8G` 设置 DeviceMapper 的 data 为 8G
    - `--storage-opt dm.loopmetadatasize=4G` 设置 metadata 为 4G
    - `--storage-opt dm.basesize=8G` 设置 单个镜像的大小不能大于 8G
4. [Docker 存储空间设置](https://blog.csdn.net/a12345676abc/article/details/100557202)
    - 采用dockerfile构建镜像时，出现the device has no space to left. 提示设备空间不足，或者 docker commit 提交容器保存镜像时，提示空间不足，往往时由于生成的目标镜像的尺寸大于docker默认配置的值。

docker 采用 devicemapper 存储驱动，容器内部读写、删除文件后，docker used space 未被释放，导致需要频繁删除、重建容器。

```bash
docker ps -qa | xargs docker inspect --format='{{ .State.Pid }}' | xargs -IZ fstrim /proc/Z/root/
```

还未测试, 貌似会清理已存在容器的硬盘空间.
