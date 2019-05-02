```
Step 1/11 : FROM hub.c.163.com/sunfmin/loraiotbase as builder
 ---> 94de96069609
Step 2/11 : WORKDIR ${GOPATH}/src/github.com/sunfmin/hupan3/backend/bms/anningAPP
 ---> Running in 20e497b103bb
Removing intermediate container 20e497b103bb
 ---> e85c5b9750f5
Step 3/11 : COPY . .
 ---> 6379e44aa71e
Step 4/11 : RUN cd app && go get -d -v
 ---> Running in c4e4747c751d
Removing intermediate container c4e4747c751d
 ---> cc798bb05479
Step 5/11 : RUN set -x && CGO_ENABLED=0 GOOS=linux go build -a -o /app ./app
 ---> Running in 38bd0ca83502
+ CGO_ENABLED=0
+ GOOS=linux
+ go build -a -o /app ./app
max depth exceeded
```