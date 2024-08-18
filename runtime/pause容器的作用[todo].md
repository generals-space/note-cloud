# pause容器的作用

参考文章

1. [The Almighty Pause Container](https://www.ianlewis.org/en/almighty-pause-container)
    - nginx+ghost实例展示容器间可以共享network, pid, ipc等命名空间
    - 解释了僵尸进程产生的原因以及在使用docker共享ns, kuber多容器pod场景下产生僵尸进程的可能性, 由此引出了kube集群中pause容器的作用.
