helm-dev *args:
    helm install -f deployments/helm/valkyrie/values.dev.yaml valkyrie deployments/helm/valkyrie {{ args }}

helm-dev-purge *args:
    helm uninstall valkyrie
