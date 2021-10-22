import axios from 'axios'
import { DEVLAKE_ENDPOINT, GRAFANA_PORT } from './config'
import request from '../../../utils/request'
export default async function handler(req, res) {
  const r = await request.get(`${DEVLAKE_ENDPOINT}/task?status=TASK_CREATED`)

  res.json({
    grafanaPort: GRAFANA_PORT,
    tasks: r.data.tasks
  })
}
