# Provision the servers to support Kubernetes
The document is going to guide you through all the steps of provisioning the VMs to support Kubernetes from scratch.

1. Disable swap memory usage because Kubernetes doesn't support memory swap
```sh
$ swapoff -a
```
Note: you may need to comment out the line for memory swap in /etc/fstab
```sh
# e.g.
...
/dev/sdb         none            swap    sw 0    0
...
```

Reference:
[Why disable swap on kubernetes](https://serverfault.com/questions/881517/why-disable-swap-on-kubernetes)

2. Update /etc/hosts by adding the mapping of IPs and server names
```sh
# in my case, the following lines are added to /etc/hosts
100.100.1.241   kube-master
100.100.1.242   kube-worker-1
100.100.1.243   kube-worker-2
100.100.1.244   kube-worker-3
```

3. Install Docker in the servers
```sh
$ sudo apt update && \
  sudo apt install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common && \
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
  sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && \
  sudo apt update && \
  sudo apt install -y docker-ce

# check if Docker is running properly
$ sudo systemctl status docker # click q to exit
```

4. Install Kubernetes
```sh
# for both master and worker
$ curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
  echo 'deb http://apt.kubernetes.io/ kubernetes-xenial main' | tee /etc/apt/sources.list.d/kubernetes.list && \
  sudo apt update && \
  sudo apt install -y kubelet kubeadm kubectl

# for master only, it's recommended to spin up the Kubernetes master with 2 CPU cores; but in my case, the master only has one CPU core
$ kubeadm init  --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=200.200.1.241 --token-ttl=0 --ignore-preflight-errors=NumCPU

# move config file to home dir
$ sudo mkdir -p $HOME/.kube && \
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

# deploy the CNI to your cluster; in this tutorial, I deployed Calico and the other option is Flannel
$  kubectl apply -f https://docs.projectcalico.org/v3.7/manifests/calico.yaml

# or Flannel
$ kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

# check if all kube-system pods are successfully running (it takes couple of seconds to bring up all pods); once everything is up, move on to provision Kubernetes workers
$ kubectl get pods -n kube-system


# for worker only, once the initialization is done, note down the command for Kubernetes workers to join the cluster
...
kubeadm join 100.100.1.241:6443 --token <token-hash>  --discovery-token-ca-cert-hash <cert-hash> 
...

# note: If you get the error that token is expired while executing the command to join the cluster, you need to generate a new token.
$ kubeadm token create
```

5. Configuration for private Docker regiestry
```sh
# Create a new secret with provided information
$ kubectl create secret docker-registry <secret-key> --docker-server=quay.io --docker-username=<username> --docker-password=<password>
secret/<secret-key> created
...

# Then edit serviceaccounts
$ kubectl edit serviceaccounts default

Add

```sh
imagePullSecrets:
- name: <secret-key>
```
To the end after `Secrets`
```

References:
https://stackoverflow.com/questions/32726923/pulling-images-from-private-registry-in-kubernetes
