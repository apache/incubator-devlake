const express = require('express')
const fs = require('fs')
const os = require('os')
const cors = require('cors')
const path = require('path')
const dotenv = require('dotenv')
const axios = require('axios')
const app = express()
const DEVLAKE_ENDPOINT = require('./config').DEVLAKE_ENDPOINT
const CLIENT_ROOT = require('./config').CLIENT_ROOT

app.use(express.static(__dirname))
app.use(cors({
  origin: CLIENT_ROOT
}))

// Main pages

app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'index.html'))
})

app.get('/triggers', (req, res) => {
  res.sendFile(path.join(__dirname, 'triggers.html'))
})

// Plugins

app.get('/plugins/jira', (req, res) => {
  res.sendFile(path.join(__dirname, 'plugins/jira.html'))
})

app.get('/plugins/gitlab', (req, res) => {
  res.sendFile(path.join(__dirname, 'plugins/gitlab.html'))
})

app.get('/plugins/jenkins', (req, res) => {
  res.sendFile(path.join(__dirname, 'plugins/jenkins.html'))
})

// Api

app.get('/api/triggers/task', async (req, res) => {
  const r = await axios.post(`${DEVLAKE_ENDPOINT}/task`, req.body)
  res.json(r.data)
})

app.get('/api/getenv', async (req, res) => {
  const filePath = process.env.ENV_FILEPATH || path.join(process.cwd(), 'data', '../../../../.env')

  try {
    const fileData = fs.readFileSync(filePath)
    const env = dotenv.parse(fileData)

    return res.status(200).json(env)
  } catch (e) {
    console.error('Could not read env file', e)
    return res.status(500).send(e)
  }
})

app.get('/api/setenv/:key/:value', (req, res) => {
  const key = req.params.key
  const value = req.params.value

  const envFilePath = process.env.ENV_FILEPATH || path.join(process.cwd(), 'data', '../../../../.env')

  console.log(key, value, envFilePath)

  const readEnvVars = () => fs.readFileSync(envFilePath, 'utf-8').split(os.EOL)

  const envVars = readEnvVars()
  const targetLine = envVars.find((line) => line.split('=')[0] === key)

  if (targetLine !== undefined) {
    const targetLineIndex = envVars.indexOf(targetLine)
    envVars.splice(targetLineIndex, 1, `${key}=${value}`)
  } else {
    envVars.push(`${key}=${value}`)
  }

  fs.writeFileSync(envFilePath, envVars.join(os.EOL))
  res.status(200).json({ key: value, status: 'updated' })
})

app.listen(5000, () => console.log('Live on port 5000'))
