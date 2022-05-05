FROM grafana/grafana:8.0.6
COPY ./provisioning/dashboards /etc/grafana/provisioning/dashboards
COPY ./provisioning/datasources /etc/grafana/provisioning/datasources
COPY ./dashboards /etc/grafana/dashboards
COPY ./img/grafana_icon.svg /usr/share/grafana/public/img/grafana_icon.svg
COPY ./img /usr/share/grafana/public/img/lake
RUN grafana-cli plugins install grafana-piechart-panel
