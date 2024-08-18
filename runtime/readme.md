参考文章

1. [lxcfs容器隔离技术实现原理分析之loadavg、cpuonline](https://blog.csdn.net/ZVAyIVqt0UFji/article/details/103193083)
2. [使用 Go 和 Linux Kernel 技术探究容器化原理](https://zhuanlan.zhihu.com/p/512715825)
3. [nestybox/sysbox](https://github.com/nestybox/sysbox)
    - 与 runc、gVisor 同级的 runtime, 可以被 docker/kubernetes 使用.

1. [Docker背后的内核知识（一）](https://www.cnblogs.com/beiluowuzheng/p/10004132.html)
    - namespace的API包括`clone()`, `setns()`以及`unshare()`, 还有`/proc`下的部分文件.
    - 很多示例代码
2. [Docker背后的内核知识（二）](https://www.cnblogs.com/beiluowuzheng/p/10015177.html)
    - 主要是cgroup相关的知识
3. [Docker安全](https://www.bookstack.cn/read/dockerdocs/Articles-security.md)
    - 类似于引言一样的文章, 不涉及技术细节.
4. [docker基础知识之mount namespace](http://kuring.me/post/namespace_mount/)
    - c语言代码
5. [linux和docker的capabilities介绍](https://www.cnblogs.com/charlieroro/p/10108577.html)
