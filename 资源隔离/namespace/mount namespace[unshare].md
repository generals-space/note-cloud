# mount namespace

参考文章

1. [Linux Namespace : Mount](https://www.cnblogs.com/sparkdev/p/9424649.html)
    - `unshare --mount`的使用方法与隔离表现
    - `mount`的`--make-shared`, `--make-private`参数的行为表现.
    - `unshare --propagation`选项的作用.

`mount namespace`是用来隔离`mountpoint`挂载点的, 我们创建一个新的mount ns, 这个新的ns会继承父ns的所有挂载点, 但是之后双方再进行挂载/卸载操作时, 就不会再互相影响了.


## 示例1-`unshare --mount`的隔离表现

```
mkdir /demo && sudo chmod 777 /demo && cd $_
mkdir -p iso1/subdir1
mkdir -p iso2/subdir2
## 将 iso1, iso2 整个目录封装进 1.iso, 2.iso 镜像, 里面包含 iso1, iso2 的完整文件系统.
mkisofs -o 1.iso ./iso1
mkisofs -o 2.iso ./iso2

mkdir /mnt/iso1 /mnt/iso2
mount /demo/1.iso /mnt/iso1
```

以下在另一个shell中执行

```
unshare --mount
readlink /proc/$$/ns/mnt
## 此处 umount, 并不会影响到原 namespace 的 iso1 挂载点
umount /mnt/iso1
mount /demo/2.iso /mnt/iso2
```

在两个shell中分别执行`mount | grep iso`, 会发现不同.
