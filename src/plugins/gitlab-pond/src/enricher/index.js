require('module-alias/register')

const collectionManager = require('../collector/collection-manager')

module.exports = {
  async enrich (rawDb, enrichedDb, options) {
    try {
      console.log('INFO: Gitlab Enrichment for projectIds: ', options.projectIds)
      await module.exports.saveProjectsToPsql(
        rawDb,
        enrichedDb,
        options.projectIds
      )
      await module.exports.saveCommitsToPsqlBasedOnProjectIds(
        rawDb,
        enrichedDb,
        options.projectIds
      )
      console.log('Done enriching issues')
    } catch (error) {
      console.error(error)
    }
  },

  async saveCommitsToPsqlBasedOnProjectIds (rawDb, enrichedDb, projectIds) {
    const {
      GitlabCommit
    } = enrichedDb

    // find the project in mongo
    const commits = await collectionManager.findCollection('gitlab_commits',
      { projectId: { $in: projectIds } }
      , rawDb)

    // mongo always returns an array
    const creationPromises = []
    const updatePromises = []

    commits.forEach(commit => {
      commit = {
        projectId: commit.projectId,
        id: commit.id,
        shortId: commit.short_id,
        title: commit.title,
        message: commit.message,
        authorName: commit.author_name,
        authorEmail: commit.author_email,
        authoredDate: commit.authored_date,
        committerName: commit.committer_name,
        committerEmail: commit.committer_email,
        committedDate: commit.committed_date,
        webUrl: commit.web_url,
        additions: commit.stats.additions,
        deletions: commit.stats.deletions,
        total: commit.stats.total
      }

      creationPromises.push(GitlabCommit.findOrCreate({
        where: {
          id: commit.id
        },
        defaults: commit
      }))

      updatePromises.push(GitlabCommit.update(commit, {
        where: {
          id: commit.id
        }
      }))
    })

    await Promise.all(creationPromises)
    await Promise.all(updatePromises)
  },

  async saveProjectsToPsql (rawDb, enrichedDb, projectIds) {
    const {
      GitlabProject
    } = enrichedDb

    // find the project in mongo
    const projects = await collectionManager.findCollection('gitlab_projects',
      { id: { $in: projectIds } }
      , rawDb)

    const creationPromises = []
    const updatePromises = []

    // mongo always returns an array
    projects.forEach(project => {
      project = {
        name: project.name,
        id: project.id,
        pathWithNamespace: project.path_with_namespace,
        webUrl: project.web_url,
        visibility: project.visibility,
        openIssuesCount: project.open_issues_count,
        starCount: project.star_count
      }

      creationPromises.push(GitlabProject.findOrCreate({
        where: {
          id: project.id
        },
        defaults: project
      }))

      updatePromises.push(GitlabProject.update(project, {
        where: {
          id: project.id
        }
      }))
    })

    await Promise.all(creationPromises)
    await Promise.all(updatePromises)
  }
}
