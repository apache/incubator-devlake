const express = require('express')
const path = require('path')
const app = express()
const SERVER_PORT = process.env.CONFIG_UI_PORT ?? 9000

app.use(express.static(path.join(__dirname, 'dist')))

app.get('/*', function (req, res) {
  res.sendFile(path.join(__dirname, 'dist', 'index.html'))
})

app.listen(SERVER_PORT, () => console.log(`lake / config-ui => listening on port : ${SERVER_PORT}`))