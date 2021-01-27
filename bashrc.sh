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
alias ipp='ip -6'
alias d='docker'
alias dk='docker'
alias dc='docker-compose'
alias k='kubectl'
alias kde='kubectl describe'
alias kap='kubectl apply -f'
alias kgl='kubectl get --show-labels'

export kcr='controllerrevision'

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
## @function:   进入容器内部, 优先尝试 bash 终端
## @note:       alias kex='kubectl exec -it'
## $1:          目标 Pod 名称(必选)
## $2:          目标资源的所属空间, 如 -n namespace(可选)
function kex() {
    ## 如果把标准输出也 null 掉, 当目标 Pod 有正常 bash 时就无法进行交互, 终端会卡死.
    ## kubectl exec -it $1 bash 1>/dev/null 2>/dev/null
    kubectl exec -it $@ bash 2>/dev/null
    if (( $? != 0 )); then
        echo 'bash does exist, try to use sh'
        kubectl exec -it $@ sh
    fi
}
## @function:   打印目标资源的 yaml 信息
## $1:          目标资源类型(必选)
## $2:          目标资源名称(必选)
## $3:          目标资源的所属空间, 如 -n namespace(可选)
function kya() {
    ## $@ 表示传入的所有参数.
    kubectl get -o yaml $@
}
## @function:   打印目标资源的详细信息
## $1:          目标资源类型(必选)
## $2:          目标资源名称(必选)
## $3:          目标资源的所属空间, 如 -n namespace(可选)
function kwd() {
    ## $@ 表示传入的所有参数.
    kubectl get -o wide $@
}
## @function:   强制删除一个pod
## $1:          目标Pod名称(必选)
## $2:          目标资源的所属空间, 如 -n namespace(可选)
function kkill() {
    ## $@ 表示传入的所有参数.
    kubectl delete pod --force --grace-period 0 $@
}
###################################################################
## @function: 清空目标容器的日志(docker clear log)
## $1:        目标容器名称或id
function dclog() {
    local log_path=$(docker inspect -f '{{.LogPath}}' $1)
    ## 清空目标文件
    :>$log_path
}
## @function: 进入目标容器网络空间
## $1:        目标容器名称/ID
function denter(){
    local dpid=$(docker inspect -f "{{.State.Pid}}" $1)
    nsenter -t $dpid --net /bin/sh
}
## @function: 进入目标容器网络空间
## $1:        目标容器名称/ID
function dnsenter() {
    local dpid=$(docker inspect -f "{{.State.Pid}}" $1)
    nsenter -t $dpid --net /bin/sh 
}
## @function: 清理不用的容器和镜像
function dclean() {
    ## 删除旧容器, 不会删除当前正在运行的容器
    docker rm $(docker ps -a -q)
    ## 移除本地多余的历史镜像
    docker images | grep '<none>' | awk '{print $3}' | xargs docker rmi
}
## @function:   进入容器内部, 优先尝试 bash 终端
## @note:       alias dex='docker exec -it'
## $1:          目标 Pod 名称
function dex() {
    ## 如果把标准输出也 null 掉, 当目标 Pod 有正常 bash 时就无法进行交互, 终端会卡死.
    ## docker exec -it $1 bash 1>/dev/null 2>/dev/null
    docker exec -it $1 bash 2>/dev/null
    if (( $? != 0 )); then
        echo 'bash does exist, try to use sh'
        docker exec -it $1 sh
    fi
}
