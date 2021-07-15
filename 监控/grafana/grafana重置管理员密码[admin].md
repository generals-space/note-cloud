# grafana重置管理员密码[admin]

参考文章

1. [grafana忘记登陆密码](https://www.cnblogs.com/yexiuer/p/11287994.html)
    - sqlite3修改数据库
2. [grafana忘记admin密码，重置](https://www.cnblogs.com/ccielife/p/12802670.html)
    - `grafana-cli admin reset-admin-password xxxxx`

参考文章1需要找到grafana的`.db`文件, 但是我的环境是docker容器, 里面没有`sqlite3`命令, 而且用的是共享存储, 在外面改后再放回去有点风险(主要是不清楚怎么放回去).

然后找到了参考文章2, 内置命令`grafana-cli`真的香.

grafana版本: v6.4.3

