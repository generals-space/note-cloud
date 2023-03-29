```yaml
version: '3'
services:
  mysql:
    image: mysql:5.7
    restart: always
    container_name: mysql
    networks:
    - mydb
    volumes:
    - /opt/mydb/mysql:/var/lib/mysql
    environment:
    - "MYSQL_ROOT_PASSWORD=5u)nu8ssM-sjeexD"
    - "MYSQL_DATABASE=mydb"
    - "TZ=Asia/Shanghai"
    - "LANG=C.UTF-8"

  phpmyadmin:
    image: phpmyadmin:5.1
    restart: always
    container_name: phpmyadmin
    networks:
    - mydb
    ports:
    - 8081:80
    environment:
    - "TZ=Asia/Shanghai"
    - "LANG=C.UTF-8"
    - "PMA_HOST=mysql"
    depends_on:
    - mysql

networks:
  mydb:
    driver: bridge
```
