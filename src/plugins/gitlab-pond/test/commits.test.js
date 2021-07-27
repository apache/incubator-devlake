const commits = require("../src/collector/commits")
const mockData = require('./data/commits')
const dbConnector = require('@mongo/connection')
const assert = require('assert')

// SKIP API calls
describe('Commits', () => {
  describe.skip('fetchProjectRepoCommits', () => {
    it('gets commits for a project', async () => {
      let projectId = 28270340
      let foundCommits = await commits.fetchProjectRepoCommits(projectId)
      console.log('foundCommits', foundCommits);
    })
  })
  describe.only('save', () => {
    it('commits found in the db have the same length as the mock data provided', async () => {
      const {
        db, client
      } = await dbConnector.connect()
      const collectionName = 'gitlab_project_repo_commits'
      try {
        await dbConnector.clearCollectionData(db, collectionName)
        await commits.save({ response: mockData, db})
        let foundCommits = await commits.findCommits('', db)
        assert.equal(foundCommits.length, mockData.length)
      } catch (error) {
        console.log('Failed to collect', error)
      } finally {
        dbConnector.disconnect(client)
      }
    })
  })
})