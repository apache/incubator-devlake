module.exports = {
  async findOrCreateCollection (db, collectionName) {
    const foundCollectionsCursor = await db.listCollections()
    const foundCollections = await foundCollectionsCursor.toArray()

    // check if Jira collection exists
    const collectionExists = foundCollections
      .some(collection => collection.name === collectionName)

    return collectionExists === true
      ? await db.collection(collectionName)
      : await db.createCollection(collectionName)
  }
}
