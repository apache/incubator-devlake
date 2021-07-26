const commits = require("../src/collector/commits")
const mockData = require('./data/commits')

describe('Commits', () => {
  describe('fetchProjectRepoCommits', () => {
    it('gets commits for a project', async () => {
      let projectId = 28270340
      let foundCommits = await commits.fetchProjectRepoCommits(projectId)
      console.log('foundCommits', foundCommits);
    })
  })
  describe.only('save', () => {
    it('stores commits for a project', async () => {
      let db = ''
      let foundCommits = await commits.save({response: mockData, db})
      console.log('foundCommits', foundCommits);
    })
  })
})