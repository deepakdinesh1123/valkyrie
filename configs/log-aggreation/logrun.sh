podman run -d \
  --name fluent-bit \
  --network podman \
  -v ../fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf \
  -v ../../internal/logs:/etc/data \
  -e WORKER_ID=$WORKER_ID \
  -e LOKI_HOST=$LOKI_HOST \
  fluent/fluent-bit

podman run -d \
  --name loki \
  --network podman \
  -p 3100:3100 \
  docker.io/grafana/loki:2.7.0 \
  -config.file=/etc/loki/local-config.yaml

podman run -d \
  --name grafana \
  --network podman \
  -p 3000:3000 \
  -e GF_PATHS_PROVISIONING=/etc/grafana/provisioning \
  -e GF_AUTH_ANONYMOUS_ENABLED=true \
  -e GF_AUTH_ANONYMOUS_ORG_ROLE=Admin \
  docker.io/grafana/grafana:latest \
  sh -euc "
        mkdir -p /etc/grafana/provisioning/datasources &&
        cat <<EOF > /etc/grafana/provisioning/datasources/ds.yaml
        apiVersion: 1
        datasources:
        - name: Loki
          type: loki
          access: proxy
          orgId: 1
          url: http://loki:3100
          basicAuth: false
          isDefault: true
          version: 1
          editable: true
        EOF
        /run.sh
  "

podman run -d \
  --name loki \
  --network loki \
  -p 3100:3100 \
  docker.io/grafana/loki:2.7.0 \
  -config.file=/etc/loki/local-config.yaml
