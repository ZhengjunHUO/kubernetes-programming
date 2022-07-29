## HOWTO
```bash
kubectl create ns controller
kubectl apply -f kubernetes/crd.yaml
kubectl apply -f kubernetes/rbac.yaml
kubectl apply -f kubernetes/deployment.yaml
# observe the leader election
kubectl logs -f -l app=hzjoperator
# delete the leader to see a relection
kubectl delete pod hzjoperator-6bdf4598c4-9qbbn
# add the cr 
kubectl apply -f kubernetes/fufu.yaml
# observe the cr info in operator
kubectl logs -f -l app=hzjoperator
```
