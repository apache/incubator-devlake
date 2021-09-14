import axios from 'axios'
import { DEVLAKE_ENDPOINT } from './config'

export default function handler(req, res) {

  axios.post(`${DEVLAKE_ENDPOINT}/task`, req.body ).then(res => {
    console.log(res.data)
  }).catch(e => {
    console.log(e)
  })
}
