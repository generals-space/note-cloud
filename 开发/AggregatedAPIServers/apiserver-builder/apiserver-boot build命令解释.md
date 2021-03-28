# apiserver-boot build命令解释

参考文章

1. [使用Aggregated APIServer的方式构建API服务](https://jeremyxu2010.github.io/2019/07/%E4%BD%BF%E7%94%A8aggregated-apiserver%E7%9A%84%E6%96%B9%E5%BC%8F%E6%9E%84%E5%BB%BAapi%E6%9C%8D%E5%8A%A1/)
    - 全面的使用手册, 构建流程

```
apiserver-boot build executables
apiserver-boot build container --image myapiserver:0.0.1
apiserver-boot build config --name podgroup-apiserver --namespace default --image myapiserver:0.0.1
```

## `apiserver-boot build executables`

执行这条命令时, 将会先执行`apiserver-boot build generated`生成代码, 然后编译`cmd/apiserver`与`cmd/manager`, 在`bin`目录下生成`apiserver`和`controller-manager`两个可执行文件, 基本等同于如下

```bash
apiserver-boot build generated
CGO_ENABLED=0
go build -o bin/apiserver cmd/apiserver/main.go
go build -o bin/controller-manager cmd/manager/main.go
```

## `apiserver-boot build container --image myapiserver:0.0.1`

同样先`generated`, 然后编译`apiserver`和`controller-manager`, 然后将两个二进制文件封装到 docker 镜像.

```
apiserver-boot build generated

CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
go build -o /tmp/apiserver-boot-build-container073667538/apiserver cmd/apiserver/main.go
go build -o /tmp/apiserver-boot-build-container073667538/controller-manager cmd/manager/main.go

docker build -t myapiserver:0.0.1 /tmp/apiserver-boot-build-container073667538 
```

`/tmp/apiserver-boot-build-container073667538`目录下存在`Dockerfile`, 其内容如下

```dockerfile
FROM ubuntu:14.04

RUN apt-get update
RUN apt-get install -y ca-certificates

ADD apiserver .

ADD controller-manager .
```

## `apiserver-boot build config --name podgroup-apiserver --namespace default --image myapiserver:0.0.1`

在工程根目录下创建`config`子目录, 结构如下:

```
[root@k8s-master-01 config]# tree
.
├── apiserver.yaml
└── certificates
    ├── apiserver_ca.crt
    ├── apiserver_ca.key
    ├── apiserver_ca.srl
    ├── apiserver.crt
    ├── apiserver.csr
    └── apiserver.key

1 directory, 7 files
```

其中证书是调用`openssl`工具生成的, `apiserver.yaml`的内容则包含了很多东西, 当前目录有名为`apiserver.yaml`的文件可以查看.
