require('module-alias/register')

const {
  findOrCreateCollection
} = require('../../../commondb')

const fetcher = require('./fetcher')

module.exports = {
  async collectProjectCommits (options, db) {
    for (let index = 0; index < options.projectIds.length; index++) {
      const projectId = options.projectIds[index]

      // TODO: this only get 20 commits... we need to page through all of them
      let response = await module.exports.fetchCollectionData('projects', projectId, 'repository/commits?with_stats=true&per_page=100')
      response = response.map(res => {
        return {
          projectId,
          ...res
        }
      })
      await module.exports.saveMany(response, db, 'gitlab_commits')
    }
  },

  async collectProjectsDetails (options, db) {
    for (let index = 0; index < options.projectIds.length; index++) {
      const projectId = options.projectIds[index]

      const response = await module.exports.fetchCollectionData('projects', projectId, '')
      await module.exports.saveOne(response, db, 'gitlab_projects')
    }
  },

  async saveOne (response, db, collectionName) {
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

  async saveMany (response, db, collectionName) {
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
    const requestUri = `${modelName}/${id}${uriComponent && '/' + uriComponent}`
    console.log('INFO: requestUri', requestUri)
    return fetcher.fetch(requestUri)
  },
  async findCollection (collectionName, where, db, limit = 99999999) {
    console.log(`INFO >>> ${collectionName} where`, where)
    const collection = await findOrCreateCollection(db, collectionName)
    const collectionDataCursor = await collection.find(where).limit(limit)
    return await collectionDataCursor.toArray()
  }
}
