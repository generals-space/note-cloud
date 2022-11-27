kube-ipvs 设备不属于 CNI 范畴, ta实现的是 Service 到 Pod 的映射, 或者说, 是Service端口到Pod端口的映射(因为Service与Pod之间只通过端口而不是IP进行连接).
