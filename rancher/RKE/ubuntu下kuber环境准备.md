使用rke也是需要预先安装依赖的, 包括防火墙, 内核模块, docker.

kubectl的源还要参考kubernetes的官方文档, 从kubeadm文档中挑出需要的步骤好了. 但是官方的源太慢了, 这个其实也可以从阿里云的kubernetes源中找到.

ubuntu下docker的安装步骤可见阿里云的镜像站. 然后不要忘了给ubuntu用户赋予docker的执行权限.

```
usermod -aG docker ubuntu
```
