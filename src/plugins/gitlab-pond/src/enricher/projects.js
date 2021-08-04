const projectsCollector = require('../collector/projects')

async function enrich ({ rawDb, enrichedDb, projectId }) {
  if (!projectId) {
    throw new Error('Failed to enrich gitlab project, projectId is required')
  }

  console.info('INFO >>> gitlab enriching project', projectId)
  await enrichProjectById(rawDb, enrichedDb, projectId)
  console.info('INFO >>> gitlab enriching project done!', projectId)
}

async function enrichProjectById (rawDb, enrichedDb, projectId) {
  const projectsCollection = await projectsCollector.getCollection(rawDb)
  const project = await projectsCollection.findOne({ id: projectId })
  const enriched = mapResponseToSchema(project)
  await enrichedDb.GitlabProject.upsert(enriched)
}

function mapResponseToSchema (project) {
  return {
    name: project.name,
    id: project.id,
    pathWithNamespace: project.path_with_namespace,
    webUrl: project.web_url,
    visibility: project.visibility,
    openIssuesCount: project.open_issues_count,
    starCount: project.star_count
  }
}

module.exports = { enrich }
