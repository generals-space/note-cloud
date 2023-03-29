# docker-commit时容器仍在运行带来的麻烦

参考文章

1. [保存对容器的修改](http://www.docker.org.cn/book/docker/docer-save-changes-10.html)

真正有帮助的是一名游客的评论:

> 没写清楚`docker commit`前需要停止容器，我安装完ping之后试了下ping www.baidu.com 然后 commit，结果就是查看镜像信息时里面COMMAND有“ping www.baidu.com”, 造成每次以这个镜像启动一个容器就会去ping百度了