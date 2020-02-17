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
