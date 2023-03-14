A: 172.16.156.172
B: 172.16.42.38


┌-------------------------------------------------┐ 
|          10.254.0.2/24     10.254.0.3/24        | 
|           ┌--------┐        ┌--------┐          | 
|           |  eht0  |        |  eht0  |          | 
|           └----┬---┘        └----┬---┘          | 
|                └────────┬────────┘              | 
|                         |                       | 
|  ┌-----------┐    ┌-----┴-----┐                 | 
|  | flannel.1 ├----┤    cni0   | 10.254.0.1/24   | 
|  └-----------┘    └-----------┘   (bridge)      | 
|  10.254.0.0/32                ┌--------┐        | 
|     (vxlan)                   |  eth0  |        | 
|                               └----┬---┘        | 
|                  172.16.91.128/24  |            | 
└------------------------------------|------------┘ 
                                     |
                                     |     +----------------+      |
                                     └────>| 172.16.91.1/24 |<─────┘
                                           +----------------+
                                               网关/路由器


### vm1

#### 路由

```console
$ ip r
default via 172.16.159.253 dev eth0
10.254.0.0/24 dev cni0 proto kernel scope link src 10.254.0.1
10.254.1.0/24 via 10.254.1.0 dev flannel.1 onlink
169.254.0.0/16 dev eth0 scope link metric 1002
172.16.144.0/20 dev eth0 proto kernel scope link src 172.16.156.172
```

#### 转发表

```console
$ bridge fdb
1a:7e:bf:63:36:ad dev flannel.1 dst 172.16.42.38 self permanent
```

其中`1a:7e:bf:63:36:ad`是vm2上`flannel.1`接口的mac地址.

### vm2

#### 路由

```console
$ ip r
default via 172.16.47.253 dev eth0
10.254.0.0/24 via 10.254.0.0 dev flannel.1 onlink
10.254.1.0/24 dev cni0 proto kernel scope link src 10.254.1.1
169.254.0.0/16 dev eth0 scope link metric 1002
```

#### 转发表

```
$ bridge fdb
1e:08:95:6b:9c:ec dev flannel.1 dst 172.16.156.172 self permanent
```

`1e:08:95:6b:9c:ec`是vm1上`flannel.1`接口的mac地址.
