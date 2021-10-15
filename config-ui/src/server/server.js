const express = require('express')
const fs = require('fs')
const os = require('os')
const path = require('path')
const app = express()

app.use(express.static(__dirname))

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

app.get('/api/setenv/:key/:value', (req, res) => {

  const key = req.params.key
  const value = req.params.value

  const envFilePath = process.env.ENV_FILEPATH || path.join(process.cwd(), 'data', '../../.env')
  const readEnvVars = () => fs.readFileSync(envFilePath, "utf-8").split(os.EOL)

  const envVars = readEnvVars()
  const targetLine = envVars.find((line) => line.split("=")[0] === key)

  if (targetLine !== undefined) {
    const targetLineIndex = envVars.indexOf(targetLine)
    envVars.splice(targetLineIndex, 1, `${key}=${value}`)
  }
  else {
    envVars.push(`${key}=${value}`)
  }

  fs.writeFileSync(envFilePath, envVars.join(os.EOL))
  res.status(200).json({ key: value, status: 'updated' })
})

app.listen(4000, () => console.log(`Live on port 4000`))
