const express = require('express')
const bodyParser = require('body-parser')

const dispatch = require('./dispatch')

const app = express()
app.use(bodyParser.json())

const port = 3001

app.get('/', async (req, res) => {
  res.send('Collector is running')
})

app.post('/', async (req, res) => {
  await dispatch.createJobs(req.body[0])

  res.status(200).send(req.body)
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})
