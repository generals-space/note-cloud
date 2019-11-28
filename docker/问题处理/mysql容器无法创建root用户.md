# mysql容器无法创建root用户

参考文章

1. [ERROR 1396 (HY000): Operation CREATE USER failed for 'root'@'%'](https://github.com/docker-library/mysql/issues/129)

```yaml
      containers:
      - name: mysql
        image: daocloud.io/library/mysql:5.7.20
        env:
        - name: MYSQL_PASSWORD
          value: "123456"
        - name: MYSQL_ROOT_PASSWORD
          value: "123456"
        - name: MYSQL_USER
          value: root
```

> ERROR 1396 (HY000) at line 1: Operation CREATE USER failed for 'root'@'%'

尝试过把`MYSQL_USER`的值指定成其他名称, 但应用程序上是不行的, 因为程序需要root用户创建数据库, 而普通用户没有创建db的权限, 问题还是在无法创建root用户上面.

按照参考文章1所说, root用户是默认就存在的, 不需要用`MYSQL_USER`字段再添加一遍.

于是把`MYSQL_USER`字段删掉, 创建就成功了.
