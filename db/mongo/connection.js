require('module-alias/register')

const { MongoClient } = require('mongodb')
const MONGO_URI = require('@config/resolveConfig').mongo.connectionString

module.exports = {
  connect: async () => {
    try {
      console.log(MONGO_URI)
      const client = new MongoClient(MONGO_URI)
      await client.connect()
      const db = client.db()
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
  },
  async clearCollectionData (db, collectionName, where = {}) {
    try {
      const foundCollectionsCursor = await db.listCollections()
      const foundCollections = await foundCollectionsCursor.toArray()

      let collectionExists = false
      foundCollections.forEach(collection => {
        if (collection.name === collectionName) {
          collectionExists = true
        }
      })
      return collectionExists
        ? await db.collection(collectionName).deleteMany(where)
        : { exists: false }
    } catch (e) {
      console.error('MONGO.DB clearCollectionData >> ERROR: ', e)
    }
  }
}
