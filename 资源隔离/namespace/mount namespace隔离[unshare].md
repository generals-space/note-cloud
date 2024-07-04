# mount namespace

参考文章

1. [Linux Namespace : Mount](https://www.cnblogs.com/sparkdev/p/9424649.html)
    - `unshare --mount`的使用方法与隔离表现
    - `mount`的`--make-shared`, `--make-private`参数的行为表现.
    - `unshare --propagation`选项的作用.
2. [黄东升: mount namespace和共享子树](https://cloud.tencent.com/developer/article/1518101)
    - 共享子树(Shared subtrees)
    - 四种传递类型：
        - MS_SHARED
        - MS_PRIVATE
        - MS_SLAVE
        - MS_UNBINDABLE
    - Peer groups（对等组）
3. [Mount namespaces, mount propagation, and unbindable mounts](https://blog.csdn.net/u012319493/article/details/102887094)
    - MS_SHARED 和 MS_PRIVATE 示例
    - MS_SLAVE 示例
    - MS_UNBINDABLE 示例
4. [Building a container by hand using namespaces: The mount namespace](https://www.redhat.com/sysadmin/mount-namespaces)
5. [Linux Namespace系列（04）：mount namespaces (CLONE_NEWNS)](https://segmentfault.com/a/1190000006912742)
    - mount ns 是第一个被加入Linux的 ns, 由于当时没想到还会引入其它的 ns, 所以取名为`CLONE_NEWNS`, 而没有叫`CLONE_NEWMOUNT`
    - mount ns 用来隔离文件系统的挂载点, 使得不同的 mount ns 拥有自己独立的挂载点信息, 不同的 ns 之间不会相互影响, 这对于构建用户或者容器自己的文件系统目录非常有用
    - 当前进程所在 mount ns 里的所有挂载信息可以在`/proc/[pid]/mounts`、`/proc/[pid]/mountinfo`和`/proc/[pid]/mountstats`里面找到

`mount ns`是用来隔离`mountpoint`挂载点的, 我们创建一个新的 mount ns, 这个新的ns会继承父ns的所有挂载点, 但是之后双方再进行挂载/卸载操作时, 就不会再互相影响了.

## 1. `unshare --mount`的隔离表现

示例来自参考文章5.

```bash
mkdir /demo && sudo chmod 777 /demo && cd /demo
mkdir -p iso1/subdir1
mkdir -p iso2/subdir2
## 将 iso1, iso2 整个目录封装进 1.iso, 2.iso 镜像, 里面包含 iso1, iso2 的完整文件系统.
mkisofs -o 1.iso ./iso1
mkisofs -o 2.iso ./iso2

mkdir /mnt/iso1 /mnt/iso2
mount /demo/1.iso /mnt/iso1
```

在当前 mount ns 下, 挂载了 iso1, 而没有挂载 iso2, 执行`mount`命令可以看到.

```log
$ mount | grep iso
/demo/1.iso on /mnt/iso1 type iso9660 (ro,relatime)
```

------

以下在另一个shell中执行

```bash
unshare --mount
readlink /proc/$$/ns/mnt
## 此处 umount, 并不会影响到原 namespace 的 iso1 挂载点
umount /mnt/iso1
mount /demo/2.iso /mnt/iso2
```

在执行完`unshare`与`readlink`命令后, 执行`mount`命令, 可以看到与原 mount ns 一样的挂载情况, 因为新的 ns 继承了原有的 ns.

```log
$ unshare --mount
$ readlink /proc/$$/ns/mnt
$ mount | grep iso
/demo/1.iso on /mnt/iso1 type iso9660 (ro,relatime)
```

然后在新的 mount ns 中, 执行`umount`与`mount`, 卸载 iso1, 挂载 iso2, 现在再查看`mount`情况.

```log
$ umount /mnt/iso1

$ mount /demo/2.iso /mnt/iso2
mount: /dev/loop1 is write-protected, mounting read-only
$ mount | grep iso
/demo/2.iso on /mnt/iso2 type iso9660 (ro,relatime)
```

此时, 回到第1个shell(原 mount ns), 再次查看`mount`情况.

```log
$ mount | grep iso
/demo/1.iso on /mnt/iso1 type iso9660 (ro,relatime)
```

没有变化, 可以得出, 两个 mount ns 之间不会相互影响.

## 2. 不同 mount ns 挂载到同一目录

单纯从挂载点角度看好像没什么用, "挂载点"的概念还是太抽象. 其实"挂载点"就是挂载的"目标目录", 在同一个 mount ns 下, 将不同的源目录挂载到同一个"目标目录", "目标目录"的内容会被相互覆盖. 而如果创建独立的 mount ns 进行挂载, 则会出现, 在不同的 mount ns 中的同一目录下, 内容不同的情况.

### 2.1 同一 mount ns 下多次 mount 相互覆盖

我们沿用上一个示例的命令.

```bash
mkdir /demo && sudo chmod 777 /demo && cd /demo
mkdir -p iso1/subdir1
mkdir -p iso2/subdir2
## 将 iso1, iso2 整个目录封装进 1.iso, 2.iso 镜像, 里面包含 iso1, iso2 的完整文件系统.
mkisofs -o 1.iso ./iso1
mkisofs -o 2.iso ./iso2
## 从此处开始不同, 只创建一个目录.
mkdir /mnt/iso
```

先 mount 一次, 看看效果.

```log
$ mount /demo/1.iso /mnt/iso
mount: /dev/loop2 is write-protected, mounting read-only
/mnt
$ ls /mnt/iso
subdir1
/mnt
$ mount | grep iso
/demo/1.iso on /mnt/iso type iso9660 (ro,relatime)
```

在当前 shell 再挂载一次.

```log
$ mount /demo/2.iso /mnt/iso
mount: /dev/loop3 is write-protected, mounting read-only
/mnt
$ ll /mnt/iso
total 2
dr-xr-xr-x 1 root root 2048 Jan 16 22:28 subdir2
/mnt
$ mount | grep iso
/demo/1.iso on /mnt/iso type iso9660 (ro,relatime)
/demo/2.iso on /mnt/iso type iso9660 (ro,relatime)
```

`/mnt/iso`目录中的内容被替换成了`iso2`的内容.

不过`mount`显示的挂载信息却同时包含`iso1`和`iso2`, 很神奇, 这样的话, `umount`也需要执行2次才行.

```log
$ mount | grep iso
/demo/1.iso on /mnt/iso type iso9660 (ro,relatime)
/demo/2.iso on /mnt/iso type iso9660 (ro,relatime)

$ umount /mnt/iso
$ mount | grep iso
/demo/1.iso on /mnt/iso type iso9660 (ro,relatime)

$ umount /mnt/iso
/mnt
$ mount | grep iso
```

### 2.2 不同 mount ns 下分别 mount 各自独立

现在重新来一次, 在当前 ns 下, 挂载`1.iso`到`/mnt/iso`.

```log
$ mount /demo/1.iso /mnt/iso
mount: /dev/loop2 is write-protected, mounting read-only

$ ls /mnt/iso
subdir1

$ mount | grep iso
/demo/1.iso on /mnt/iso type iso9660 (ro,relatime)
```

ok, 接下来在另一个 shell 中, 创建独立的 mount ns, 并将`2.iso`也挂载到`/mnt/iso`目录. 当然, 要先卸载继承上一个 ns 的`iso1`挂载点.

```log
$ unshare --mount
$ readlink /proc/$$/ns/mnt
$ mount | grep iso
/demo/1.iso on /mnt/iso type iso9660 (ro,relatime)

$mount /demo/2.iso /mnt/iso
mount: /dev/loop3 is write-protected, mounting read-only

$ls /mnt/iso
subdir2
```

回到前一个 mount ns 的 shell, 查看`/mnt/iso`的内容, 会发现还是`iso1`的内容. 

**在不同 mount ns 中, 同一个目录的内容是可以不同的, 而且手动修改也互不影响.**
