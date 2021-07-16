const express = require('express')

const app = express()
const port = 3000

app.get('/', async (req, res) => {
  res.send("Let's enrich!w")
})

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})