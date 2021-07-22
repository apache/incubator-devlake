const dbConnector = require('@mongo/connection')

const mongo = {
  async storeRawData (data, collectionName) {
    const { client, db } = await dbConnector.connect()
    try {
      const collection = await dbConnector.findOrCreateCollection(db, collectionName)
      await collection.insert(data)
    } catch (error) {
      console.error(error)
    } finally {
      dbConnector.disconnect(client)
    }
  }
}

module.exports = mongo
