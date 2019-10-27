kubectl create ns harbor

kubectl create configmap harbor-db-env -n harbor --from-env-file=../config/db/env
kubectl create configmap registryctl-env -n harbor --from-env-file=../config/registryctl/env
kubectl create configmap core-env -n harbor --from-env-file=../config/core/env
kubectl create configmap jobservice-env -n harbor --from-env-file=../config/jobservice/env

kubectl create configmap registry-cfg-map -n harbor --from-file=../config/registry
kubectl create configmap registryctl-cfg-map -n harbor --from-file=../config/registryctl
kubectl create configmap core-cfg-map -n harbor --from-file=../config/core
kubectl create configmap jobservice-cfg-map -n harbor --from-file=../config/jobservice

kubectl create configmap proxy-crt-map -n harbor --from-file=../config/nginx
cp ../config/nginx/nginx.conf /mnt/nfsvol/harbor/proxy_nginx/
chown nfsnobody:nfsnobody /mnt/nfsvol/harbor/proxy_nginx/nginx.conf

## kubectl create configmap proxy-cfg-map -n harbor --from-file=../config/nginx/nginx.conf

## ingress

kubectl create secret tls https-certs -n harbor --cert=../config/nginx/server.crt --key=../config/nginx/server.key
