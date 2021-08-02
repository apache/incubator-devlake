module.exports = {
  async findOrCreateCollection (db, collectionName, options = {}) {
    const foundCollectionsCursor = await db.listCollections()
    const foundCollections = await foundCollectionsCursor.toArray()

    // check if Jira collection exists
    const collectionExists = foundCollections
      .some(collection => collection.name === collectionName)

    return collectionExists
      ? await db.collection(collectionName)
      : await db.createCollection(collectionName, options)
  }
}
