# kubelet无法启动connection refused

`kubeadm init`时输出如下结果

```
[kubelet-check] Initial timeout of 40s passed.

Unfortunately, an error has occurred:
	timed out waiting for the condition

This error is likely caused by:
	- The kubelet is not running
	- The kubelet is unhealthy due to a misconfiguration of the node in some way (required cgroups disabled)

If you are on a systemd-powered system, you can try to troubleshoot the error with the following commands:
	- 'systemctl status kubelet'
	- 'journalctl -xeu kubelet'

```

`ps`查看时kubelet的确没启动, 然后去`/var/log/message`查看有如下日志.

```
Jan 30 13:45:27 k8s-master-01 kubelet: E0130 13:45:27.704233   38608 reflector.go:153] k8s.io/kubernetes/pkg/kubelet/kubelet.go:458: Failed to list *v1.Node: Get https://k8s-server-lb:8443/api/v1/nodes?fieldSelector=metadata.name%3Dk8s-master-01&limit=500&resourceVersion=0: dial tcp 172.16.91.10:8443: connect: connection refused
```

费了半天劲, 最后发现是`nginx`负载均衡pod没启动, 我只有一个master节点, 但是在nginx配置文件upstream中写了3个后端节点.

```conf
    upstream kube-apiserver {
        least_conn;
        server k8s-master-01:6443 max_fails=1 fail_timeout=10s;
        server k8s-master-02:6443 max_fails=1 fail_timeout=10s;
        server k8s-master-03:6443 max_fails=1 fail_timeout=10s;
    }
```

于是nginx pod日志如下

```
2020/01/30 04:35:23 [emerg] 1#1: host not found in upstream "k8s-master-02:6443" in /etc/nginx/nginx.conf:13
nginx: [emerg] host not found in upstream "k8s-master-02:6443" in /etc/nginx/nginx.conf:13
```

因为这个问题, nginx pod没能启动, 8443端口也就没能监听, kubelet就连接不上apiserver. 把`k8s-master-02`和`k8s-master-03`删掉, 重启nginx pod就行了.

------

2020-03-21 更新

有次部署时, 将集群从1.17.2降到1.16.2, 把`kubelet`, `kubectl`和`kubeadm`都卸掉, 然后安装1.16.2版本的, 就出现了这样的错误, 是因为没有再次执行`systemctl enable kubelet`和`systemctl restart kubelet`, 卸载重装后需要重新执行这两条命令.

