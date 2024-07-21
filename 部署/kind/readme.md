
d cp /etc/yum.repos.d/kubernetes.repo root-kube-node-01-1:/etc/yum.repos.d/
yum install -y kubeadm-1.17.2 kubelet-1.17.2 kubectl-1.17.2

kubeadm join kube-apiserver.generals.space:6443 --token abcdef.0123456789abcdef --discovery-token-ca-cert-hash sha256:b969766a8c9d9dfc3615dff8767c5bd6b8aa7930bdd699d6bb2213c434904c61 --ignore-preflight-errors=all

/usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --container-runtime=remote --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock
