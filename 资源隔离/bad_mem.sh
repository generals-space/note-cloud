#!/bin/bash
x='a'
while [ True ]; do
    ## 变量x以2倍速增长
    x=$x$x
done;

## 如果开启了OOM(sysctl -w vm.panic_on_oom=1), 会被系统kill, 并在/var/log/message留下日志.
## 如果
## $ bash ./mem.sh
## ./mem.sh: xrealloc: cannot allocate 18446744071562068096 bytes (32768 bytes allocated)

## 20230421 更新, 可以考虑使用 memtester, 看起来会更直观, 可控性也更高.
