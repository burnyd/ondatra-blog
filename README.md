## Create the KIND cluster with the cluster.yaml file
kind create cluster --config cluster.yaml
## Create meshnet daemonset.
kubectl apply -k /home/dhertzberg/projects/meshnet-cni/manifests/base
## Apply ceos operator
https://github.com/aristanetworks/arista-ceoslab-operator
## Apply metalblb
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml
## IPaddress pool for kind..
kubectl apply -f metallb-config.yaml
## Move ceos images
ceoslab                                  4.31.0F       ecc52bdfdf25   7 weeks ago     2.42GB
kind load docker-image  ceoslab:4.31.0F
## Apply KNE Binary
git pull kne && cd kne_cli && go build.. then move file here
./kne_cli create ceos.pb.txt
```
kubectl get pods -n ceos
NAME              READY   STATUS    RESTARTS   AGE
keng-controller   2/2     Running   0          26m
keng-port-eth1    2/2     Running   0          26m
keng-port-eth2    2/2     Running   0          26m
r1                1/1     Running   0          26m
r2                1/1     Running   0          26m
```

## Run the unit test.
### Change the topology key first to where this repo is.  Mine looks like the following.
### topology: /home/dhertzberg/projects/ondatrablog/ceos.pb.txt
### This will not take env vars like $PWD.
go test -testbed=testbed.textproto -config=$PWD/config.yaml

Testing output can be found with the [testingoutput.txt](testsoutput.txt) file.