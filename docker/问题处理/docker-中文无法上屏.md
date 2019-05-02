# docker-中文无法上屏

参考文章

1. [Docker下不支持中文解决办法-已验证有效](https://www.jianshu.com/p/ecf13d88534b)

mysql容器, 运行mysql客户端时, 无法输入中文, 粘贴时的中文也被过滤掉了.

启动命令

```
$ docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=mydb -e MYSQL_USER=mydb -e MYSQL_PASSWORD=123456 -p 3306:3306 mysql
```

连接方式

```
$ docker exec -it mysql mysql -u ptcms -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
...
mysql>
```

复制如下命令时, 中文被过滤.

```sql
##下一行是原文：
mysql> insert into user(id,username) values('1','张三');
##变成了：
mysql> insert into user(id,username) values('1','');
```

参考文章1中详细讲解了解决方法, 不过最简单的方法是在容器启动时传入`LANG`环境变量.

```
docker run -d --name mysql -e LANG=C.UTF-8 -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=mydb -e MYSQL_USER=mydb -e MYSQL_PASSWORD=123456 -p 3306:3306 mysql
```

然后就可以了.