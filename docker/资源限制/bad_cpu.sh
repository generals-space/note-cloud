#!/bin/bash
x=1
while [ True ]; do
    ## 变量x以2倍速增长
    x=$((x+1))
done;

## 该脚本会将CPU占用到100%, 但不会被kill, 而且由于是单线程, 多核主机上并不会造成太大影响, top输出为
##    PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
##  83019 general   20   0  113184   1192   1016 R 100.0  0.0   0:44.23 bash
