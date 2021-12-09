FROM grafana/grafana:8.0.6
COPY ./provisioning/dashboards /etc/grafana/provisioning/dashboards
COPY ./provisioning/datasources /etc/grafana/provisioning/datasources
COPY ./dashboards /etc/grafana/dashboards
COPY ./img/grafana_icon.svg /usr/share/grafana/public/img/grafana_icon.svg
COPY ./img /usr/share/grafana/public/img/lake

# These three lines disable login pages
ENV GF_AUTH_ANONYMOUS_ENABLED=true
ENV GF_AUTH_ORG_ROLE=viewer
ENV GF_AUTH_DISABLE_LOGIN_FORM=true
