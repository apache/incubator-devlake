require('module-alias/register')

const { findOrCreateCollection } = require('../../../commondb')

const fetcher = require('./fetcher')

module.exports = {
  async collect (options) {
    try {
      let { id, db, modelName, uriComponent } = options
      const response = await module.exports.fetchCollectionData(modelName, id, uriComponent)

      await module.exports.save({ response, db })
    } catch (error) {
      console.log(error)
    }
  },
  async save ( {response, db}, collectionName ){
    try {
      const promises = []
      const collection = await findOrCreateCollection(db, collectionName)
      response.forEach(item => {
        item.primaryKey = item.id

        promises.push(collection.findOneAndUpdate({
          primaryKey: item.primaryKey
        }, {
          $set: item
        }, {
          upsert: true
        }))
      })

      await Promise.all(promises)
    } catch (error) {
      console.error(error)
    }
  },
  async fetchCollectionData (modelName, id, uriComponent = '') {
    const requestUri = `${modelName}/${id}${uriComponent && `/` + uriComponent}`
    console.log('requestUri', requestUri);
   return fetcher.fetch(requestUri)
 },
  async findCommits (collectionName, where, db, limit = 99999999) {
    console.log(`INFO >>> ${collectionName} where`, where)
    const collection = await findOrCreateCollection(db, collectionName)
    const collectionDataCursor = await collection.find(where).limit(limit)
    return await collectionDataCursor.toArray()
  }
}