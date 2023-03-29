```yaml
version: "3"
```

在`docker-compose.yml`文件中声明`version`字段可以使用不同版本的字段, 各版本可能会有所出入.

通过`docker-compose`可以创建`services`, `networks`和`volumes`, 且ta们的生命周期可以通过`up`和`down`等命令控制.

[version 3](https://docs.docker.com/compose/compose-file/)

[version 2](https://docs.docker.com/compose/compose-file/compose-file-v2/)

[version 1](https://docs.docker.com/compose/compose-file/compose-file-v1/)

------

docker client: 19.03.4

docker server 20.10.12

可以用docker compose子命令, 代替docker-compose了, 后者比前者支持的语法版本更高, 比如deploy.replicas.

