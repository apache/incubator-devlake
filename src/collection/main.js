const express = require('express')

const app = express()
const port = 3001

app.get('/', async (req, res) => {
  res.send("Let's collect!")
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})