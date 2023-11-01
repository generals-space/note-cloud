# 为Pod配置hosts[域名]

参考文章

1. [k8s给pod添加hosts](https://www.cnblogs.com/route/p/16119323.html)

注意: 这里配置 hosts 只是指的 /etc/hosts 配置文件, 并不是指定 hostname 主机名.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
## ...省略
spec:
  template:
    metadata:
    ## ...省略
    spec:
      hostAliases:
      - ip: "10.236.9.220"
        hostnames:
        - "xxx.xxx.local"
      containers:
      - name: 
        image: 

```
