# Deploy Kubernetes cluster on OpenStack
This tutorial aims to guide the developers how to deploy a Kubernetes cluster on OpenStack, the solution may also work for private clouds. One of advantages of deploying Kubernetes cluster from scratch because the developers will be able to seemlessly set up Kubernetes clusters on different platforms, especially private ones. In this tutorial, the Kubernetes cluster will be deployed in OpenStack inside the private cloud.

References:
[Basic and DSR load balancing with network load balancers (NLB)](https://cloud.ibm.com/docs/containers?topic=containers-loadbalancer)

## The steps of deploying Kubernetes cluster are (say the OpenStack is running properly):
1. Create a network and a host server (this step is skipped because this tutorial is for learning deployment of Kubernets cluster from scratch)
2. Create the servers that will run Kubernets master and workers. (I suggest at least creating three servers, one is running as the master, and the other two are running as the workers.)
```sh
# In my case, a server can be instantiated by the following command and runs in Ubuntu 16.04.1 LTS. Once the server is instantiated, remember to note down the server IP
$ openstack server create --image <image-id> --flavor <flavor-id> --network <network-id> --key-name <key-name> --wait <server-host-name>
```

Reference:
[OpenStack CLI cheatsheet](https://docs.openstack.org/ocata/user-guide/cli-cheat-sheet.html)

3. Provision the servers to support Kubernetes
Please refer to provision/guide.md

4. Deployment of Redis cluster
4-1. Write yaml files
4-2. Kubectl create/apply the yaml files

Reference:
[Deploying PHP Guestbook application with Redis](https://kubernetes.io/docs/tutorials/stateless-application/guestbook/#scale-the-web-frontend)

5. Development of applications in Golang (WIP)
In this tutorial, I wrote the applications in Golang. Feel free to write the applications in your favirote languages


Build the Golang application
```sh
$ docker build -t alantai/web-state:v0.1.1 -f dockerfiles/Dockerfile.apis_state .
```

Tag the Docker image
```sh
$ docker tag alantai/web-state:v0.1.1 quay.io/ocedo/scm-ui:web-state-v0.1.1
```

Push to Docker registry
```sh
$ docker push quay.io/ocedo/scm-ui:web-state-v0.1.1
```

6. Deploy the Golang application and the service
6-1. go to the kube-master console and update the image tag in the yaml file
6-2. make sure the LoadBalancer IP is exported and run the following command
```sh
$ envsubst < kube-apis-state-quay.yaml | kubectl apply -f -
```

References:
[dumb-init](https://github.com/Yelp/dumb-init)
[Introducing dumb-init, an init system for Docker containers](https://engineeringblog.yelp.com/2016/01/dumb-init-an-init-for-docker.html)

6. Manage the Kubernetes cluster
Here I just wrote some simple examples. Feel free to try otehr Kubernetes mechanisms.

```sh
# scale deployment
$ kubectl scale deployment apis-state --replicas=6

# rollback deployment
$ kubectl rollout status deployment/<deployment-name>
$ kubectl rollout history deployment/<deployment-name>
$ kubectl rollout undo deployment/<deployment-name>
```
