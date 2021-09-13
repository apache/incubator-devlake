import axios from 'axios'

export default function handler(req, res) {

  axios.post('http://localhost:8080/task', req.body ).then(res => {
    console.log(res.data)
  }).catch(e => {
    console.log(e)
  })
}
