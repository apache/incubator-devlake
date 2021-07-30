require('module-alias/register')

const mongo = require('../util/mongo')

module.exports = {
  async enrich (rawDb, enrichedDb, options) {
    try {
      console.log('INFO: Gitlab Enrichment for projectId: ', options.projectId)
      await module.exports.saveProjectsToPsql(
        rawDb,
        enrichedDb,
        options.projectId
      )
      await module.exports.saveCommitsToPsqlBasedOnProjectId(
        rawDb,
        enrichedDb,
        options.projectId
      )
      await module.exports.saveMergeRequestsToPsqlBasedOnProjectId(
        rawDb,
        enrichedDb,
        options.projectId
      )
      console.log('Done enriching issues')
    } catch (error) {
      console.error(error)
    }
  },
  async saveMergeRequestsToPsqlBasedOnProjectId (rawDb, enrichedDb, projectId) {
    const {
      GitlabMergeRequest
    } = enrichedDb

    // find the project in mongo
    const mergeRequests = await mongo.findCollection('gitlab_merge_requests',
      { projectId: projectId }
      , rawDb)

    // mongo always returns an array
    const upsertPromises = []

    mergeRequests.forEach(mergeRequest => {
      mergeRequest = {
        projectId: mergeRequest.project_id,
        id: mergeRequest.id,
        numberOfReviewers: mergeRequest.reviewers && mergeRequest.reviewers.length,
        state: mergeRequest.state,
        title: mergeRequest.title,
        webUrl: mergeRequest.web_url,
        userNotesCount: mergeRequest.user_notes_count,
        workInProgress: mergeRequest.work_in_progress,
        sourceBranch: mergeRequest.source_branch,
        mergedAt: mergeRequest.merged_at,
        gitlabCreatedAt: mergeRequest.created_at,
        closedAt: mergeRequest.closed_at,
        mergedByUsername: mergeRequest.merged_by && mergeRequest.merged_by.username,
        description: mergeRequest.description,
        authorUsername: mergeRequest.author && mergeRequest.author.username
      }

      upsertPromises.push(GitlabMergeRequest.upsert(mergeRequest))
    })

    await Promise.all(upsertPromises)
  },
  async saveCommitsToPsqlBasedOnProjectId (rawDb, enrichedDb, projectId) {
    const {
      GitlabCommit
    } = enrichedDb

    // find the project in mongo
    const commits = await mongo.findCollection('gitlab_commits',
      { projectId: projectId }
      , rawDb)

    // mongo always returns an array
    const upsertPromises = []

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

      upsertPromises.push(GitlabCommit.upsert(commit))
    })

    await Promise.all(upsertPromises)
  },

  async saveProjectsToPsql (rawDb, enrichedDb, projectId) {
    const {
      GitlabProject
    } = enrichedDb

    // find the project in mongo
    const projects = await mongo.findCollection('gitlab_projects',
      { id: projectId }
      , rawDb)

    const upsertPromises = []

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

      upsertPromises.push(GitlabProject.upsert(project))
    })

    await Promise.all(upsertPromises)
  }
}
