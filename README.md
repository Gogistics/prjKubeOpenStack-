# Deploy Kubernetes cluster on OpenStack
This tutorial aims to guide the developers how to deploy a Kubernetes cluster on OpenStack, the solution may also work for private clouds. One of advantages of deploying Kubernetes cluster from scratch because the developers will be able to seemlessly set up Kubernetes clusters on different platforms, especially private ones. In this tutorial, the Kubernetes cluster will be deployed in OpenStack inside the private cloud.

## The steps of deploying Kubernetes cluster are (say the OpenStack is running properly):
1. Create a network and a host server (this step is skipped because this tutorial is for learning deployment of Kubernets cluster from scratch)
2. Create the servers that will run Kubernets master and workers. (I suggest at least creating three servers, one is running as the master, and the other two are running as the workers.)
```sh
# In my case, a server can be instantiated by the following command and runs in Ubuntu 16.04.1 LTS. Once the server is instantiated, remember to note down the server IP
$ openstack server create --image <image-id> --flavor <flavor-id> --network <network-id> --key-name <key-name> --wait <server-host-name>
```

Reference:
[OpenStack CLI cheatsheet](https://docs.openstack.org/ocata/user-guide/cli-cheat-sheet.html)

3. Provision the servers
3-1. Disable swap memory usage because Kubernetes doesn't support memory swap
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

3-2. Update /etc/hosts by adding the mapping of IPs and server names
```sh
# in my case, the following lines are added to /etc/hosts
100.100.1.241   kube-master
100.100.1.242   kube-worker-1
100.100.1.243   kube-worker-2
100.100.1.244   kube-worker-3
```

3-3. Install Docker in the servers
```sh
$ sudo apt update && \
  sudo apt install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common && \
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
  sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && \
  sudo apt update && \
  sudo apt install -y docker-ce

# check if Docker is running properly
$ sudo systemctl status docker
```

3-4. Install Kubernetes
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

4. Development of applications in Golang (WIP)

