k create cm prometheus-config -n monitoring --from-file=prometheus.yaml
k create cm prometheus-rules -n monitoring --from-file=./config/rules
k create cm blackbox-exporter-config -n monitoring --from-file=config.yml=./config/blackbox-exporter.yaml
