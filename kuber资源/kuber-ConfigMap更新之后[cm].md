# kuber-ConfigMap更新之后[cm]

`ConfigMap`通过`volume`挂载入`Pod`, 然后更新`ConfigMap`中的信息后, `Pod`内部的`ConfigMap`是会同步变动的, 但是由于Pod内的进程没有重启, 所以大部分场景还是需要重启一下Pod才会生效.
