# join添加节点

```
$ kubeadm token list
TOKEN                     TTL         EXPIRES                     USAGES                   DESCRIPTION   EXTRA GROUPS
1ugnpn.7qwbt58paanfov8a   <invalid>   2019-07-19T11:17:01+08:00   authentication,signing   <none>        system:bootstrappers:kubeadm:default-node-token
```

如下生成worker节点的join命令

```
$ kubeadm token create --print-join-command
kubeadm join k8s-master-7-13:8443 --token fw6ywo.1sfp61ddwlg1we27     --discovery-token-ca-cert-hash sha256:52cab6e89be9881e2e423149ecb00e610619ba0fd85f2eccc3137adffa77bb04
```

要添加master节点, 还要执行如下命令, 得到`certificate-key`

```
## --experimental-upload-certs已被废弃
## kubeadm init phase upload-certs --experimental-upload-certs
$ kubeadm init phase upload-certs --upload-certs
[upload-certs] Storing the certificates in Secret "kubeadm-certs" in the "kube-system" Namespace
[upload-certs] Using certificate key:
70f399e275cabef0bb2794ea76303da0220574f59e994e755378a359edb5a233
```

使用第一步生成的worker的join命令, 加上上面生成的`certificate-key`, 拼接起来组成master的join命令.

```
kubeadm join k8s-master-7-13:8443 --token fw6ywo.1sfp61ddwlg1we27     --discovery-token-ca-cert-hash sha256:52cab6e89be9881e2e423149ecb00e610619ba0fd85f2eccc3137adffa77bb04 --certificate-key 70f399e275cabef0bb2794ea76303da0220574f59e994e755378a359edb5a233
```