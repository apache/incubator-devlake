# Deploy devlake with helm

## Prerequest

- Helm >= 3.6.0
- Kubernetes >= 1.19.0

## Quick Install

clone the code, and enter the deployment/helm folder.
```
helm install devlake .
```

And visit your devlake from the node port (32001 by default).

## Parameters

Some useful parameter for the chart, you could also check them in values.yaml

| Parameter | Description | Default |
|-----------|-------------|---------|
| replicaCount  | Replica Count for devlake, currently not used  | 1  |
| mysql.useExternal  | If use external mysql server, currently not used  |  false  |
| mysql.externalServer  | External mysql server address  | 127.0.0.1  |
| mysql.externalPort  | External mysql server port  | 3306  |
| mysql.username  | username for mysql | merico  |
| mysql.password  | password for mysql | merico  |
| mysql.database  | database for mysql | lake  |
| mysql.rootPassword  | root password for mysql | admin  |
| mysql.storage.class  | storage class for mysql's volume | ""  |
| mysql.storage.size  | volume size for mysql's data | 5Gi  |
| mysql.image.repository  | repository for mysql's image | mysql  |
| mysql.image.tag  | image tag for mysql's image | 8.0.26  |
| mysql.image.pullPolicy  | pullPolicy for mysql's image | IfNotPresent  |
| grafana.image.repository  | repository for grafana's image | mericodev/grafana  |
| grafana.image.tag  | image tag for grafana's image | latest  |
| grafana.image.pullPolicy  | pullPolicy for grafana's image | Always  |
| lake.storage.class  | storage class for lake's volume | ""  |
| lake.storage.size  | volume size for lake's data | 100Mi  |
| lake.image.repository  | repository for lake's image | mericodev/lake  |
| lake.image.tag  | image tag for lake's image | latest  |
| lake.image.pullPolicy  | pullPolicy for lake's image | Always  |
| ui.image.repository  | repository for ui's image | mericodev/config-ui  |
| ui.image.tag  | image tag for ui's image | latest  |
| ui.image.pullPolicy  | pullPolicy for ui's image | Always  |
| service.type  | Service type for exposed service | NodePort  |
| service.grafanaPort  | Service port for grafana | 32000  |
| service.uiPort  | Service port for config ui | 32001  |
| service.grafanaEndpoint  | The external grafana endpoint, used when ingress not configured  |  http://127.0.0.1:32000  |