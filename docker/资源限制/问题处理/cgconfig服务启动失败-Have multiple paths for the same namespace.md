# cgconfig服务启动失败-Have multiple paths for the same namespace

`/etc/cgconfig.d/gopls.conf`

```conf
group gopls {
    cpu {
        cpu.cfs_period_us = 100000;
        cpu.cfs_quota_us = 400000;
    }
    cpuset {
        cpuset.cpus = 0-3;
        cpuset.mems = 0;
    }
}
```

`/etc/cgrules.conf`

```
root:gopls cpu,cpuset gopls
```

但是启动失败.

```
$ systemctl restart cgconfig
Job for cgconfig.service failed because the control process exited with error code. See "systemctl status cgconfig.service" and "journalctl -xe" for details.
```

```
$ journalctl -xe

1月 20 14:40:06 general-work cgconfigparser[58050]: error at line number 18 at b0VIM:syntax error
1月 20 14:40:06 general-work cgconfigparser[58050]: /usr/sbin/cgconfigparser; error loading /etc/cgconfig.d/.gopls.conf.swp: Have multiple paths for the same namespace
1月 20 14:40:06 general-work systemd[1]: cgconfig.service: main process exited, code=exited, status=104/n/a
1月 20 14:40:06 general-work cgconfigparser[58050]: error at line number 18 at 7.4:syntax error
1月 20 14:40:06 general-work cgconfigparser[58050]: /usr/sbin/cgconfigparser; error loading /etc/cgconfig.d/gopls.conf: Have multiple paths for the same namespace
1月 20 14:40:06 general-work cgconfigparser[58050]: Error: failed to parse file /etc/cgconfig.d/.gopls.conf.swp
1月 20 14:40:06 general-work cgconfigparser[58050]: Error: failed to parse file /etc/cgconfig.d/gopls.conf
1月 20 14:40:06 general-work systemd[1]: Failed to start Control Group configuration service.
```

网上查了查, 没有找到正确解决方法, 有人说是`.conf`文件不能是dos格式, 要转换成unix, 无效.

后来发现是因为我在一个标签页中使用vim编辑, 未退出, 另一个标签页使用`systemctl`重启服务, 导致配置文件加载了两次(`.swp`缓存文件也算), group名称重复而失败.
