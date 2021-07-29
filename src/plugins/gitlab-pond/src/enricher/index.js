require('module-alias/register')

const collectionManager = require('../collector/collection-manager')

module.exports = {
  async enrich(rawDb, enrichedDb, options) {
    try {
      console.log('INFO: Gitlab Enrichment for projectId: ', options.projectId)
      await module.exports.saveProjectsToPsql(
        rawDb,
        enrichedDb,
        options.projectId
      )
      console.log('Done enriching issues')
    } catch (error) {
      console.error(error)
    }
  },

  async saveProjectsToPsql(rawDb, enrichedDb, projectId) {
    const {
      GitlabProject
    } = enrichedDb

    // find the project in mongo
    let project = await collectionManager.findCollection('projects', {
      'id': projectId
    }, rawDb)

    // mongo always returns an array
    project = project[0]

    project = {
      name: project.name,
      id: project.id,
      pathWithNamespace: project.path_with_namespace,
      webUrl: project.web_url,
      visibility: project.visibility,
      openIssuesCount: project.open_issues_count,
      starCount: project.star_count,
    }

    console.log('JON >>> project', project)
    // save the project in psql
    await GitlabProject.findOrCreate({
      where: {
        id: project.id
      },
      defaults: project
    })

    await GitlabProject.update(project, {
      where: {
        id: project.id
      }
    })
  },
}