# ubuntu容器安装nsenter

ubuntu14的`util-linux`版本为2.20, 但想要进入docker容器, 不能低于`2.24`. 需要手动编译安装. 安装命令如下, 注意要首先安装依赖包.

```shell
sudo apt-get install autopoint autoconf libtool automake
wget https://www.kernel.org/pub/linux/utils/util-linux/v2.24/util-linux-2.24.tar.gz
tar xzvf util-linux-2.24.tar.gz
cd util-linux-2.24
./configure --without-ncurses
make && make install
```
