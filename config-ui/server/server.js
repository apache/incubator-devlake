const express = require('express')
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

app.listen(4000, () => console.log(`Live on port 4000`))
