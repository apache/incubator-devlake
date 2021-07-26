require('module-alias/register')

const { findOrCreateCollection } = require('../../../commondb')

const fetcher = require('./fetcher')

const collectionName = 'gitlab_project_repo_commits'

module.exports = {
  async collect (options) {
    try {
      const commitsResponse = await module.exports.fetchProjectRepoCommits(options.projectId)

      await module.exports.save({ commitsResponse, db: options.db })
    } catch (error) {
      console.log(error)
    }
  },
  async fetchProjectRepoCommits (projectId) {
     const requestUri = `projects/${projectId}/repository/commits?all=true&with_stats=true`

    return fetcher.fetch(requestUri)
  },
  async save ( {response, db} ){
    try {
      const promises = []
      const commitsCollection = await findOrCreateCollection(db, collectionName)

      response.commits.forEach(commit => {
        commit.primaryKey = Number(commit.id)

        promises.push(commitsCollection.findOneAndUpdate({
          primaryKey: commit.primaryKey
        }, {
          $set: commit
        }, {
          upsert: true
        }))
      })

      await Promise.all(promises)
    } catch (error) {
      console.error(error)
    }
  },
  async findCommits (where, db, limit = 99999999) {
    console.log('INFO >>> findCommits where', where)
    const commitsCollection = await findOrCreateCollection(db, collectionName)
    const foundCommitsCursor = await commitsCollection.find(where).limit(limit)
    return await foundCommitsCursor.toArray()
  }
}