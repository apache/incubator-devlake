import axios from 'axios'
import { DEVLAKE_ENDPOINT } from './config'

export default async function handler(req, res) {
  const r = await axios.post(`${DEVLAKE_ENDPOINT}/task`, req.body)
  res.json(r.data)
}
