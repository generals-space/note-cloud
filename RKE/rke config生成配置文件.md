rke_linux-amd64 config
[+] Cluster Level SSH Private Key Path [~/.ssh/id_rsa]:
[+] Number of Hosts [1]: 3
[+] SSH Address of host (1) [none]: 192.168.0.211
[+] SSH Port of host (1) [22]:
[+] SSH Private Key Path of host (192.168.0.211) [none]: ~/.ssh/id_rsa
[+] SSH User of host (192.168.0.211) [ubuntu]: root
[+] Is host (192.168.0.211) a Control Plane host (y/n)? [y]: y
[+] Is host (192.168.0.211) a Worker host (y/n)? [n]: n
[+] Is host (192.168.0.211) an etcd host (y/n)? [n]: y
[+] Override Hostname of host (192.168.0.211) [none]: rke-master-01
[+] Internal IP of host (192.168.0.211) [none]: 192.168.0.211
[+] Docker socket path on host (192.168.0.211) [/var/run/docker.sock]:
[+] SSH Address of host (2) [none]: 192.168.0.214
[+] SSH Port of host (2) [22]:
[+] SSH Private Key Path of host (192.168.0.214) [none]: ~/.ssh/id_rsa
[+] SSH User of host (192.168.0.214) [ubuntu]: root
[+] Is host (192.168.0.214) a Control Plane host (y/n)? [y]: n
[+] Is host (192.168.0.214) a Worker host (y/n)? [n]: y
[+] Is host (192.168.0.214) an etcd host (y/n)? [n]: n
[+] Override Hostname of host (192.168.0.214) [none]: rke-worker-01
[+] Internal IP of host (192.168.0.214) [none]: 192.168.0.214
[+] Docker socket path on host (192.168.0.214) [/var/run/docker.sock]:
[+] SSH Address of host (3) [none]: 192.168.0.215
[+] SSH Port of host (3) [22]:
[+] SSH Private Key Path of host (192.168.0.215) [none]: ~/.ssh/id_rsa
[+] SSH User of host (192.168.0.215) [ubuntu]: root
[+] Is host (192.168.0.215) a Control Plane host (y/n)? [y]: n
[+] Is host (192.168.0.215) a Worker host (y/n)? [n]: y
[+] Is host (192.168.0.215) an etcd host (y/n)? [n]: n
[+] Override Hostname of host (192.168.0.215) [none]: rke-worker-02
[+] Internal IP of host (192.168.0.215) [none]: 192.168.0.215
[+] Docker socket path on host (192.168.0.215) [/var/run/docker.sock]:
[+] Network Plugin Type (flannel, calico, weave, canal) [canal]: calico
[+] Authentication Strategy [x509]:
[+] Authorization Mode (rbac, none) [rbac]:
[+] Kubernetes Docker image [rancher/hyperkube:v1.16.3-rancher1]:
[+] Cluster domain [cluster.local]: k8s-server-lb
[+] Service Cluster IP Range [10.43.0.0/16]:
[+] Enable PodSecurityPolicy [n]:
[+] Cluster Network CIDR [10.42.0.0/16]:
[+] Cluster DNS Service IP [10.43.0.10]:
[+] Add addon manifest URLs or YAML files [no]:
