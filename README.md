# Documents - Commands to Setup Kubernetes Cluster
## 1. Download The Installer and Bootstrap Scripts

```bash wget https://raw.githubusercontent.com/zealvora/kcsa-repo/refs/heads/main/install_k8s_stack.sh```
 
 
```bash wget https://raw.githubusercontent.com/zealvora/kcsa-repo/refs/heads/main/bootstrap_cluster.sh```

## 2. Setup the Environment

```bash chmod +x install_k8s_stack.sh bootstrap_cluster.sh
./install_k8s_stack.sh
./bootstrap_cluster.sh```

## 3. Verify
```bash
kubectl get nodes
 
kubectl run nginx --image=nginx
 
kubectl delete pod nginx
 
docker ps
```
