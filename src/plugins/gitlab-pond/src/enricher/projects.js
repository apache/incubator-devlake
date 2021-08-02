const projectsCollector = require('../collector/projects')

async function enrich ({ rawDb, enrichedDb, projectId }) {
  if (!projectId) {
    throw new Error('Failed to enrich gitlab project, projectId is required')
  }

  await enrichProjectById(rawDb, enrichedDb, projectId)
}

async function enrichProjectById (rawDb, enrichedDb, projectId) {
  console.info('INFO >>> gitlab enriching project', projectId)
  const projectsCollection = await projectsCollector.getCollection(rawDb)
  const project = await projectsCollection.findOne({ id: Number(projectId) })
  const enriched = {
    name: project.name,
    id: project.id,
    pathWithNamespace: project.path_with_namespace,
    webUrl: project.web_url,
    visibility: project.visibility,
    openIssuesCount: project.open_issues_count,
    starCount: project.star_count
  }
  await enrichedDb.GitlabProject.upsert(enriched)
  console.info('INFO >>> gitlab enriching project done!', projectId)
}

module.exports = { enrich }
