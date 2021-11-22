import { DEVLAKE_ENDPOINT, GRAFANA_ENDPOINT } from './config'
import request from '../../../utils/request'
export default async function handler (req, res) {
  const r = await request.get(`${DEVLAKE_ENDPOINT}/task?status=TASK_CREATED`)

  res.json({
    grafanaEndpoint: GRAFANA_ENDPOINT,
    tasks: r.data.tasks
  })
}
