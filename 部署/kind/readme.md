参考文章

[在 Kubernetes 中运行 Kubernetes](https://www.qikqiak.com/post/k8s-in-k8s/)

docker run -d --name k8s-master-02 --hostname k8s-master-02 --privileged=true -v /sys/fs/cgroup:/sys/fs/cgroup:ro registry.cn-hangzhou.aliyuncs.com/generals-space/centos7-systemd

d cp /etc/yum.repos.d/kubernetes.repo root-kube-node-01-1:/etc/yum.repos.d/
yum install -y kubeadm-1.17.2 kubelet-1.17.2 kubectl-1.17.2

kubeadm join kube-apiserver.generals.space:6443 --token abcdef.0123456789abcdef --discovery-token-ca-cert-hash sha256:b969766a8c9d9dfc3615dff8767c5bd6b8aa7930bdd699d6bb2213c434904c61 --ignore-preflight-errors=all

/usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --container-runtime=remote --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock

kubelet 对 /var/lib/kubelet 目录做了 remount, 因此容器内运行 kubelet 无法与宿主机共享此目录.

```log
[preflight] Running pre-flight checks
[preflight] The system verification failed. Printing the output from the verification:
KERNEL_VERSION: 3.10.0-1062.el7.x86_64
OS: Linux
CGROUPS_CPU: enabled
CGROUPS_CPUACCT: enabled
CGROUPS_CPUSET: enabled
CGROUPS_DEVICES: enabled
CGROUPS_FREEZER: enabled
CGROUPS_MEMORY: enabled
error execution phase preflight: [preflight] Some fatal errors occurred:
        [ERROR FileContent--proc-sys-net-bridge-bridge-nf-call-iptables]: /proc/sys/net/bridge/bridge-nf-call-iptables does not exist
        [ERROR SystemVerification]: failed to parse kernel config: unable to load kernel module: "configs", output: "", err: exit status 1
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
To see the stack trace of this error execute with --v=5 or higher
```


------

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
