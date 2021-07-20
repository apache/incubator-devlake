const express = require('express')
const bodyParser = require('body-parser')

const dispatch = require('./dispatch')

const app = express()
app.use(bodyParser.json())

const port = 3000

app.get('/', async (req, res) => {
  res.send("Let's enrich!")
})

app.post('/', async (req, res) => {
  await dispatch.createJob(req.body)

  res.status(200).send(req.body)
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})
