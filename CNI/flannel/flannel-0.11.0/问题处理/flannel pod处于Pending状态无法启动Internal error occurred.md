# flannel pod处于Pending状态无法启动Internal error occurred

kuber集群版本: 1.17.2
flannel: quay.io/coreos/flannel:v0.11.0-amd64

```console
$ k get pod -A
NAMESPACE     NAME                                    READY   STATUS    RESTARTS   AGE
kube-system   coredns-7f9c544f75-4f9lq                0/1     Pending   0          3h36m
kube-system   coredns-7f9c544f75-b6t4c                0/1     Pending   0          3h36m
kube-system   kube-flannel-ds-amd64-vn2q5             0/1     Pending   0          3m31s
$ k logs -f kube-flannel-ds-amd64-vn2q5 -n kube-system
Error from server (InternalError): Internal error occurred: Authorization error (user=kube-apiserver-kubelet-client, verb=get, resource=nodes, subresource=proxy)
```

明明之前在kuber集群为1.16.2的时候还好的, 1.17.2的时候貌似`kubectl logs`命令不好用了, 只能直接使用`docker logs`查对应容器的日志, 不过flannel容器根本没有创建...

后来把集群拆了, 重新setup一个1.16.2的版本, 就可以了...

