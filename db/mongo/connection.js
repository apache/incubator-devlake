const { MongoClient } = require('mongodb')
const MONGO_URI = require('@config/resolveConfig').mongo.connectionString
const DEBUG = process.env.DEBUG

module.exports = {
	connect: async (database = 'test') => {
		try {
			const client = new MongoClient(MONGO_URI)
	  	await client.connect()
	  	const db = client.db(database)
	  	return { client, db }
  	} catch (e) {
  		console.log('MONGO.DB connect() >> ERROR: ', e)
  	}
	},
	disconnect: (client) => {
		try {
			client.close()
		} catch (e) {
			console.log('MONGO.DB disconnect() >> ERROR: ', e)
		}
	},
	createCollection: async (db, collection, options = { capped : true, size : 5242880, max : 5000 }) => {
		try {
			const collection = await db.createCollection(collection, options)
			return collection
		} catch (e) {
			console.log('MONGO.DB createCollection() >> ERROR: ', e)
		}
	}
}
