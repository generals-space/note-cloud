ls /var/lib/containerd

io.containerd.content.v1.content: 本地镜像的所有layer层, 但是是压缩过的, 应该是从harbor仓库拉取到的原数据(镜像大小统计的不是这种层的信息).
io.containerd.metadata.v1.bolt
io.containerd.runtime.v1.linux
io.containerd.runtime.v2.task
io.containerd.snapshotter.v1.native: 镜像的每层layer解压后都放在这里, 每个layer层按 1, 2, 3等数字作为目录名, 每个layer层都包含本层发生变动的部分
io.containerd.snapshotter.v1.overlayfs
tmpmounts
