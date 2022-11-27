# Service NodePort 只能通过Pod运行的主机访问

参考文章

1. [NodePort 只能在node节点上访问，外部无法访问](http://dockone.io/question/1504)
    - 在各主机上执行`iptables -P FORWARD ACCEPT`即可
2. [设置service的nodeport以后外部无法访问对应的端口的问题](https://blog.51cto.com/11288550/2378289)
3. [kubernetes: cannot access NodePort from other machines](https://stackoverflow.com/questions/46667659/kubernetes-cannot-access-nodeport-from-other-machines)

## 场景描述

部署了 NodePort 的 Service, 假设端口为 30080, 但是在集群外(同一局域网)却无法通过`宿主机IP:30080`访问目录服务. 详细的排查过程

1. 服务正常, `exec`进入集群中某一Pod, 通过`PodIP:Pod内部端口`可以访问
2. 通过`netstat`查看, 各主机节点上都监听了`30080`端口, 由`kube-proxy`服务维护
3. 登录某一宿主机节点, 访问自身的`宿主机IP:30080`可以访问到该服务
4. 局域网内其他主机都无法通过`集群中宿主IP:30080`访问到服务, **只有该服务对应的 Pod 所在的主机可以被正常访问**
5. 集群内的节点只能访问自身的`30080`, 互相也不能访问(除了Pod所在的主机)

按照参考文章1中所说, 是因为 iptables 配置中`Forward`链有问题, 但是我在排查的时候`policy`明明是`ACCEPT`, 并不是`DROP`, 现在没时间, 以后再具体分析吧.

宿主机上的 iptables 配置如下

```
[root@dev-k8s-node4 ~]# iptables -nvL
Chain INPUT (policy ACCEPT 6801 packets, 2908K bytes)
 pkts bytes target     prot opt in     out     source               destination
 101M   41G cali-INPUT  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:Cz_u1IQiXIMmKD4c */
 210K   13M KUBE-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            ctstate NEW /* kubernetes service portals */
 210K   13M KUBE-EXTERNAL-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            ctstate NEW /* kubernetes externally-visible service portals */
6300K 3705M KUBE-FIREWALL  all  --  *      *       0.0.0.0/0            0.0.0.0/0

Chain FORWARD (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
 174M   68G cali-FORWARD  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:wUHhoiAYhphO9Mso */
   10   752 KUBE-FORWARD  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes forwarding rules */
    0     0 KUBE-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            ctstate NEW /* kubernetes service portals */
    0     0 DOCKER-USER  all  --  *      *       0.0.0.0/0            0.0.0.0/0
    0     0 DOCKER-ISOLATION-STAGE-1  all  --  *      *       0.0.0.0/0            0.0.0.0/0
    0     0 ACCEPT     all  --  *      docker0  0.0.0.0/0            0.0.0.0/0            ctstate RELATED,ESTABLISHED
    0     0 DOCKER     all  --  *      docker0  0.0.0.0/0            0.0.0.0/0
    0     0 ACCEPT     all  --  docker0 !docker0  0.0.0.0/0            0.0.0.0/0
    0     0 ACCEPT     all  --  docker0 docker0  0.0.0.0/0            0.0.0.0/0

Chain OUTPUT (policy ACCEPT 14315 packets, 1258K bytes)
 pkts bytes target     prot opt in     out     source               destination
  83M   33G cali-OUTPUT  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:tVnHkvAo15HuiPy0 */
 425K   25M KUBE-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            ctstate NEW /* kubernetes service portals */
6571K  543M KUBE-FIREWALL  all  --  *      *       0.0.0.0/0            0.0.0.0/0

Chain DOCKER (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain DOCKER-ISOLATION-STAGE-1 (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 DOCKER-ISOLATION-STAGE-2  all  --  docker0 !docker0  0.0.0.0/0            0.0.0.0/0
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0

Chain DOCKER-ISOLATION-STAGE-2 (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 DROP       all  --  *      docker0  0.0.0.0/0            0.0.0.0/0
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0

Chain DOCKER-USER (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0

Chain KUBE-EXTERNAL-SERVICES (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain KUBE-FIREWALL (2 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes firewall for dropping marked packets */ mark match 0x8000/0x8000

Chain KUBE-FORWARD (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            ctstate INVALID
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes forwarding rules */ mark match 0x4000/0x4000
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes forwarding conntrack pod source rule */ ctstate RELATED,ESTABLISHED
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes forwarding conntrack pod destination rule */ ctstate RELATED,ESTABLISHED

Chain KUBE-KUBELET-CANARY (0 references)
 pkts bytes target     prot opt in     out     source               destination

Chain KUBE-PROXY-CANARY (0 references)
 pkts bytes target     prot opt in     out     source               destination

Chain KUBE-SERVICES (3 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-FORWARD (1 references)
 pkts bytes target     prot opt in     out     source               destination
 174M   68G MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:vjrMJCRpqwy5oRoX */ MARK and 0xfff1ffff
 174M   68G cali-from-hep-forward  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:A_sPAO0mcxbT9mOV */ mark match 0x0/0x10000
  78M   31G cali-from-wl-dispatch  all  --  cali+  *       0.0.0.0/0            0.0.0.0/0            /* cali:8ZoYfO5HKXWbB3pk */
  95M   36G cali-to-wl-dispatch  all  --  *      cali+   0.0.0.0/0            0.0.0.0/0            /* cali:jdEuaPBe14V2hutn */
 539K   34M cali-to-hep-forward  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:12bc6HljsMKsmfr- */
 539K   34M ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:MH9kMp5aNICL-Olv */ /* Policy explicitly accepted packet. */ mark match 0x10000/0x10000

Chain cali-INPUT (1 references)
 pkts bytes target     prot opt in     out     source               destination
  94M   38G ACCEPT     4    --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:PajejrV4aFdkZojI */ /* Allow IPIP packets from Calico hosts */ match-set cali40all-hosts-net src ADDRTYPE match dst-type LOCAL
    0     0 DROP       4    --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:_wjq-Yrma8Ly1Svo */ /* Drop IPIP packets from non-Calico hosts */
 299K   25M cali-wl-to-host  all  --  cali+  *       0.0.0.0/0            0.0.0.0/0           [goto]  /* cali:8TZGxLWh_Eiz66wc */
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:6McIeIDvPdL6PE1T */ mark match 0x10000/0x10000
6297K 3703M MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:YGPbrUms7NId8xVa */ MARK and 0xfff0ffff
6297K 3703M cali-from-host-endpoint  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:2gmY7Bg2i0i84Wk_ */
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:q-Vz2ZT9iGE331LL */ /* Host endpoint policy accepted packet. */ mark match 0x10000/0x10000

Chain cali-OUTPUT (1 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:Mq1_rAdXXH3YkrzW */ mark match 0x10000/0x10000
 323K   25M RETURN     all  --  *      cali+   0.0.0.0/0            0.0.0.0/0            /* cali:69FkRTJDvD5Vu6Vl */
  77M   33G ACCEPT     4    --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:AnEsmO6bDZbQntWW */ /* Allow IPIP packets to other Calico hosts */ match-set cali40all-hosts-net dst ADDRTYPE match src-type LOCAL
6245K  518M MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:9e9Uf3GU5tX--Lxy */ MARK and 0xfff0ffff
6245K  518M cali-to-host-endpoint  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:OB2pzPrvQn6PC89t */
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:tvSSMDBWrme3CUqM */ /* Host endpoint policy accepted packet. */ mark match 0x10000/0x10000

Chain cali-failsafe-in (0 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:wWFQM43tJU7wwnFZ */ multiport dports 22
    0     0 ACCEPT     udp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:LwNV--R8MjeUYacw */ multiport dports 68
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:QOO5NUOqOSS1_Iw0 */ multiport dports 179
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:cwZWoBSwVeIAZmVN */ multiport dports 2379
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:7FbNXT91kugE_upR */ multiport dports 2380
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:ywE9WYUBEpve70WT */ multiport dports 6666
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:l-WQSVBf_lygPR0J */ multiport dports 6667

Chain cali-failsafe-out (0 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 ACCEPT     udp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:82hjfji-wChFhAqL */ multiport dports 53
    0     0 ACCEPT     udp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:TNM3RfEjbNr72hgH */ multiport dports 67
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:ycxKitIl4u3dK0HR */ multiport dports 179
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:hxjEWyxdkXXkdvut */ multiport dports 2379
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:cA_GLtruuvG88KiO */ multiport dports 2380
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:Sb1hkLYFMrKS6r01 */ multiport dports 6666
    0     0 ACCEPT     tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:UwLSebGONJUG4yG- */ multiport dports 6667

Chain cali-from-hep-forward (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-from-host-endpoint (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-from-wl-dispatch (2 references)
 pkts bytes target     prot opt in     out     source               destination
1536K  636M cali-fw-cali291dce4466b  all  --  cali291dce4466b *       0.0.0.0/0            0.0.0.0/0           [goto]  /* cali:vkm68zJCyubKYfyX */
 9172  943K cali-fw-calie50f5cd5607  all  --  calie50f5cd5607 *       0.0.0.0/0            0.0.0.0/0           [goto]  /* cali:t453q5rE8nmXIIpy */
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:mKJD_Zhs63CJcduS */ /* Unknown interface */

Chain cali-fw-cali291dce4466b (1 references)
 pkts bytes target     prot opt in     out     source               destination
  60M   25G ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:GUWRhz0gvc8WDxPo */ ctstate RELATED,ESTABLISHED
  200 10400 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:DBloaTnRWp1UlKXk */ ctstate INVALID
 115K 8174K MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:TisMqKpa8uSRerLb */ MARK and 0xfffeffff
    0     0 DROP       udp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:iowOjT3c15hqh8Nd */ /* Drop VXLAN encapped packets originating in pods */ multiport dports 4789
    0     0 DROP       4    --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:MqireO6Bik3CVqyn */ /* Drop IPinIP encapped packets originating in pods */
 115K 8174K cali-pro-kns.default  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:6RdwKgWR7iQ8MasR */
 115K 8174K RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:MnPBr6Fk1oBdF8qj */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 cali-pro-ksa.default.default  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:RPnJ3X3xAYUkdrBz */
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:3NzjRg4l1r5OdcNe */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:2rCy18RgS5G25Jrf */ /* Drop if no profiles matched */

Chain cali-fw-calie50f5cd5607 (1 references)
 pkts bytes target     prot opt in     out     source               destination
 9171  943K ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:Tn-ExPhP6VLBKDmp */ ctstate RELATED,ESTABLISHED
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:uzQwrH61KgvQLDqG */ ctstate INVALID
    1    60 MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:MS3ujCLFwlGu3yMg */ MARK and 0xfffeffff
    0     0 DROP       udp  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:luwFiSAd71iG85tj */ /* Drop VXLAN encapped packets originating in pods */ multiport dports 4789
    0     0 DROP       4    --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:zG1ZgwLDMKcr9_AO */ /* Drop IPinIP encapped packets originating in pods */
    1    60 cali-pro-kns.default  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:wKWVkKVCH46CnaG7 */
    1    60 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:i9FPhyqdbct5qPMc */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 cali-pro-_6syrmAMTuGHlzH1yqi  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:OsbgudgCLFtfeqSk */
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:4_lYosJ-d_bh1H__ */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:n6FWAtaMKXztkVyD */ /* Drop if no profiles matched */

Chain cali-pri-_6syrmAMTuGHlzH1yqi (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-pri-kns.default (2 references)
 pkts bytes target     prot opt in     out     source               destination
 239K   14M MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:7Fnh7Pv3_98FtLW7 */ MARK or 0x10000
 239K   14M RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:ZbV6bJXWSRefjK0u */ mark match 0x10000/0x10000

Chain cali-pri-ksa.default.default (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-pro-_6syrmAMTuGHlzH1yqi (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-pro-kns.default (2 references)
 pkts bytes target     prot opt in     out     source               destination
 115K 8175K MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:oLzzje5WExbgfib5 */ MARK or 0x10000
 115K 8175K RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:4goskqvxh5xcGw3s */ mark match 0x10000/0x10000

Chain cali-pro-ksa.default.default (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-to-hep-forward (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-to-host-endpoint (1 references)
 pkts bytes target     prot opt in     out     source               destination

Chain cali-to-wl-dispatch (1 references)
 pkts bytes target     prot opt in     out     source               destination
1793K  599M cali-tw-cali291dce4466b  all  --  *      cali291dce4466b  0.0.0.0/0            0.0.0.0/0           [goto]  /* cali:f_Ux9oyFvY4QbiNW */
 3160 2337K cali-tw-calie50f5cd5607  all  --  *      calie50f5cd5607  0.0.0.0/0            0.0.0.0/0           [goto]  /* cali:pBgpxRIQcmlJP3eY */
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:kLmvq4EvhFAVUA9c */ /* Unknown interface */

Chain cali-tw-cali291dce4466b (1 references)
 pkts bytes target     prot opt in     out     source               destination
  70M   23G ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:FfjEKGD81xi-95yK */ ctstate RELATED,ESTABLISHED
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:pav_SuunNxwldWZ5 */ ctstate INVALID
 239K   14M MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:Drcva8-MmGwB-Et3 */ MARK and 0xfffeffff
 239K   14M cali-pri-kns.default  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:7xCiskIsKHYe5_iP */
 239K   14M RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:2THp5JkE8Xtfzr4K */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 cali-pri-ksa.default.default  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:cknB6iVskY-FQLWY */
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:OMQjt-DWaF7k6LWF */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:vAuURCOx8cse5gur */ /* Drop if no profiles matched */

Chain cali-tw-calie50f5cd5607 (1 references)
 pkts bytes target     prot opt in     out     source               destination
 3153 2336K ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:mQiGM86s-mTudCN8 */ ctstate RELATED,ESTABLISHED
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:ACUmzOPQkr_zfwVH */ ctstate INVALID
    7   436 MARK       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:8FZbvxjvLQumdc7z */ MARK and 0xfffeffff
    7   436 cali-pri-kns.default  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:-Z-SQY0yWgZM7tGU */
    7   436 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:4PdtG-NHhTxAt-mC */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 cali-pri-_6syrmAMTuGHlzH1yqi  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:7uCzuhZCMVTMcZD6 */
    0     0 RETURN     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:31BSwgtPzlkNHd2p */ /* Return if profile accepted */ mark match 0x10000/0x10000
    0     0 DROP       all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:UONgDTeqHRULguOb */ /* Drop if no profiles matched */

Chain cali-wl-to-host (1 references)
 pkts bytes target     prot opt in     out     source               destination
 299K   25M cali-from-wl-dispatch  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:Ee9Sbo10IpVujdIY */
    0     0 ACCEPT     all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* cali:nSZbcOoG1xPONxb8 */ /* Configured DefaultEndpointToHostAction */
[root@dev-k8s-node4 ~]#
```
