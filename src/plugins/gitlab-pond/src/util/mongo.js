const { findOrCreateCollection } = require('../../../commondb')

module.exports = {
  async findCollection (collectionName, where, db, limit = 99999999) {
    const collection = await findOrCreateCollection(db, collectionName)
    const collectionDataCursor = await collection.find(where).limit(limit)
    return await collectionDataCursor.toArray()
  }
}
