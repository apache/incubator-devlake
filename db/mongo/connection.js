require('module-alias/register')

const { MongoClient } = require('mongodb')
const MONGO_URI = require('@config/resolveConfig').mongo.connectionString

module.exports = {
  connect: async (database = 'test') => {
    try {
      const client = new MongoClient(MONGO_URI)
      await client.connect()
      const db = client.db(database)
      return { client, db }
    } catch (e) {
      console.log('MONGO.DB connect() >> ERROR: ', e)
    }
  },
  disconnect: (client) => {
    try {
      console.log('INFO >>> closing mongo connection')
      client.close()
    } catch (e) {
      console.log('MONGO.DB disconnect() >> ERROR: ', e)
    }
  },
  async findOrCreateCollection (db, collectionName, options = {}) {
    try {
      const foundCollectionsCursor = await db.listCollections()
      const foundCollections = await foundCollectionsCursor.toArray()

      // check if Jira collection exists
      const collectionExists = foundCollections
        .some(collection => collection.name === collectionName)

      return collectionExists
        ? await db.collection(collectionName)
        : await db.createCollection(collectionName, options)
    } catch (e) {
      console.log('MONGO.DB createCollection() >> ERROR: ', e)
    }
  }
}
