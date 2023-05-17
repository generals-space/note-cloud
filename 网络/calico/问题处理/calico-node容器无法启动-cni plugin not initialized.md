参考文章

1. [Container runtime network not ready: cni config uninitialized [closed]](https://stackoverflow.com/questions/49112336/container-runtime-network-not-ready-cni-config-uninitialized)
    - systemctl stop apparmor
2. [network plugin is not ready: cni config uninitialized ](https://github.com/kubernetes/kubernetes/issues/48798)
    - KUBELET_NETWORK_ARGS

```
  Normal  Scheduled       11m                  default-scheduler  Successfully assigned kube-system/calico-node-cfd9b to dev-k8s-node6
  Normal  Started         11m                  kubelet            Started container calico-node
  Normal  SandboxChanged  9m36s                kubelet            Pod sandbox changed, it will be killed and re-created.
  Normal  Started         9m35s (x2 over 11m)  kubelet            Started container upgrade-ipam
  Normal  Pulled          9m35s (x2 over 11m)  kubelet            Container image "docker.io/calico/cni:v3.25.1" already present on machine
  Normal  Created         9m35s (x2 over 11m)  kubelet            Created container install-cni
  Normal  Pulled          9m35s (x2 over 11m)  kubelet            Container image "docker.io/calico/cni:v3.25.1" already present on machine
  Normal  Created         9m35s (x2 over 11m)  kubelet            Created container upgrade-ipam
  Normal  Started         9m34s (x2 over 11m)  kubelet            Started container install-cni
  Normal  Created         9m32s (x2 over 11m)  kubelet            Created container mount-bpffs
  Normal  Pulled          9m32s (x2 over 11m)  kubelet            Container image "docker.io/calico/node:v3.25.1" already present on machine
  Normal  Started         9m31s (x2 over 11m)  kubelet            Started container mount-bpffs
  Normal  Pulled          9m31s (x2 over 11m)  kubelet            Container image "docker.io/calico/node:v3.25.1" already present on machine
  Normal  Created         9m30s (x2 over 11m)  kubelet            Created container calico-node
  Normal  Killing         39s (x6 over 9m39s)  kubelet            Stopping container calico-node
```

```log
May 16 21:51:43 dev-k8s-node6 kubelet[536349]: E0516 21:51:43.430484  536349 kubelet.go:2349] "Container runtime network not ready" networkReady="NetworkReady=false reason:NetworkPluginNotReady message:Network plugin returns error: cni plugin not initialized"
May 16 21:51:48 dev-k8s-node6 kubelet[536349]: E0516 21:51:48.433749  536349 kubelet.go:2349] "Container runtime network not ready" networkReady="NetworkReady=false reason:NetworkPluginNotReady message:Network plugin returns error: cni plugin not initialized"
May 16 21:51:50 dev-k8s-node6 kubelet[536349]: I0516 21:51:50.352422  536349 scope.go:110] "RemoveContainer" containerID="e8a10d46c889026c15dc7647c52686164d2d53a6610b07324b9b91d76be7280d"
May 16 21:51:50 dev-k8s-node6 kubelet[536349]: E0516 21:51:50.353536  536349 pod_workers.go:951] "Error syncing pod, skipping" err="failed to \"StartContainer\" for \"calico-node\" with CrashLoopBackOff: \"back-off 20s restarting failed container=calico-node pod=calico-node-rzc8n_kube-system(1364aab9-1f25-48cc-8e07-aed35485566b)\"" pod="kube-system/calico-node-rzc8n" podUID=1364aab9-1f25-48cc-8e07-aed35485566b
May 16 21:51:53 dev-k8s-node6 kubelet[536349]: E0516 21:51:53.434967  536349 kubelet.go:2349] "Container runtime network not ready" networkReady="NetworkReady=false reason:NetworkPluginNotReady message:Network plugin returns error: cni plugin not initialized"
May 16 21:51:56 dev-k8s-node6 containerd[621]: time="2023-05-16T21:51:56.327171210+08:00" level=info msg="StopContainer for \"42aa16e4c9075093f7ecab00ce127531eae867ed7cb43c2463988df9bf36ec22\" with timeout 30 (s)"
May 16 21:51:56 dev-k8s-node6 containerd[621]: time="2023-05-16T21:51:56.328537350+08:00" level=info msg="Stop container \"42aa16e4c9075093f7ecab00ce127531eae867ed7cb43c2463988df9bf36ec22\" with signal terminated"
```

```
[root@dev-k8s-master1 ~]# calicoctl node status
Calico process is running.

IPv4 BGP status
+--------------+-------------------+-------+----------+--------------------------------+
| PEER ADDRESS |     PEER TYPE     | STATE |  SINCE   |              INFO              |
+--------------+-------------------+-------+----------+--------------------------------+
| 10.20.0.99   | node-to-node mesh | up    | 08:52:17 | Established                    |
| 10.20.0.101  | node-to-node mesh | up    | 08:52:09 | Established                    |
| 10.20.0.102  | node-to-node mesh | up    | 08:52:09 | Established                    |
| 10.20.0.103  | node-to-node mesh | up    | 08:52:20 | Established                    |
| 10.20.0.104  | node-to-node mesh | up    | 08:51:58 | Established                    |
| 10.20.0.105  | node-to-node mesh | up    | 08:52:19 | Established                    |
| 10.20.0.107  | node-to-node mesh | up    | 08:52:16 | Established                    |
| 10.20.0.108  | node-to-node mesh | up    | 08:52:09 | Established                    |
| 10.20.0.106  | node-to-node mesh | start | 13:58:55 | Passive Socket: Connection     |
|              |                   |       |          | closed                         |
+--------------+-------------------+-------+----------+--------------------------------+

IPv6 BGP status
No IPv6 peers found.
```


```
root@dev-k8s-node6:/var/log# systemctl status apparmor
● apparmor.service - Load AppArmor profiles
     Loaded: loaded (/lib/systemd/system/apparmor.service; enabled; vendor preset: enabled)
     Active: active (exited) since Mon 2023-05-15 11:50:47 CST; 1 day 10h ago
       Docs: man:apparmor(7)
             https://gitlab.com/apparmor/apparmor/wikis/home/
   Main PID: 535 (code=exited, status=0/SUCCESS)
      Tasks: 0 (limit: 19153)
     Memory: 0B
        CPU: 0
     CGroup: /system.slice/apparmor.service

5月 15 11:50:48 dev-k8s-node6 apparmor.systemd[535]: Restarting AppArmor
5月 15 11:50:48 dev-k8s-node6 apparmor.systemd[535]: Reloading AppArmor profiles
5月 15 11:50:47 dev-k8s-node6 systemd[1]: Starting Load AppArmor profiles...
5月 15 11:50:47 dev-k8s-node6 systemd[1]: Finished Load AppArmor profiles.
root@dev-k8s-node6:/var/log# systemctl stop apparmor
root@dev-k8s-node6:/var/log# systemctl status apparmor
● apparmor.service - Load AppArmor profiles
     Loaded: loaded (/lib/systemd/system/apparmor.service; enabled; vendor preset: enabled)
     Active: inactive (dead) since Tue 2023-05-16 22:06:49 CST; 2s ago
       Docs: man:apparmor(7)
             https://gitlab.com/apparmor/apparmor/wikis/home/
    Process: 542949 ExecStop=/bin/true (code=exited, status=0/SUCCESS)
   Main PID: 535 (code=exited, status=0/SUCCESS)
        CPU: 3ms

5月 15 11:50:48 dev-k8s-node6 apparmor.systemd[535]: Restarting AppArmor
5月 15 11:50:48 dev-k8s-node6 apparmor.systemd[535]: Reloading AppArmor profiles
5月 15 11:50:47 dev-k8s-node6 systemd[1]: Starting Load AppArmor profiles...
5月 15 11:50:47 dev-k8s-node6 systemd[1]: Finished Load AppArmor profiles.
5月 16 22:06:49 dev-k8s-node6 systemd[1]: Stopping Load AppArmor profiles...
5月 16 22:06:49 dev-k8s-node6 systemd[1]: apparmor.service: Succeeded.
5月 16 22:06:49 dev-k8s-node6 systemd[1]: Stopped Load AppArmor profiles.
root@dev-k8s-node6:/var/log# systemctl restart containerd
```


```
May 16 22:17:28 dev-k8s-node6 containerd[543200]: time="2023-05-16T22:17:28.076426672+08:00" level=error msg="failed to reload cni configuration after receiving fs change event(\"/etc/cni/net.d/calico-kubeconfig\": WRITE)" error="cni config load failed: failed to load CNI config list file /etc/cni/net.d/10-calico.conflist: error parsing configuration list: unexpected end of JSON input: invalid cni config: failed to load cni config"
```






```
ssh -y -NTf -R 10001:127.0.0.1:22 root@120.55.15.143

crictl --runtime-endpoint=unix:///run/containerd/containerd.sock ps -a

calico-node -confd

bird -R -s /var/run/calico/bird.ctl -d -c /etc/calico/confd/config/bird.cfg

calicoctl node status
```

```yaml
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
    type: RollingUpdate

        livenessProbe:
          exec:
            command:
            - /bin/calico-node
            - -felix-live
            - -bird-live
          failureThreshold: 6
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10
        readinessProbe:
          exec:
            command:
            - /bin/calico-node
            - -felix-ready
            - -bird-ready
          failureThreshold: 3
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 10

```
