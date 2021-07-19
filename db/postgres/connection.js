const { Sequelize } = require('sequelize')
const defaultDbConfig = require('@config/resolveConfig').postgres.connectionString
const DEBUG = process.env.DEBUG

module.exports = {
  // variables: ['Hello world!']
  execute: async (sql, variables, transaction, dbConfig = defaultDbConfig) => {
    const options = {
      // TODO: try to put paging here for all queries
      // Maybe start date end date and sort too
      replacements: variables
    }
    if (transaction) {
      options.transaction = transaction
    }

    let connection
    try {
      connection = module.exports.sequelize(dbConfig)
      const data = await connection.query(sql, options)

      // If we do not close these manually, we eventually run out of connections
      await connection.close()

      return data
    } catch (error) {
      console.error(error)
      if (connection) {
        await connection.close()
      }
      throw error
    }
  },
  sequelize: (dbConfig) => {
    return new Sequelize(dbConfig.database, dbConfig.username, dbConfig.password, {
      host: dbConfig.host,
      dialect: 'postgres',
      port: dbConfig.port,
      logging: DEBUG ? console.log : false,
      pool: {
        max: 100,
        min: 0,
        acquire: 30000,
        idle: 10000
      }
    })
  }
}
