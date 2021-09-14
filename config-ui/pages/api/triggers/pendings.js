import axios from 'axios'
import { DEVLAKE_ENDPOINT, GRAFANA_PORT } from './config'

export default function handler(req, res) {

  axios.get(`${DEVLAKE_ENDPOINT}/task?status=TASK_CREATED`).then(r => {
      res.status(200).json({
          grafanaPort: GRAFANA_PORT,
          tasks: r.data.tasks
      })
    console.log(res.data)
  }).catch(e => {
    console.log(e)
  })
}
