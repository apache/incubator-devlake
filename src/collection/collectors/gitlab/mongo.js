const dbConnector = require('@mongo/connection')
const collectionName = 'gitlab_projects'

const mongo = {
  async storeRawData(data){
    const { client, db } = await dbConnector.connect()
    try {
      const projectCollection = await dbConnector.findOrCreateCollection(db, collectionName)
      await projectCollection.insert(data)
    } catch (error) {
      console.error(error)
    } finally {
      dbConnector.disconnect(client)
    }
  }
}

module.exports = mongo