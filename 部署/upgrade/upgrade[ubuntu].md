参考文章

1. [kubeadm upgrade](https://kubernetes.io/zh-cn/docs/reference/setup-tools/kubeadm/kubeadm-upgrade/#cmd-upgrade-apply)


```
controlplane $ k get node -owide
NAME           STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
controlplane   Ready    control-plane   16d   v1.26.1   172.30.1.2    <none>        Ubuntu 20.04.5 LTS   5.4.0-131-generic   containerd://1.6.12
node01         Ready    <none>          16d   v1.26.1   172.30.2.2    <none>        Ubuntu 20.04.5 LTS   5.4.0-131-generic   containerd://1.6.12
```

```
controlplane $ kubeadm upgrade plan
[upgrade/config] Making sure the configuration is correct:
[upgrade/config] Reading configuration from the cluster...
[upgrade/config] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -o yaml'
[preflight] Running pre-flight checks.
[upgrade] Running cluster health checks
[upgrade] Fetching available versions to upgrade to
[upgrade/versions] Cluster version: v1.26.1
[upgrade/versions] kubeadm version: v1.26.1
[upgrade/versions] Target version: v1.26.2
[upgrade/versions] Latest version in the v1.26 series: v1.26.2

Components that must be upgraded manually after you have upgraded the control plane with 'kubeadm upgrade apply':
COMPONENT   CURRENT       TARGET
kubelet     2 x v1.26.1   v1.26.2

Upgrade to the latest version in the v1.26 series:

COMPONENT                 CURRENT   TARGET
kube-apiserver            v1.26.1   v1.26.2
kube-controller-manager   v1.26.1   v1.26.2
kube-scheduler            v1.26.1   v1.26.2
kube-proxy                v1.26.1   v1.26.2
CoreDNS                   v1.9.3    v1.9.3
etcd                      3.5.6-0   3.5.6-0

You can now apply the upgrade by executing the following command:

        kubeadm upgrade apply v1.26.2

Note: Before you can perform this upgrade, you have to update kubeadm to v1.26.2.
```

```
$ apt install kubeadm=1.26.0-00
controlplane $ kubeadm version
kubeadm version: &version.Info{Major:"1", Minor:"26", GitVersion:"v1.26.0", GitCommit:"b46a3f887ca979b1a5d14fd39cb1af43e7e5d12d", GitTreeState:"clean", BuildDate:"2022-12-08T19:57:06Z", GoVersion:"go1.19.4", Compiler:"gc", Platform:"linux/amd64"}
```


$ kubeadm upgrade apply 1.26.0
[upgrade/successful] SUCCESS! Your cluster was upgraded to "v1.26.0". Enjoy!

[upgrade/kubelet] Now that your control plane is upgraded, please proceed with upgrading your kubelets if you haven't already done so.


controlplane $ k get pod -A
NAMESPACE            NAME                                       READY   STATUS    RESTARTS        AGE
kube-system          calico-kube-controllers-5f94594857-9rx4h   1/1     Running   3 (2m19s ago)   16d
kube-system          canal-zdvkp                                2/2     Running   0               69m
kube-system          canal-zjj65                                2/2     Running   0               69m
kube-system          coredns-787d4945fb-pxx4v                   1/1     Running   0               32s
kube-system          coredns-787d4945fb-qwdd4                   1/1     Running   0               32s
kube-system          etcd-controlplane                          1/1     Running   0               2m13s
kube-system          kube-apiserver-controlplane                1/1     Running   0               83s
kube-system          kube-controller-manager-controlplane       1/1     Running   0               62s
kube-system          kube-proxy-jdqgh                           1/1     Running   0               31s
kube-system          kube-proxy-lm47w                           1/1     Running   0               26s
kube-system          kube-scheduler-controlplane                1/1     Running   0               46s
local-path-storage   local-path-provisioner-8bc8875b-7f4nb      1/1     Running   0               16d


controlplane $ apt install kubelet=1.26.0
Reading package lists... Done
Building dependency tree       
Reading state information... Done
E: Version '1.26.0' for 'kubelet' was not found
controlplane $ apt install kubelet=1.26.0-00
Reading package lists... Done
Building dependency tree       
Reading state information... Done
The following packages will be DOWNGRADED:
  kubelet
0 upgraded, 0 newly installed, 1 downgraded, 0 to remove and 100 not upgraded.
Need to get 20.5 MB of archives.
After this operation, 4096 B disk space will be freed.
Do you want to continue? [Y/n] y
Get:1 https://packages.cloud.google.com/apt kubernetes-xenial/main amd64 kubelet amd64 1.26.0-00 [20.5 MB]
Fetched 20.5 MB in 4s (5494 kB/s)  
dpkg: warning: downgrading kubelet from 1.26.1-00 to 1.26.0-00
(Reading database ... 72923 files and directories currently installed.)
Preparing to unpack .../kubelet_1.26.0-00_amd64.deb ...
Unpacking kubelet (1.26.0-00) over (1.26.1-00) ...
Setting up kubelet (1.26.0-00) ...
controlplane $ systemctl restart kubelet
controlplane $ k get node -owide
NAME           STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
controlplane   Ready    control-plane   16d   v1.26.0   172.30.1.2    <none>        Ubuntu 20.04.5 LTS   5.4.0-131-generic   containerd://1.6.12
node01         Ready    <none>          16d   v1.26.1   172.30.2.2    <none>        Ubuntu 20.04.5 LTS   5.4.0-131-generic   containerd://1.6.12

