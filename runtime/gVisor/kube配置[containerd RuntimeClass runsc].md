参考文章

1. [Kubernetes Quick Start](https://gvisor.dev/docs/user_guide/quick_start/kubernetes/)
    - 选择 containerd 那一节, 见参考文章2
2. [Containerd Quick Start](https://gvisor.dev/docs/user_guide/containerd/quick_start/)
    - containerd 下的配置
    - `containerd-shim-runsc-v1`的安装链接, 见参考文章3
3. [Installation](https://gvisor.dev/docs/user_guide/install/)
4. [Runtime Class](https://kubernetes.io/docs/concepts/containers/runtime-class/)

在containerd中配置gvisor, 还要求存在`containerd-shim-runsc-v1`, 我看了下, 直接用 yum 安装的 containerd 是没有这个包的.

```console
## 这里按tab补全
$ containerd
containerd               containerd-shim          containerd-shim-runc-v1  containerd-shim-runc-v2  containerd-stress
```

```bash
ARCH=$(uname -m)
URL=https://storage.googleapis.com/gvisor/releases/release/latest/${ARCH}
wget ${URL}/runsc ${URL}/containerd-shim-runsc-v1
chmod a+rx runsc containerd-shim-runsc-v1
sudo mv runsc containerd-shim-runsc-v1 /usr/local/bin
```

按照参考文章3中的步骤, 下载并安装了gvisor后, 就可以使用了.

### docker

可以使用 docker 直接调用 gvisor 作为 runtime, 需要先执行一下`runsc install`子命令行.

```console
$ /usr/local/bin/runsc install
2023/03/11 08:56:58 Runtime runsc not found: adding
2023/03/11 08:56:58 Successfully updated config.
$ systemctl reload docker
$ docker run -it --name test --rm --runtime=runsc centos:7 bash
```

启动后相关进程列表如下.

```log
root       80188   52521  0 08:58 pts/0    00:00:00 docker run -it --name test --rm --runtime=runsc centos:7 bash
root       80338       1  0 08:58 ?        00:00:00 /usr/bin/containerd-shim-runc-v2 -namespace moby -id 978aa1a28992090f6fb02793cd65f623697cfba5d18cf0af0456b5a60e6ee785 -address /run/containerd/containerd.sock
root       80357   80338  0 08:58 ?        00:00:00 runsc-gofer --root=/var/run/docker/runtime-runc/moby --log=/run/containerd/io.containerd.runtime.v2.task/moby/978aa1a28992090f6fb02793cd65f623697cfba5d18cf0af0456b5a60e6ee785/log.json --log-format=json --systemd-cgroup=true --log-fd=3 gofer --bundle /run/containerd/io.containerd.runtime.v2.task/moby/978aa1a28992090f6fb02793cd65f623697cfba5d18cf0af0456b5a60e6ee785 --spec-fd=4 --mounts-fd=5 --io-fds=6 --io-fds=7 --io-fds=8 --io-fds=9 --apply-caps=false --setup-root=false --sync-userns-fd=-1 --proc-mount-sync-fd=16
nobody     80358   80338  0 08:58 pts/2    00:00:00 runsc-sandbox --root=/var/run/docker/runtime-runc/moby --log=/run/containerd/io.containerd.runtime.v2.task/moby/978aa1a28992090f6fb02793cd65f623697cfba5d18cf0af0456b5a60e6ee785/log.json --log-format=json --systemd-cgroup=true --log-fd=3 boot --proc-mount-sync-fd=22 --product-name None --bundle=/run/containerd/io.containerd.runtime.v2.task/moby/978aa1a28992090f6fb02793cd65f623697cfba5d18cf0af0456b5a60e6ee785 --io-fds=4 --io-fds=5 --io-fds=6 --io-fds=7 --mounts-fd=8 --start-sync-fd=9 --controller-fd=10 --spec-fd=11 --cpu-num 1 --total-memory 2079682560 --stdio-fds=12 --stdio-fds=13 --stdio-fds=14 978aa1a28992090f6fb02793cd65f623697cfba5d18cf0af0456b5a60e6ee785
```

由于 docker 底层是 containerd, 之后由 containerd 再调用 gvisor. 

### containerd

不过containerd 自己通过 crictl 创建的容器还是用的 runc, 需要先修改一下配置文件`/etc/containerd/config.toml`.

```bash
cat <<EOF | sudo tee /etc/containerd/config.toml
version = 2
[plugins."io.containerd.runtime.v1.linux"]
  shim_debug = true
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
  runtime_type = "io.containerd.runc.v2"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runsc]
  runtime_type = "io.containerd.runsc.v1"
EOF
```

> 这一步类似于上面的`runsc install`.

一般来说, `containerd config default`打印出来的内容中, 已经包含了`containerd.runtimes.runc`部分, 我们只需要添加`runsc`部分即可.

修改完配置, 并重启containerd后, 执行如下命令使用.

## 在 kube 集群中使用

先创建 gvisor 的 runtime class.

```console
$ cat <<EOF | kubectl apply -f -
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: gvisor
handler: runsc
EOF

$ k get runtimeclass
NAME     HANDLER   AGE
gvisor   runsc     7s
```

然后创建pod, 指定使用 gvisor 作为 runtime(默认为 runc)

```
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: nginx-gvisor
spec:
  runtimeClassName: gvisor
  containers:
  - name: nginx
    image: nginx
EOF
```

> 注意, 安装gvisor与配置containerd要在所有节点上进行, 否则pod被调度到的主机上没有启用gvisor的话, 就无法正常启动.

```
$ kubectl describe pod nginx-gvisor
3m47s       Warning   FailedCreatePodSandBox   pod/nginx-gvisor    Failed to create pod sandbox: rpc error: code = Unknown desc = RuntimeHandler "runsc" not supported 
```
