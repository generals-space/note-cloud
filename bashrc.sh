export ETCDCTL_API=3
export ETCDCTL_ENDPOINTS=https://192.168.0.101:2379
export ETCDCTL_CACERT=/etc/kubernetes/pki/etcd/ca.crt
export ETCDCTL_CERT=/etc/kubernetes/pki/etcd/server.crt
export ETCDCTL_KEY=/etc/kubernetes/pki/etcd/server.key

export GO111MODULE=on
export GOPROXY=https://goproxy.cn
## export GOPROXY=direct

function enproxy()
{
    export http_proxy=http://192.168.0.8:1081
    export https_proxy=http://192.168.0.8:1081
}
function deproxy()
{
    unset http_proxy
    unset https_proxy
}

alias c='curl'
alias d='docker'
alias dk='docker'
alias dc='docker-compose'
alias k='kubectl'
alias kde='kubectl describe'
alias kap='kubectl apply -f'

## @function: 获取当前的context列表
## @note:     kc -> kubectl context
function kc() {
    kubectl config get-contexts
}
## @function: 切换当前的context为目标$1
## @note:     ks -> kubectl switch
function ks() {
    kubectl config use-context $1
}
## @function: 查看目标容器的进程pid, 以便之后进入其网络空间
## $1:        容器名称/ID
function dpid(){
    docker inspect -f "{{.State.Pid}}" $1
}
## @function: 清空目标容器的日志
## $1:        目标容器名称或id
function dclear() {
    local log_path=$(docker inspect -f '{{.LogPath}}' $1)
    ## 清空目标文件
    :>$log_path
}
## @function: 进入目标容器
## $1:        目标容器名称或id
function denter(){
    docker exec -it $1 /bin/bash
}
## @function: 进入目标容器网络空间
## $1:        目标容器进程pid
function dnsenter() {
    nsenter -t $1 -n /bin/sh
}
