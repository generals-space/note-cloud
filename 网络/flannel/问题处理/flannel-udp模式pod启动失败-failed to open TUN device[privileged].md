# flannel-udp模式pod启动失败-failed to open TUN device

参考文章

1. [Failed to create tun device: open /dev/net/tun: no such file or directory](https://blog.csdn.net/wen_dy/article/details/78856079)
    - 没废话, 不bb
2. [failed to open TUN device: open /dev/net/tun: no such file or directory](https://github.com/coreos/flannel/issues/1267)
    - `privileged: true`

将`02-cm.yaml`中的`Type`从`vxlan`修改成`udp`后创建Pod, 发现Pod一直处于`CrashLoopBackOff`状态. 

查看日志如下

```log
$ k logs -f kube-flannel-ds-amd64-mng7n
...
I0424 10:54:03.024305       1 main.go:386] Found network config - Backend type: udp
E0424 10:54:03.024484       1 main.go:289] Error registering network: failed to open TUN device: open /dev/net/tun: no such file or directory
I0424 10:54:03.024553       1 main.go:366] Stopping shutdownHandler...
```

执行如下命令创建`tun`设备即可.

```
mknod /dev/net/tun c 10 200
```

`c`表示`character`字符设备.

...我错了

```log
$ mknod /dev/net/tun c 10 200
mknod: "/dev/net/tun": 文件已存在
```

要把`03-ds.yaml`中daemonset的`securityContext`设置为`privileged: true`.
