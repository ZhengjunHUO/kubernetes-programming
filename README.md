# kubernetes-programming

## operator
Rewrite [k8s-custom-controller](https://github.com/ZhengjunHUO/k8s-custom-controller) in a better way
```bash
# Build the operator image
docker build -t hzjoperator:v2 .
# If using kind, load the image into cluster
kind load docker-image hzjoperator:v2 --name huo
```
