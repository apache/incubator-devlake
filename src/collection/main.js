require('module-alias/register')
const express = require('express')
const bodyParser = require('body-parser')
const config = require('@config/resolveConfig').lake || {}
const dispatch = require('./dispatch')

const app = express()
app.use(bodyParser.json())

const port = process.env.COLLECTION_PORT || 3001
const host = process.env.COLLECTION_HOST || 'localhost'

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

app.listen(port, host, () => {
  console.log(`Collection API listening at http://${host}:${port}`)
})
