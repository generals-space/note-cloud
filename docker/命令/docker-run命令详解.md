`-e`/`--env`: 指定容器启动时的环境变量

```
$ docker run -it -e name=general -e age=12 centos6:1.0.0 /bin/bash
```

在容器内可以通过`$name`与`$age`读取到指定的环境变量.

docker run -d -p 8800:8822      [宿主机端口:容器端口]