网卡, 黑屏, 需要等待...

从外面官网拷贝东西拷不进去, 需要在桌面里用火狐访问官网(这样能打开多个标签页了吗)

ssh 到其他节点的密码? 要用 root 用户, 免密登录.

是否能刷新网页

支持多终端


https://killer.sh/attendee/72625245-e04b-4f37-808a-7d4aca267718/content

## 1

Task weight: 1%

You have access to multiple clusters from your main terminal through kubectl contexts. Write all those context names into /opt/course/1/contexts.

Next write a command to display the current context into /opt/course/1/context_default_kubectl.sh, the command should use kubectl.

Finally write a second command doing the same thing into /opt/course/1/context_default_no_kubectl.sh, but without the use of kubectl.
    - 不使用 kubectl 获取当前 context, 从 .kube/config 文件中获取?

kubectl config view

## 9

Task weight: 5%

Use context: kubectl config use-context k8s-c2-AC

Ssh into the controlplane node with ssh cluster2-controlplane1. Temporarily stop the kube-scheduler, this means in a way that you can start it again afterwards.

Create a single Pod named manual-schedule of image httpd:2.4-alpine, confirm it's created but not scheduled on any node.

Now you're the scheduler and have all its power, manually schedule that Pod on node cluster2-controlplane1. Make sure it's running.

Start the kube-scheduler again and confirm it's running correctly by creating a second Pod named manual-schedule2 of image httpd:2.4-alpine and check if it's running on cluster2-node1.

手动调度pod, 通过指定 nodeName, 可以越过调度器.

## 18

Task weight: 8%

Use context: kubectl config use-context k8s-c3-CCC

There seems to be an issue with the kubelet not running on cluster3-node1. Fix it and confirm that cluster has node cluster3-node1 available in Ready state afterwards. You should be able to schedule a Pod on cluster3-node1 afterwards.

Write the reason of the issue into /opt/course/18/reason.txt.

启动脚本处的 kubelect 路径不正确.

## 22(证书检测)

Task weight: 2%

Use context: kubectl config use-context k8s-c2-AC

Check how long the kube-apiserver server certificate is valid on cluster2-controlplane1. Do this with openssl or cfssl. Write the exipiration date into /opt/course/22/expiration.

Also run the correct kubeadm command to list the expiration dates and confirm both methods show the same date.

Write the correct kubeadm command that would renew the apiserver server certificate into /opt/course/22/kubeadm-renew-certs.sh.

## 23

Task weight: 2%

Use context: kubectl config use-context k8s-c2-AC

Node cluster2-node1 has been added to the cluster using kubeadm and TLS bootstrapping.

Find the "Issuer" and "Extended Key Usage" values of the cluster2-node1:

kubelet client certificate, the one used for outgoing connections to the kube-apiserver.
kubelet server certificate, the one used for incoming connections from the kube-apiserver.
Write the information into file /opt/course/23/certificate-info.txt.

Compare the "Issuer" and "Extended Key Usage" fields of both certificates and make sense of these.

## 25

Task weight: 8%

Use context: kubectl config use-context k8s-c3-CCC

Make a backup of etcd running on cluster3-controlplane1 and save it on the controlplane node at /tmp/etcd-backup.db.

Then create a Pod of your kind in the cluster.

Finally restore the backup, confirm the cluster is still working and that the created Pod is no longer with us.

备份卡住

而且无法从pod中拷贝出文件(crictl 没 cp 命令)

[openshift 4.5.9 etcd损坏+脑裂修复过程](https://zhangguanzhang.github.io/2021/06/08/ocp4.5.9-restore-etcd/)

## Extra Question 1(驱逐)

Use context: kubectl config use-context k8s-c1-H

Check all available Pods in the Namespace project-c13 and find the names of those that would probably be terminated first if the Nodes run out of resources (cpu or memory) to schedule all Pods. Write the Pod names into /opt/course/e1/pods-not-stable.txt.

## Preview Question 1
Use context: kubectl config use-context k8s-c2-AC

The cluster admin asked you to find out the following information about etcd running on cluster2-controlplane1:

Server private key location
Server certificate expiration date
Is client certificate authentication enabled
Write these information into /opt/course/p1/etcd-info.txt

Finally you're asked to save an etcd snapshot at /etc/etcd-snapshot.db on cluster2-controlplane1 and display its status.

etcd备份恢复

## Preview Question 2
Use context: kubectl config use-context k8s-c1-H

You're asked to confirm that kube-proxy is running correctly on all nodes. For this perform the following in Namespace project-hamster:

Create a new Pod named p2-pod with two containers, one of image nginx:1.21.3-alpine and one of image busybox:1.31. Make sure the busybox container keeps running for some time.

Create a new Service named p2-service which exposes that Pod internally in the cluster on port 3000->80.

Find the kube-proxy container on all nodes cluster1-controlplane1, cluster1-node1 and cluster1-node2 and make sure that it's using iptables. Use command crictl for this.

Write the iptables rules of all nodes belonging the created Service p2-service into file /opt/course/p2/iptables.txt.

Finally delete the Service and confirm that the iptables rules are gone from all nodes.

使用 crictl 确认 kube-proxy 使用了 iptables

