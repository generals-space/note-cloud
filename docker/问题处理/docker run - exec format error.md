# docker run - exec format error

参考文章

1. [“standard_init_linux.go:195: exec user process caused ”exec format error“” when run gitlab by docker](https://stackoverflow.com/questions/49765276/standard-init-linux-go195-exec-user-process-caused-exec-format-error-when)
2. [docker "exec format error"](https://blog.csdn.net/liduanwh/article/details/79999196)

在适配 arm 服务器上的 docker 镜像时, 在使用`docker run`启动一个 busybox-arm 的镜像, 报了上述错误. 最初以为是因为镜像中没有`bash`, 于是换了`sh`, 结果还是报这个错误.

参考文章1给出了答案, 是因为我把 arm 平台的镜像拉取下来, 在 x86 的服务器运行的, 所以才会出错, 如果在 arm 服务器上运行就没问题了.

------

还有一种情况, 就是启动脚本`init.sh`里我没写`#!/bin/sh`, 但是 docker/kuber 的 command 里写了 `/bin/bash -c init.sh`, 导致解释器不匹配. 参考文章2中给出了解释.
