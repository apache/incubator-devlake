import { DEVLAKE_ENDPOINT } from '../../../utils/config'
import request from '../../../utils/request'

export default async function handler (req, res) {
  const r = await request.post(`${DEVLAKE_ENDPOINT}/task`, req.body)
  res.json(r.data)
}
