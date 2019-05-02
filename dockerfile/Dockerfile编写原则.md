# Dockerfile编写原则

参考文章

1. [如何写好Dockerfile，Dockerfile最佳实践]()

2. [Dockerfile 最佳实践](https://my.oschina.net/u/2612999/blog/1036388)

## `--no-cache`的使用.

我们知道, Docker在构建镜像的过程就是顺序执行 Dockerfile 每个指令. 执行过程中，Docker将在缓存中查找可重用的镜像.

如果dockerfile内部引用的资源本身变化, 但路径(或者版本)没变, build的时候依然不会获取新资源.

比如, dockerfile中有如下语句

```
RUN curl http://xxx.xxx.xx.xx/project.tar.gz
```

每次项目更新都只要把新版本的工程包打成`project.tar.gz`就可以了. 当然, 你也可以将不同版本的工程包带上版本号, 但dockerfile不能像shell脚本那样传入参数, 这样就需要你每次更新工程也连带着更新dockerfile文件.

当然, 目的是好的, 实际上, 由于dockerfile没有变化, 而原来构建的镜像层也还在, 所以docker会认为可以直接使用已经缓存的镜像层替代, 这就很尴尬了...

所以, 如果你不想使用缓存，就需要`--no-cache=true`选项了.

```
$ docker build --no-cache=true -f 你的Docker路径 -t 镜像名称及标签 Dockerfile所在目录
```