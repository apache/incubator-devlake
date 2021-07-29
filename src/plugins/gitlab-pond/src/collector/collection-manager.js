require('module-alias/register')

const { findOrCreateCollection } = require('../../../commondb')

const fetcher = require('./fetcher')

module.exports = {
  async collectProjectsDetails(options, db) {
    let modelName = 'projects'

    for (let index = 0; index < options.projectIds.length; index++) {
      const projectId = options.projectIds[index];
      
      let response = await module.exports.fetchCollectionData(modelName, projectId, '')
      await module.exports.saveOne(response, db, modelName)
    }
  },

  async saveOne (response, db, collectionName){
    try {
      const collection = await findOrCreateCollection(db, collectionName)
      response.primaryKey = response.id
  
      await collection.findOneAndUpdate({
        primaryKey: response.primaryKey
      }, {
        $set: response
      }, {
        upsert: true
      })
    } catch (error) {
      console.error(error)
    }
  },

  async saveMany (response, db, collectionName ){
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
    console.log('INFO: requestUri', requestUri);
   return fetcher.fetch(requestUri)
 },
  async findCollection (collectionName, where, db, limit = 99999999) {
    console.log(`INFO >>> ${collectionName} where`, where)
    const collection = await findOrCreateCollection(db, collectionName)
    const collectionDataCursor = await collection.find(where).limit(limit)
    return await collectionDataCursor.toArray()
  }
}