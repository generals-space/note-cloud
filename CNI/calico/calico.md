进入到calico-node容器, 查看其中启动的进程.

```console
$ ps -ef
UID         PID   PPID  C STIME TTY          TIME CMD
root          1      0  0 04:20 ?        00:00:00 /usr/bin/runsvdir -P /etc/service/enabled
root         47      1  0 04:20 ?        00:00:00 runsv felix
root         48      1  0 04:20 ?        00:00:00 runsv bird
root         49      1  0 04:20 ?        00:00:00 runsv bird6
root         50      1  0 04:20 ?        00:00:00 runsv confd
root         51     47 10 04:20 ?        00:05:25 calico-node -felix
root         52     50  0 04:20 ?        00:00:02 calico-node -confd
root        136     48  0 04:20 ?        00:00:02 bird -R -s /var/run/calico/bird.ctl -d -c /etc/calico/confd/config/bird.cfg
root        137     49  0 04:20 ?        00:00:02 bird6 -R -s /var/run/calico/bird6.ctl -d -c /etc/calico/confd/config/bird6.cfg
```

