require('module-alias/register')
const express = require('express')
const bodyParser = require('body-parser')
const config = require('@config/resolveConfig').api || {}
const dispatch = require('./dispatch')

const app = express()
app.use(bodyParser.json())

const port = 3001

app.get('/', async (req, res) => {
  res.send('Collector is running')
})

app.post('/', async (req, res) => {
  if (config.token && req.headers['x-token'] !== config.token) {
    return res.status(401).json({ message: 'UNAUTHORIZED' })
  }

  await dispatch.createJobs(req.body)
  res.status(200).send(req.body)
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})
