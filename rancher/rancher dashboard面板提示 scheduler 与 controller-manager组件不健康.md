# rancher dashboard面板提示 scheduler 与 controller-manager组件不健康

参考文章

1. [Rancher管理K8S组件controller-manager 和scheduler状态 Unhealthy](https://my.oschina.net/u/1431757/blog/4550843)
    - 有效
2. [Kubectl get componentstatus fails for scheduler and controller-manager](https://forums.rancher.com/t/kubectl-get-componentstatus-fails-for-scheduler-and-controller-manager/15801)
3. [Rancher 2.0 - Troubleshooting and fixing “Controller Manager Unhealthy Issue”](https://stackoverflow.com/questions/54827814/rancher-2-0-troubleshooting-and-fixing-controller-manager-unhealthy-issue)
4. [Unhealthy controller manager and scheduler after leaving it running overnight](https://github.com/rancher/rancher/issues/14036)

- kuber: 1.19.0
- rancher: v2.4.24, v2.5.7

kuber是通过`kubeadm`创建的集群, 然后通过rancher的"添加已有集群"进行纳管. 通过`kubectl`查看的集群状态为正常的, 但是 rancher 的 dashboard 上查看时发现`scheduler`和`controller-manager`两个组件都是`Unhealthy`.

按照参考文章1所说的方式进行排查, 实测有效.

查看`contrller-manager`的可用参数

```
$ k exec -it kube-controller-manager-k8s-master-1 -n kube-system -- kube-controller-manager --help
      --port int
                The port on which to serve unsecured, unauthenticated access. Set to 0 to disable. (default 10252) (DEPRECATED: see --secure-port instead.)
```

```
k get cs -n kube-system

```
