## 本质上说, --from-file 指定一个文件创建的 cm, 还是一个目录, 只不过目录下只有一个文件而已.
## 如果命令行中指定的是一个目录, 则会包含该目录下所有文件.
kubectl create configmap es-config --from-file=./cm/es/elasticsearch.yml
kubectl create configmap kibana-config --from-file=./cm/kibana/kibana.yml
kubectl create configmap logstash-config --from-file=./cm/logstash/logstash.yml
kubectl create configmap logstash-pipeline-config --from-file=./cm/logstash/pipeline
kubectl create configmap nginx-config --from-file=./cm/nginx/nginx.conf
