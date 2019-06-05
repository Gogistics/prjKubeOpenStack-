

# Disable swap memory usage
```
$ swapoff -a
```


## Docker installation
```
$ sudo apt update
$ sudo apt install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
$ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
$ sudo apt update

# Install Docker
$ sudo apt install -y docker-ce

# Check if Docker is running
$ sudo systemctl status docker
```

## Kubernetes installation and configuration
```
# Installation:
$ curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
$ echo 'deb http://apt.kubernetes.io/ kubernetes-xenial main' | tee /etc/apt/sources.list.d/kubernetes.list
$ sudo apt update
$ sudo apt install -y kubelet kubeadm kubectl

# Configuration:
# Configuration for the master-
$ kubeadm init  --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=200.200.1.241 --token-ttl=0 --ignore-preflight-errors=NumCPU
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Deploy the CNI to your cluster
# Calico
$ kubectl apply -f https://docs.projectcalico.org/v3.7/manifests/calico.yaml

Or

# Flannel
$ kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
# \Flannel
# \Configuration for the master-


# Check if all kube pods are successfully running (it takes couple of seconds to bring up all pods)
$ kubectl get pods -n kube-system
```



# Configuration for private Docker regiestry
```
# Create a new secret with provided information
$ kubectl create secret docker-registry quaySCM --docker-server=quay.io --docker-username="ocedo+jenkins" --docker-password="E4OB2QWD9A8ZE52X05M7SNXYTOOTVCXNBH2V1MQS4GXPLBI24T79O4WEHE8AQAS7"
secret/quay-scm created
...

# Then edit serviceaccounts
$ kubectl edit serviceaccounts default
```

References:
https://stackoverflow.com/questions/32726923/pulling-images-from-private-registry-in-kubernetes
