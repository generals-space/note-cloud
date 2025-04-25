[failed to pull image k8s.gcr.io/kube-apiserver:v1.22.4 failed to convert whiteout file \"usr/local/.wh..wh..opq\": operation not supported: unknown" #2608](https://github.com/kubernetes/kubeadm/issues/2608)
[[NFS] failed to pull image; failed to convert; operation not supported: unknown #6273](https://github.com/containerd/containerd/issues/6273)
[kinder: the latest etcd image cannot be pulled in the containerd base image #2223](https://github.com/kubernetes/kubeadm/issues/2223)
[[Failing Test] kubeadm-kinder-master (ci-kubernetes-e2e-kubeadm-kinder-master) #93223](https://github.com/kubernetes/kubernetes/issues/93223)
[kinder: prepull the additional images (pause, etcd..) on the host #2228](https://github.com/kubernetes/kubeadm/pull/2228)

kubeadm config images pull 或是手动用 ctr 下载 kube-proxy 镜像时, 会报如下错误.

```log
[root@k8s-master-02 ~]# ctr -n k8s.io images pull registry.cn-hangzhou.aliyuncs.com/google_containers/kube-proxy:v1.17.2
registry.cn-hangzhou.aliyuncs.com/google_containers/kube-proxy:v1.17.2:           resolved       |++++++++++++++++++++++++++++++++++++++| 
index-sha256:9a2940bd7718e1b7060b96bb33b0af44ebb4e3de0a0f80b1e6f1622118e4f430:    done           |++++++++++++++++++++++++++++++++++++++| 

elapsed: 0.4 s                                                                    total:   0.0 B (0.0 B/s)                                         
unpacking linux/amd64 sha256:9a2940bd7718e1b7060b96bb33b0af44ebb4e3de0a0f80b1e6f1622118e4f430...
INFO[0000] apply failure, attempting cleanup             error="failed to extract layer sha256:682fbb19de80799fed8b83bd8172050774c83294f952bdd8013d9cce2ab2f2a6: failed to convert whiteout file \"usr/lib/x86_64-linux-gnu/xtables/.wh..wh..opq\": operation not supported: unknown" key="extract-763670182-suGC sha256:561eaf129957ad575ddb973e7455175c1b74ed4bbaa925ff79080483c621087e"
ctr: failed to extract layer sha256:682fbb19de80799fed8b83bd8172050774c83294f952bdd8013d9cce2ab2f2a6: failed to convert whiteout file "usr/lib/x86_64-linux-gnu/xtables/.wh..wh..opq": operation not supported: unknown
```

ctr -n k8s.io images import --no-unpack ./kube-proxy.tar

## 修改 containerd snapshotter 为 native

[failed to create ship task: failed to mount rootfs component: invalid argument: unknown](https://github.com/containerd/containerd/issues/9260)
[failed to create shim task: failed to mount rootfs component: invalid argument: unknown](https://github.com/containerd/containerd/discussions/9667)
[Failed to mount rootfs component with overlay filesystem](https://github.com/k3s-io/k3s/issues/2755)
[Error: "failed to create ship task: failed to mount rootfs component: invalid argument: unknown" when executing keadm join.](https://github.com/kubeedge/kubeedge/issues/5088)

```log
$ dmesg -T
[Fri Nov 29 15:49:23 2024] overlayfs: filesystem on '/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/135/fs' not supported as upperdir
```

./nerdctl -n=k8s.io --snapshotter=native ps -a
./nerdctl --snapshotter=native ps -a
./nerdctl --snapshotter=native rm 
./nerdctl --snapshotter=native run -d --name testpause --network host registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6
