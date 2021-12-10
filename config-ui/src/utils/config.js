export const DEVLAKE_ENDPOINT = '/api'
export const GRAFANA_BASE_URL = '/grafana'
export const LOCAL_BASE_URL = 'http://localhost:3002'
export const GRAFANA_ENDPOINT = '/d/0Rjxknc7z/demo-homepage?orgId=1'
export const GRAFANA_URL = process.env.LOCAL ? LOCAL_BASE_URL + GRAFANA_ENDPOINT : GRAFANA_BASE_URL + GRAFANA_ENDPOINT
