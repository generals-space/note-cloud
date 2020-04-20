# runc-rootless

一般情况下, 普通用户是无法执行`runc`命令的, 因为权限不足, 而且也无法创建和启动容器.

```
$ useradd general
$ su -l general
$ cd /mycontainer
$ ls
config.json  rootfs
$ runc list
ERRO[0000] open /run/runc: permission denied
open /run/runc: permission denied
$ runc create xxx01
ERRO[0000] rootless container requires user namespaces
rootless container requires user namespaces
```

不过`runc`有一个`--rootless`选项, 可以生成能够让普通用户也能`create/start/delete`容器的`spec`配置文件.

```
runc spec --rootless
```

> 注意: 如果此时的`mycontainer`仍然是前面用root用户创建的目录, 仍然会创建失败, 需要使用普通用户自己创建的目录.

```bash
mkdir /mycontainer
cd /mycontainer
mkdir rootfs
## docker 命令需要 root 权限, 所以需要 root 用户先写入, 再修改其属主.
docker export $(docker create busybox) | tar -C rootfs -xf -
chown -R 普通用户:普通用户组 /mycontainer/rootfs
```

接下来可以以普通用户的身份执行了.

```bash
cd /mycontainer
runc spec --rootless
runc --root /tmp/runc run mycontainerid
```

> `--root`指定一个当前普通用户拥有写权限的目录, 用于存储 container 状态.

...失败了?(`--root`的路径和`container id`换了也没用)

```
$ runc --root /tmp/runc run mycontainerid
FATA[0000] nsexec:869 nsenter: failed to unshare user namespace: Invalid argument
FATA[0000] nsexec:724 nsenter: failed to sync with child: next state: Success
ERRO[0000] container_linux.go:346: starting container process caused "process_linux.go:319: getting the final child's pid from pipe caused \"EOF\""
container_linux.go:346: starting container process caused "process_linux.go:319: getting the final child's pid from pipe caused \"EOF\""
```

<???>

