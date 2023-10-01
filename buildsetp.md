docker build -t lakhlaifi/configprobe-operator:latest .
docker push configprobe-operator:latest


kubectl apply -f config/crd/bases/
kubectl apply -f config/manager/manager.yaml

