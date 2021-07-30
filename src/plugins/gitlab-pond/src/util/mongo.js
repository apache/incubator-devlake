module.exports = {
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