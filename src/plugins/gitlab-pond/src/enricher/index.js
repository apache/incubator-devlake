require('module-alias/register')

const collectionManager = require('../collector/collection-manager')

module.exports = {
  async enrich(rawDb, enrichedDb, options) {
    try {
      console.log('INFO: Gitlab Enrichment for projectId: ', options.projectId)
      await module.exports.saveProjectsToPsql(
        rawDb,
        enrichedDb,
        options.projectIds
      )
      console.log('Done enriching issues')
    } catch (error) {
      console.error(error)
    }
  },

  async saveProjectsToPsql(rawDb, enrichedDb, projectIds) {
    const {
      GitlabProject
    } = enrichedDb

    // find the project in mongo
    let projects = await collectionManager.findCollection('projects', 
      { id: { $in: projectIds } }
    , rawDb)

    // mongo always returns an array
    creationPromises = []
    updatePromises = []

    projects.forEach(project => {
      project = {
        name: project.name,
        id: project.id,
        pathWithNamespace: project.path_with_namespace,
        webUrl: project.web_url,
        visibility: project.visibility,
        openIssuesCount: project.open_issues_count,
        starCount: project.star_count,
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
  },
}