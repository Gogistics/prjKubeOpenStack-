Set up Kube master HA in OpenStack
Kubernetes High-Availability is about setting up Kubernetes, along with its supporting components in a way that there is no single point of failure. This document aims to guide the developers to deploy Kubernetes master HA cluster on OpenStack.

## Setup etcd cluster
1. Install cfssl and cfssljson
```sh
$ curl -o /usr/local/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 &&
  curl -o /usr/local/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64

$ chmod +x /usr/local/bin/cfssl* &&
  export PATH=$PATH:/usr/local/bin
```

2. Generate certificates on master-0

References:
https://medium.com/velotio-perspectives/demystifying-high-availability-in-kubernetes-using-kubeadm-3d83ed8c458b
https://coreos.com/os/docs/latest/generate-self-signed-certificates.html

```sh
$ mkdir -p /etc/kubernetes/pki/etcd &&
  cd /etc/kubernetes/pki/etcd
```

3. Create ca-config.json file in /etc/kubernetes/pki/etcd folder with following content:
```sh
{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "server": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth",
                    "client auth"
                ]
            },
            "client": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            },
            "peer": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth",
                    "client auth"
                ]
            }
        }
    }
}
```

4. Create ca-csr.json file in /etc/kubernetes/pki/etcd folder with following content.
```sh
{
    "CN": "etcd",
    "key": {
        "algo": "rsa",
        "size": 2048
    }
}
```

5. Create client.json file in /etc/kubernetes/pki/etcd folder with following content.
```sh
{
    "CN": "client",
    "key": {
        "algo": "ecdsa",
        "size": 256
    }
}

# cert_cmd
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca - &&
  cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client.json | cfssljson -bare client

# check if all certs are valid
$ openssl x509 -in ca.pem -text -noout
$ openssl x509 -in server.pem -text -noout
$ openssl x509 -in client.pem -text -noout
```

6. Create a directory /etc/kubernetes/pki/etcd on master-1 and master-2 and copy all the generated certificates into it.

7. On all masters, now generate peer and etcd certs in /etc/kubernetes/pki/etcd. To generate them, we need the previous CA certificates on all masters.
Note: the interface may be different; to check the interface run `ifconfig`
```sh
$ export PEER_NAME=$(hostname) &&
  export PRIVATE_IP=$(ip addr show ens3 | grep -Po 'inet \K[\d.]+')

$ cfssl print-defaults csr > config.json &&
  sed -i 's/www\.example\.net/'"$PRIVATE_IP"'/' config.json &&
  sed -i 's/example\.net/'"$PEER_NAME"'/' config.json &&
  sed -i '0,/CN/{s/example\.net/'"$PEER_NAME"'/}' config.json

$ cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server config.json | cfssljson -bare server &&
  cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=peer config.json | cfssljson -bare peer
```

8. On all masters, Install etcd and set it’s environment file
```sh
$ cd /usr/local/src

$ sudo wget "https://github.com/coreos/etcd/releases/download/v3.3.9/etcd-v3.3.9-linux-amd64.tar.gz" &&
  sudo tar -xvf etcd-v3.3.9-linux-amd64.tar.gz &&
  sudo mv etcd-v3.3.9-linux-amd64/etcd* /usr/local/bin/

$ sudo mkdir -p /etc/etcd /var/lib/etcd &&
  groupadd -f -g 15001 etcd &&
  useradd -c "etcd user" -d /var/lib/etcd -s /bin/false -g etcd -u 15001 etcd &&
  chown -R etcd:etcd /var/lib/etcd &&
  chown -R etcd:etcd /etc/kubernetes/pki/etcd
```

Reference:
https://devopscube.com/setup-etcd-cluster-linux/

9. Now, we will create a 3 node etcd cluster on all 3 master nodes. Starting etcd service on all three nodes as systemd. Create a file /etc/systemd/system/etcd.service on all masters

10. Ensure that you will replace the following placeholders with the proper values
Note: in this case, I use http as example because etcd always complains the error `x509 certificate signed by unknown authority`
```sh
[Unit]
Description=etcd service
Documentation=https://github.com/coreos/etcd

[Service]
User=etcd
Type=notify
ExecStart=/usr/local/bin/etcd \
 --name <host_name> \
 --data-dir /var/lib/etcd \
 --initial-advertise-peer-urls http://<host_private_ip>:2380 \
 --listen-peer-urls http://<host_private_ip>:2380 \
 --listen-client-urls http://<host_private_ip>:2379,http://127.0.0.1:2379 \
 --advertise-client-urls http://<host_private_ip>:2379 \
 --cert-file=/etc/kubernetes/pki/etcd/server.pem \
 --key-file=/etc/kubernetes/pki/etcd/server-key.pem \
 --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.pem \
 --peer-cert-file=/etc/kubernetes/pki/etcd/peer.pem \
 --peer-key-file=/etc/kubernetes/pki/etcd/peer-key.pem \
 --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.pem \
 --auto-tls=false \
 --peer-auto-tls=false \
 --initial-cluster-token etcd-cluster-1 \
 --initial-cluster kube-master-0=http://<master0_private_ip>:2380,kube-master-1=http://<master1_private_ip>:2380,kube-master-2=http://<master2_private_ip>:2380 \
 --initial-cluster-state new \
 --client-cert-auth=false \
 --peer-client-cert-auth=false \
 --heartbeat-interval 1000 \
 --election-timeout 5000
Restart=on-failure
RestartSec=5
LimitNOFILE=40000
TimeoutStartSec=2

[Install]
WantedBy=multi-user.target

<host_name> : Replace as the master’s hostname
<host_private_ip>: Replace as the current host private IP
<master0_private_ip>: Replace as the master-0 private IP
<master1_private_ip>: Replace as the master-1 private IP
<master2_private_ip>: Replace as the master-2 private IP
```

