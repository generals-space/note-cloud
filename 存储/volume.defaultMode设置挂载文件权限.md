# volume.defaultMode设置挂载文件权限

```yaml
    spec:
      containers:
        - name: kibana
          volumeMounts:
            - name: kibana-config-vol
              mountPath: /usr/share/kibana/config/kibana.yml
              subPath: kibana.yml
      volumes:
        - name: kibana-config-vol
          configMap:
            defaultMode: 600
            name: kibana-config
```
