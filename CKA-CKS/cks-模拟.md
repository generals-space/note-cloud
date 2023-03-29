## 2. falco

## 3. 

kube-apiserver 可以通过参数`--kubernetes-service-node-port=31000`将原本`default`命名空间的`kubernetes` Service, 修改为 NodePort 类型.

## 4. pod security standard

## 10.

gvisor dmesg日志

## 12. secret 

这个排查思路跟我自己的差不多, 不过 secret3 还是没搞出来.

## 14. sysdig

呃, 这个不是用 sysdig 解决的, 用的是 strace.

## 18. p.auster ServiceAccount

没在集群中查到 p.auster 的 SA, 但本题其实不需要找到这个对象.

只要在 audit log 中找到 p.auster 访问的 secret 资源, 修改其中的 password 字段即可.

## 19. filesystem readonly

## 22. mysql dockerfile?