Reference:
https://github.com/etcd-io/etcd/blob/e205d09895e6e9d810a88923a64104474002c0c4/Documentation/op-guide/security.md#example-2-client-to-server-authentication-with-https-client-certificates

11. Start the etcd service on all three master nodes and check the etcd cluster health:
```sh
$ systemctl daemon-reload &&
  systemctl disable etcd &&
  systemctl enable etcd &&
  systemctl start etcd
```

## Setup load balancer

1. Setup a virtual IP load balancer to a healthy node master. keepalived is for load balancing, install keepalived on all master nodes.
```sh
$ apt install keepalived -y
```

2. Create the following configuration file /etc/keepalived/keepalived.conf on all master nodes:
```sh
#!/bin/sh

errorExit() {
    echo "*** $*" 1>&2
    exit 1
}

curl --silent --max-time 2 --insecure https://localhost:6443/ -o /dev/null || errorExit "Error GET https://localhost:6443/"
if ip addr | grep -q <VIRTUAL-IP>; then
    curl --silent --max-time 2 --insecure https://<VIRTUAL-IP>:6443/ -o /dev/null || errorExit "Error GET https://<VIRTUAL-IP>:6443/"
fi
```

3. Restart keepalive
```sh
$ systemctl restart keepalived
```

## Setup three kube master node cluster

1. Run kubeadm init on kube-master-0:
Create config.yaml file with following content.
```sh
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
api:
  advertiseAddress: <master-private-ip>
etcd:
  external:
    endpoints:
    - http://<kube-master-0-ip-address>:2379
    - http://<kube-master-1-ip-address>:2379
    - http://<kube-master-2-ip-address>:2379
    #caFile: /etc/kubernetes/pki/etcd/ca.pem
    #certFile: /etc/kubernetes/pki/etcd/client.pem
    #keyFile: /etc/kubernetes/pki/etcd/client-key.pem
networking:
  podSubnet: <pod-CIDR>
apiServerCertSANs:
- <load-balancer-ip>
apiServerExtraArgs:
  endpoint-reconciler-type: lease
```

2. Initialize kube master
```sh
# in my case, the VM has only core so I have to ignore NumCPU error
$ kubeadm init --config=config.yaml --ignore-preflight-errors=NumCPU
```

3. Move config file to home dir
$ sudo mkdir -p $HOME/.kube && \
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && \
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Note: remember to keep the result generated by `kubeadm init ...` somewhere for joing the the cluster at worker node.
```sh
# for example
$ kubeadm join 200.200.1.248:6443 --token sbd8yp.7dmf1ddx4xcmodma \
    --discovery-token-ca-cert-hash sha256:04ead9eada173eb55270fc41326c47d0216853791290f8f9cfa5a3d68ce57e60
```

4. Deploy the CNI to your cluster; in this tutorial, I deployed Calico and the other option is Flannel
$ kubectl apply -f https://docs.projectcalico.org/v3.7/manifests/calico.yaml


Result:
```sh
...
root@kube-master-0:~/atai# kubectl get pods -n kube-system
NAME                                       READY   STATUS    RESTARTS   AGE
calico-kube-controllers-6fb584dd97-tphfm   1/1     Running   0          10m
calico-node-4n5z4                          1/1     Running   0          10m
calico-node-bngbd                          1/1     Running   0          10m
calico-node-wldz5                          1/1     Running   0          10m
coredns-5c98db65d4-92grj                   1/1     Running   0          23m
coredns-5c98db65d4-rkk2d                   1/1     Running   1          23m
kube-apiserver-kube-master-0               1/1     Running   0          22m
kube-apiserver-kube-master-1               1/1     Running   0          14m
kube-apiserver-kube-master-2               1/1     Running   1          14m
kube-controller-manager-kube-master-0      1/1     Running   2          22m
kube-controller-manager-kube-master-1      1/1     Running   2          15m
kube-controller-manager-kube-master-2      1/1     Running   1          14m
kube-proxy-4x9kp                           1/1     Running   0          16m
kube-proxy-969ln                           1/1     Running   0          15m
kube-proxy-x9mfx                           1/1     Running   0          23m
kube-scheduler-kube-master-0               1/1     Running   2          22m
kube-scheduler-kube-master-1               1/1     Running   1          14m
kube-scheduler-kube-master-2               1/1     Running   1          14m
root@kube-master-0:~/atai# 
...
```
