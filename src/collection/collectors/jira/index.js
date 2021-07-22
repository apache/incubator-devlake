const dbConnector = require('@mongo/connection')
const issues = require('./issues')
const changelogs = require('./changelogs')

module.exports = {
  async collect({
    projectId
  }) {
    console.log('Jira Collection, projectId:', projectId)
    const {
      db, client
    } = await dbConnector.connect()
    
    try {
      await issues.collect({db, projectId})
      await changelogs.collect({db, projectId})
      console.log('INFO >>> done collecting')
    } catch (error) {
      console.log('>>> error', error)
    } finally {
      dbConnector.disconnect(client)
    }
  }
}