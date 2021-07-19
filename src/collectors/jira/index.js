const { MongoClient } = require('mongodb')
const connection = require('@config/resolveConfig').mongo.connectionString
const client = new MongoClient(connection)

const issues = require('./issues')

module.exports = {
  async collect ({ projectId }) {
    await client.connect()
    const db = await client.db()

    await issues.collectIssues(db, projectId)
  },
  issues
}
