# 强制删除Terminating状态的pod

```
$ k get pod
NAME                                 READY   STATUS        RESTARTS   AGE
xd-spider-195-5764676576-4mqgc       1/1     Terminating   0          14h
xd-spider-195-5764676576-nllpn       1/1     Terminating   0          14h
```

```
k delete pod --force --grace-period=0 xd-spider-195-5764676576-4mqgc
warning: Immediate deletion does not wait for confirmation that the running resource has been terminated. The resource may continue to run on the cluster indefinitely.
pod "xd-spider-195-5764676576-4mqgc" force deleted
```

```
$ k get pod
NAME                                 READY   STATUS        RESTARTS   AGE
xd-spider-195-5764676576-nllpn       1/1     Terminating   0          14h
```
