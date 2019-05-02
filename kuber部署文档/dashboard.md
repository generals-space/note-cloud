dashboard官网地址 https://github.com/kubernetes/dashboard

kubectl -s http://172.32.100.90:8080 proxy --address='0.0.0.0'


参考文章



[kubernetes-dashboard环境搭建](http://blog.csdn.net/weixin_38011359/article/details/68065986)

    - 这篇文章里讲到可以在部署完dashboard后, 通过service在apiserver主机上的端口映射进行访问, 但1.7版本的kube没有这个映射.

insecure模式下通过proxy转发以访问dashboard, 不然没有办法在未认证的情况下连接.

```
$ kubectl -s http://172.32.100.71:8080 proxy --address='0.0.0.0' --port=8001 --accept-hosts='^*$'
Starting to serve on [::]:8001
```

浏览器访问`http://172.32.100.71:8001/ui`完成.