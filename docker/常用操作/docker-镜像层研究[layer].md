参考文章

1. [docker 镜像存储](https://xichuanliang.github.io/container/2024/03/31/docker-%E9%95%9C%E5%83%8F%E5%AD%98%E5%82%A8/)

## docker拉取镜像过程

![](https://gitee.com/generals-space/gitimg/raw/master/2025/9d27eb594b4015e0a517cd24138225b6.png)

登录镜registry，使用`~/.docker/config.json`中的认证信息在registry中进行鉴权，从而拿到token。

向registry发出请求，获取manifest。

```json
// curl -H "Accept:application/vnd.docker.distribution.manifest.v2+json" http://192.168.118.3:5000/v2/nginx/manifests/latest
{
    "schemaVersion": 2,
    "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
    "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": 7656,
        "digest": "sha256:605c77e624ddb75e6110f997c58876baa13f8754486b461117934b24a9dc3a85"
    },
    "layers": [
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 31357624,
            "digest": "sha256:a2abf6c4d29d43a4bf9fbb769f524d0fb36a2edab49819c1bf3e76f409f953ea"
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 25350007,
            "digest": "sha256:a9edb18cadd1336142d6567ebee31be2a03c0905eeefe26cb150de7b0fbc520b"
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 602,
            "digest": "sha256:589b7251471a3d5fe4daccdddfefa02bdc32ffcba0a6d6a2768bf2c401faf115"
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 894,
            "digest": "sha256:186b1aaa4aa6c480e92fbd982ee7c08037ef85114fbed73dbb62503f24c1dd7d"
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 666,
            "digest": "sha256:b4df32aa5a72e2a4316aad3414508ccd907d87b4ad177abd7cbd62fa4dab2a2f"
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 1395,
            "digest": "sha256:a0bcbecc962ed2552e817f45127ffb3d14be31642ef3548997f58ae054deb5b2"
        }
    ]
}
```

读取manifest中image config的digest，这个sha256值就是image ID，根据ID在`/var/lib/docker/image/overlay2/repositories.json`中查看是否存在相同ID的image，如果有就不用从新下载。

每一个镜像对应两个值，image ID（605c77e624ddb75e6110f997c58876baa13f8754486b461117934b24a9dc3a85）和digest（ee89b00528ff4f02f2405e4ee221743ebc3f8e8dd0bfd5c4c20a2fa2aaa7ede3）。image ID是manifest中的image config的digest值，digest是将manifest文件使用sha256计算得出的。

```json
{
    "Repositories": {
        "192.168.118.3:5000/nginx": {
            "192.168.118.3:5000/nginx:latest": "sha256:605c77e624ddb75e6110f997c58876baa13f8754486b461117934b24a9dc3a85",
            "192.168.118.3:5000/nginx@sha256:ee89b00528ff4f02f2405e4ee221743ebc3f8e8dd0bfd5c4c20a2fa2aaa7ede3": "sha256:605c77e624ddb75e6110f997c58876baa13f8754486b461117934b24a9dc3a85"
        },
        "registry": {
            "registry:latest": "sha256:b8604a3fe8543c9e6afc29550de05b36cd162a97aa9b2833864ea8a5be11f3e2",
            "registry@sha256:169211e20e2f2d5d115674681eb79d21a217b296b43374b8e39f97fcf866b375": "sha256:b8604a3fe8543c9e6afc29550de05b36cd162a97aa9b2833864ea8a5be11f3e2"
        }
    }
}
```

如果没有，向registry发出拿到image config的请求。从而获取image config文件。

```json
// curl -H "Accept:application/vnd.docker.distribution.manifest.v2+json" http://192.168.118.3:5000/v2/nginx/blobs/sha256:605c77e624ddb75e6110f997c58876baa13f8754486b461117934b24a9dc3a85
{
    ......
    "rootfs": {
        "diff_ids": [
            "sha256:2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f",
            "sha256:e379e8aedd4d72bb4c529a4ca07a4e4d230b5a1d3f7a61bc80179e8f02421ad8",
            "sha256:b8d6e692a25e11b0d32c5c3dd544b71b1085ddc1fddad08e68cbd7fda7f70221",
            "sha256:f1db227348d0a5e0b99b15a096d930d1a69db7474a1847acbc31f05e4ef8df8c",
            "sha256:32ce5f6a5106cc637d09a98289782edf47c32cb082dc475dd47cbf19a4f866da",
            "sha256:d874fd2bc83bb3322b566df739681fbd2248c58d3369cb25908d68e7ed6040a6"
        ],
        "type": "layers"
    }
}
```

加入diff_ids根据一些列的算法匹配之后，发现一个镜像层都没有，则根据diff_ids中的sha256中的值取拉去对应的层。将拉下来的层进行gzip解压成xxx.tar包，并通过`sha256sum`计算是否和image config中的diff_id相同，如果不相同，拉去失败。

```log
a2abf6c4d29d: Pull complete
a9edb18cadd1: Pull complete
589b7251471a: Pull complete
186b1aaa4aa6: Pull complete
b4df32aa5a72: Pull complete
a0bcbecc962e: Pull complete
Digest: sha256:ee89b00528ff4f02f2405e4ee221743ebc3f8e8dd0bfd5c4c20a2fa2aaa7ede3
Status: Downloaded newer image for 192.168.118.3:5000/nginx:latest
```

根据diff_ids中的sha256的值（diff id），拉取下来的压缩包进行sha256sum后，得到的值，发现与docker pull拉取镜像时的值保持一致。同时，将xxx.tar.gz解压之后成xxx.tar之后，通过sha256sum算法得出的值（digest）在本地存储的值保持一致。

此时就产生了一种对应关系，diff_id与digest的对应关系。**diff_id是layer解压之后的sha256值，digest是layer未解压时的sha256的值。**

```log
curl -H "Accept:application/vnd.docker.image.rootfs.diff.tar.gzip" http://192.168.118.3:5000/v2/nginx/blobs/sha256:a2abf6c4d29d43a4bf9fbb769f524d0fb36a2edab49819c1bf3e76f409f953ea -o layer1.tar.gz

sha256sum layer1.tar.gz
a2abf6c4d29d43a4bf9fbb769f524d0fb36a2edab49819c1bf3e76f409f953ea  layer1.tar.gz

gzip -d layer1.tar.gz
sha256sum layer1.tar
2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f  layer1.tar
```

这种对应关系存储在/var/lib/docker/image/overlay2/distribution/文件中。

```
├── diffid-by-digest（通过digest得到diffid）
│   └── sha256
│       ├── 0d96da54f60b86a4d869d44b44cfca69d71c4776b81d361bc057d6666ec0d878
│       ├── 186b1aaa4aa6c480e92fbd982ee7c08037ef85114fbed73dbb62503f24c1dd7d
│       ├── 3790aef225b922bc97aaba099fe762f7b115aec55a0083824b548a6a1e610719
│       ├── 41af1b5f0f51947706ae712999cf098bef968a7799e7cb4bb2268829e62a6ab3
│       ├── 589b7251471a3d5fe4daccdddfefa02bdc32ffcba0a6d6a2768bf2c401faf115
│       ├── 5b27040df4a23c90c3837d926f633fb327fb3af9ac4fa5d5bc3520ad578acb10
│       ├── 79e9f2f55bf5465a02ee6a6170e66005b20c7aa6b115af6fcd04fad706ea651a
│       ├── 7c457f213c7634afb95a0fb2410a74b7b5bc0ba527033362c240c7a11bef4331
│       ├── a0bcbecc962ed2552e817f45127ffb3d14be31642ef3548997f58ae054deb5b2
│       ├── a2abf6c4d29d43a4bf9fbb769f524d0fb36a2edab49819c1bf3e76f409f953ea
│       ├── a9edb18cadd1336142d6567ebee31be2a03c0905eeefe26cb150de7b0fbc520b
│       ├── b4df32aa5a72e2a4316aad3414508ccd907d87b4ad177abd7cbd62fa4dab2a2f
│       └── e2ead8259a04d39492c25c9548078200c5ec429f628dcf7b7535137954cc2df0
└── v2metadata-by-diffid（通过diffid获得digest）
    └── sha2565
        ├── 2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f
        ├── 32ce5f6a5106cc637d09a98289782edf47c32cb082dc475dd47cbf19a4f866da
        ├── 548a79621a426b4eb077c926eabac5a8620c454fb230640253e1b44dc7dd7562
        ├── 69715584ec78c168981b0925dd7c50f4537bc598dcbce814db2803a10b777b5c
        ├── aa4330046b37f18b2c8266a11687acfcb1912b3312ab6ee427668d9842672d69
        ├── ad10b481abe790a76415269ad68a67e6baeded9586f1aa4d32b22bf60a74e492
        ├── aeccf26589a7bdcad5bbde6d93db4ba6f26dd1fffcae9236e838f8546c2adb9b
        ├── b8d6e692a25e11b0d32c5c3dd544b71b1085ddc1fddad08e68cbd7fda7f70221
        ├── d874fd2bc83bb3322b566df739681fbd2248c58d3369cb25908d68e7ed6040a6
        ├── e379e8aedd4d72bb4c529a4ca07a4e4d230b5a1d3f7a61bc80179e8f02421ad8
        ├── f1db227348d0a5e0b99b15a096d930d1a69db7474a1847acbc31f05e4ef8df8c
        └── f640be0d5aadac5d1376a1ad029edc6caff948d68373888e3007f1422f912fbe
```

最终将xxx.tar包解压成下载完成的layer文件。由此可见，镜像是文件的叠加，最终生成rootfs。

```log
$ tar -xvf layer1.tar
rwxr-xr-x    2 root root      4096 Dec 20  2021 bin
drwxr-xr-x   2 root root         6 Dec 12  2021 boot
drwxr-xr-x   2 root root         6 Dec 20  2021 dev
drwxr-xr-x  30 root root      4096 Dec 20  2021 etc
drwxr-xr-x   2 root root         6 Dec 12  2021 homer
drwxr-xr-x   8 root root        96 Dec 20  2021 lib
drwxr-xr-x   2 root root        34 Dec 20  2021 lib64
drwxr-xr-x   2 root root         6 Dec 20  2021 media
drwxr-xr-x   2 root root         6 Dec 20  2021 mnt
drwxr-xr-x   2 root root         6 Dec 20  2021 opt
drwxr-xr-x   2 root root         6 Dec 12  2021 proc
drwx------   2 root root        37 Dec 20  2021 root
drwxr-xr-x   3 root root        30 Dec 20  2021 run
drwxr-xr-x   2 root root      4096 Dec 20  2021 sbin
drwxr-xr-x   2 root root         6 Dec 20  2021 srv
drwxr-xr-x   2 root root         6 Dec 12  2021 sys
drwxrwxrwt   2 root root         6 Dec 20  2021 tmp
drwxr-xr-x  11 root root       120 Dec 20  2021 usr
drwxr-xr-x  11 root root       139 Dec 20  2021 var
```

那么最终这些xxx.tar包会怎么解压到哪里呢，总不可能随便解压一下。

## 镜像在本地存储

![](https://gitee.com/generals-space/gitimg/raw/master/2025/6dbd326025116c2131e801cc41858c36.png)

在`/var/lib/docker/image/overlay2/layerdb`中记录layer的存储信息，但并不真正的存储layer。

在该文件中，发现有一个文件目录和rootfs中的最底层的diff_id（2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f）一致，且仅一个。

查阅资料可知，在/var/lib/docker/image/overlay2/layerdb/sha256中的每一个文件夹与diff_id保持一致，在此处该id成为chainid。不过对应关系为：除最底层的diff_id与chainid保持一致外，其余层均依赖于下一层的chainid。chainid用来连接镜像中的各个层

```
.
├── mounts
│   └── 9610566a284373f2da4e0c596b8360e4e5854d3aa697798aac43ac081036f311
├── sha256
│   ├── 02b80ac2055edd757a996c3d554e6a8906fd3521e14d1227440afd5163a5f1c4
│   ├── 2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f
│   ├── 780238f18c540007376dd5e904f583896a69fe620876cabc06977a3af4ba4fb5
│   ├── 7850d382fb05e393e211067c5ca0aada2111fcbe550a90fed04d1c634bd31a14
│   ├── b625d8e29573fa369e799ca7c5df8b7a902126d2b7cbeb390af59e4b9e1210c5
│   ├── b92aa5824592ecb46e6d169f8e694a99150ccef01a2aabea7b9c02356cdabe7c
└── tmp
```

计算关系：本层chainid = "下一层的chainID" + "本层的diff_id" | sha256sum

```bash
echo -n "sha256:2edcec3590a4ec7f40cf0743c15d78fb39d8326bc029073b41ef9727da6c851f sha256:e379e8aedd4d72bb4c529a4ca07a4e4d230b5a1d3f7a61bc80179e8f02421ad8" | sha256sum
780238f18c540007376dd5e904f583896a69fe620876cabc06977a3af4ba4fb5  -

echo -n "sha256:780238f18c540007376dd5e904f583896a69fe620876cabc06977a3af4ba4fb5 sha256:b8d6e692a25e11b0d32c5c3dd544b71b1085ddc1fddad08e68cbd7fda7f70221" | sha256sum
b92aa5824592ecb46e6d169f8e694a99150ccef01a2aabea7b9c02356cdabe7c  -

echo -n "sha256:b92aa5824592ecb46e6d169f8e694a99150ccef01a2aabea7b9c02356cdabe7c sha256:f1db227348d0a5e0b99b15a096d930d1a69db7474a1847acbc31f05e4ef8df8c" | sha256sum
02b80ac2055edd757a996c3d554e6a8906fd3521e14d1227440afd5163a5f1c4  -

echo -n "sha256:02b80ac2055edd757a996c3d554e6a8906fd3521e14d1227440afd5163a5f1c4 sha256:32ce5f6a5106cc637d09a98289782edf47c32cb082dc475dd47cbf19a4f866da" | sha256sum
7850d382fb05e393e211067c5ca0aada2111fcbe550a90fed04d1c634bd31a14  -

echo -n "sha256:7850d382fb05e393e211067c5ca0aada2111fcbe550a90fed04d1c634bd31a14 sha256:d874fd2bc83bb3322b566df739681fbd2248c58d3369cb25908d68e7ed6040a6" | sha256sum
b625d8e29573fa369e799ca7c5df8b7a902126d2b7cbeb390af59e4b9e1210c5  -
```

最终，计算出最终的chainid是 b625d8e29573fa369e799ca7c5df8b7a902126d2b7cbeb390af59e4b9e1210c5，进入这个文件目录可知

```
├── cache-id 真正对应的layer数据的目录
├── diff     该层的diffid
├── parent   上一层的chainid
├── size     该层的大小
└── tar-split.json.gz tar-split.json.gz，layer压缩包的split文件，通过这个文件可以还原layer的tar包，https://github.com/vbatts/tar-split
```

获取cache-id中的值

```bash
cat cache-id
3f4cb9effac5ec0d172fd92f3cd932460a785bf1b938ed8bfa5913081664003a
```

在对应的文件查找diff层。里面就是layer中真正的内容。

```
cd /var/lib/docker/overlay2/3f4cb9effac5ec0d172fd92f3cd932460a785bf1b938ed8bfa5913081664003a

cat lower
l/5COCTCQVW7JD7RE3WW62IWN3GV:l/M4YQYVZR2O2DCI6PNMWNVEGX3H:l/ZD6YRQ3T5ZXQKZXYP6MFCJD546:l/VE5FXRY5J4AEZDFLK6ILMFMU3E:l/PJT35YHDQWZPE3AEWI4BGHYSYC
```

```
.
├── 0131057585a4730358d8118159dba84e4343d2a99b20be8ec2935ac747774cde
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── 1662b3b65aae2285418da9fdad8abd6fcaad1290a7c608b6cfe3d01ecfe16f0e
│   ├── diff
│   ├── link
│   ├── lower
│   ├── merged
│   └── work
├── 1662b3b65aae2285418da9fdad8abd6fcaad1290a7c608b6cfe3d01ecfe16f0e-init
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── 21147703e267c0a9741afc92b6d84479b6d457636a2be5a84cb146ee4ee78640
│   ├── committed
│   ├── diff
│   └── link
├── 3f4cb9effac5ec0d172fd92f3cd932460a785bf1b938ed8bfa5913081664003a
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── 76383a05cead72c5aa8045be7c6dcce847621d08ca12571ed191155e09ae1146
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── 7c01c7f54d97e5890f8e4dc66849d5b4fd00c70748d35b9fe81ccaf6c637be27
│   ├── committed
│   ├── diff
│   └── link
├── b53dd778c39441181d987ee1be7f4b9fdcfa15d6eea47af0a4f31d8423925819
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── backingFsBlockDev
├── c0e68f2807d39963736c5050d63a5b805dd2d7cde1be46867f9df123e7c15e12
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── c5e9060bd948b0a4bfdd86754723bc5200fde4bee34d0c25a66e2a5414d2f11e
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── c8389c704e9938c6b7d01425eaf98dc32074d2b5eddb82536fd6b49757427d1c
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── cb4ad7aecae9ecb47e6b68b1d58a1e71a386f33cfd51d14d61ed4e3a56fccb5b
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
├── e9b1b38b87553b7995f804f6de2cdccad1ce4e24921b7ea4e8a815fdeee8cf71
│   ├── committed
│   ├── diff
│   ├── link
│   ├── lower
│   └── work
└── l
    ├── 2UJXL2KXNRPJH2MQTNOUISBQUN -> ../1662b3b65aae2285418da9fdad8abd6fcaad1290a7c608b6cfe3d01ecfe16f0e/diff
    ├── 5COCTCQVW7JD7RE3WW62IWN3GV -> ../c8389c704e9938c6b7d01425eaf98dc32074d2b5eddb82536fd6b49757427d1c/diff
    ├── ALXHZKACOJCUZ72DJ5G5DQOQCF -> ../76383a05cead72c5aa8045be7c6dcce847621d08ca12571ed191155e09ae1146/diff
    ├── C4NZT2GBJVSMOZ4QBG3RUI4EUH -> ../cb4ad7aecae9ecb47e6b68b1d58a1e71a386f33cfd51d14d61ed4e3a56fccb5b/diff
    ├── EJ7VTNUSQVI4FCHGV7SPZAUVZ5 -> ../1662b3b65aae2285418da9fdad8abd6fcaad1290a7c608b6cfe3d01ecfe16f0e-init/diff
    ├── GL3YM5S5K63VTKYNUPOMCEPSFF -> ../b53dd778c39441181d987ee1be7f4b9fdcfa15d6eea47af0a4f31d8423925819/diff
    ├── GRN25DOYR7MLG5L6ROW444UWIW -> ../0131057585a4730358d8118159dba84e4343d2a99b20be8ec2935ac747774cde/diff
    ├── M4YQYVZR2O2DCI6PNMWNVEGX3H -> ../c5e9060bd948b0a4bfdd86754723bc5200fde4bee34d0c25a66e2a5414d2f11e/diff
    ├── PJT35YHDQWZPE3AEWI4BGHYSYC -> ../7c01c7f54d97e5890f8e4dc66849d5b4fd00c70748d35b9fe81ccaf6c637be27/diff
    ├── TBJFP5EWBIA6YHDRDUNAVUTPFA -> ../3f4cb9effac5ec0d172fd92f3cd932460a785bf1b938ed8bfa5913081664003a/diff
    ├── VE5FXRY5J4AEZDFLK6ILMFMU3E -> ../c0e68f2807d39963736c5050d63a5b805dd2d7cde1be46867f9df123e7c15e12/diff
    ├── WD6DDEX6EHVPLTCDIU2MESMVSS -> ../21147703e267c0a9741afc92b6d84479b6d457636a2be5a84cb146ee4ee78640/diff
    └── ZD6YRQ3T5ZXQKZXYP6MFCJD546 -> ../e9b1b38b87553b7995f804f6de2cdccad1ce4e24921b7ea4e8a815fdeee8cf71/diff
```

## 镜像在容器中的使用

![](https://gitee.com/generals-space/gitimg/raw/master/2025/36690a304a3d1d913648713da92c86cb.png)

运行一个镜像时，会在返回容器ID，并在`/var/lib/docker/image/overlay2/layerdb/mounts/`中产生相应文件，记录容器启动所用的镜像层信息、init层信息、r/w层信息。

```
/var/lib/docker/image/overlay2/layerdb/mounts/9610566a284373f2da4e0c596b8360e4e5854d3aa697798aac43ac081036f311

├── init-id  init层的id
│   ├── 1662b3b65aae2285418da9fdad8abd6fcaad1290a7c608b6cfe3d01ecfe16f0e-init
├── mount-id r/w层的id
│   ├── 1662b3b65aae2285418da9fdad8abd6fcaad1290a7c608b6cfe3d01ecfe16f0e
└── parent   镜像层的chainid，根据chainid获取各个镜像层。
│   ├── sha256:fd39b5678fdb70fc98ac5e6b4e7383f1b74b7f1a08bc6dd74fadadcc4beaf364
```

根据parent的chainid获取镜像层。

```
cd /var/lib/docker/image/overlay2/layerdb/sha256/fd39b5678fdb70fc98ac5e6b4e7383f1b74b7f1a08bc6dd74fadadcc4beaf364

├── cache-id
│   ├── b53dd778c39441181d987ee1be7f4b9fdcfa15d6eea47af0a4f31d8423925819
├── diff
├── parent
├── size
└── tar-split.json.gz
```

从cache-id中查找对应的layer，根据lower中的链接获取其中的diff文件，最终组成镜像层。

```
cd /var/lib/docker/overlay2/b53dd778c39441181d987ee1be7f4b9fdcfa15d6eea47af0a4f31d8423925819
├── committed
├── diff
│   └── entrypoint.sh
├── link
├── lower
│   ├── b53dd778c39441181d987ee1be7f4b9fdcfa15d6eea47af0a4f31d8423925819l/
C4NZT2GBJVSMOZ4QBG3RUI4EUH:l/GRN25DOYR7MLG5L6ROW444UWIW:l/ALXHZKACOJCUZ72DJ5G5DQOQCF:l/WD6DDEX6EHVPLTCDIU2MESMVSS
└── work
```

## 镜像在registry的存储

![](https://gitee.com/generals-space/gitimg/raw/master/2025/12c40c4e87678647b9e1b882214cf626.png)

registry在本地的挂载路径为`/usr/local/image_registry`, 在registry中分为两大部分，一部分是blobs文件，另一部分是repositories文件。

blobs存储的是image manifest、image config、image layer（image的各个层）。当把相应image layer拉下来之后，使用docker-untar进程将data文件解压后的数据存放在`/var/lib/docker/overlay2/${digest}/diff`中。

repositories存储的是镜像的版本信息、历史信息等等。

- _manifests文件是在镜像上传完成之后由registry生成的，并且该目录下的文件都是一个名为link的文本，link中的文本指向blobs中的 blob digest目录。_manifests文件下包含镜像的revision和tags信息，每一个镜像的每一个tag对应着tag名相同的目录。镜像的tag并不存储在image config中，而是以目录的形式形成镜像的tag。每一个tag文件下包含current和index目录，current中的link保存了该tag当前manifest的digests信息，index中列出了该tag历史上传的所有版本的sha256信息。revision目录中存放了该repository历史上传版本的所有sha256信息。
- _layers存放镜像的image config、image layers的sha256的信息。
- _uploads是个临时文件，主要用来存放push镜像过程中的文件数据，当镜像layer上传完成之后回清空该文件夹。其中data文件上传完毕后会转移到blobs目录下，根据该文件的sha256值散列存储到相应的目录下。

## 总结

### 镜像在registry中的存储

registry中有两个文件夹，blobs和repositories。
blobs中存储image manifest、image config、image layers。即存储的是镜像的真正内容。
repositories中存储的是镜像的版本信息，tag信息等等。通过文件的link指引blob文件中真正的内容。

### 拉取镜像

获取registry的认证信息

发送url，url中包含了镜像名、tag。通过这个url去对应的repositories文件中查找对应的image manifest。将manifest返回回去。

获取到manifest之后，根据manifest中的image config对应的sha256值，去/var/lib/docker/image/overlay2/repositories.json中查找是否有相同的值，如果有就不用后续下载。如果没有就根据image config的sha256的值去registry中下载image config。

获取image config后，查看diff_ids。根据diff_ids在本地找对应的layer是否存在。

如果layer不存在，就根据diff_ids中的sha256值去registry中并行下载layer。

下载之后，进行解压（gzip），判断解压后的xxx.tar包的sha256值与image config中的sha256值是否相同，相同下载成功。

当所有的layer都下载完成之后，再将xxx.tar进行解压，解压到相应的/var/lib/docker/overlay2/${digest}/diff中。

### 推送镜像

获取registry的认证信息

向registry中发送POST /v2//blobs/uploads/ 请求，检查registry中是否已经存在镜像的layer。

客户端通过URL上传layer数据，上传镜像layer分为整块上传和分块上传。

上传完成之后，docker向registry发送求情，告知layer已经上传完毕。

上传完成之后，将manifest上传上去。

### 镜像在本地存储

在/var/lib/docker/image/overlay2/repositories.json中有对应的镜像名称和相应的值（image id、digest）。image ID是manifest中的image config的digest值，digest是将manifest文件使用sha256计算得出的。

在/var/lib/docker/image/overlay2/distribution/文件夹中有两个文件（diffid-by-digest、v2metadata-by-diffid），这两个文件中记录了diff_id和digest之间的对应关系。可根据diff_id和digest之间的关系，计算是否存在该layer以及是否下载的layer是否正确。

在/var/lib/docker/image/overlay2/layerdb/sha256中存放了chainid与layer之间的真正关系。通过最底层的diff_id与chainid相同，并且通过每一层的diff_id和上一层的chainid就能够计算出本层的chainid。进入对应的chaini目录中，查找对应的cache-id以获取layer真正的地方。

根据chainid获取每层的cache-id之后，在/var/lib/docker/overlay2中根据cache-id的值查找对应的diff文件，diff文件中的文件夹就是xxx.tar解压之后，也就是layer真实的值。

### 镜像在容器中

镜像存储在本地，转化为了filesystem bundle

当启动容器时，生成一个容器ID，并在`/var/lib/docker/image/overlay2/layerdb/mounts/容器ID`生成目录。
在该目录下有三个文件，init-id、mount-id、parent。

- init-id：容器的init层所在id
- mount-id：容器的r/w层所在id
- parent：容器的镜像层所在chainid

init-id和mount-id的文件在/var/lib/docker/overlay2/中

parent中的镜像层，是通过chainid获取镜像的每层文件。

最终使用overlayFS(联合挂载)技术，将lowerdir设置为image layer（parent）和init layer（init-id），将upperdir层设置为r/w层（mount-id），workerdir设置为挂载之后的工作目录。最终联合挂载在merger目录中。

联合挂载技术的特点：

- 上下合并时，上层文件覆盖下层同名文件。
- 写时复制。删除的文件是upper的，并且这个文件在lower层不存在直接删除。删除的文件来自于lower层，upper层没有对应的文件，overlayFS通过whiteout机制，屏蔽文件，并不真正删除文件。
