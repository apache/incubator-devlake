const assert = require('assert')
const gitlabProjectCollector = require('../../../src/collection/collectors/gitlab/projects')

// TODO: these tests are great but lets not have them actually go to the third party api
describe.skip('gitlab collector', () => {
  describe('Project collector', () => {
    describe('fetchProject', () => {
      it('Gets a project from gitlab with the expected name', async () => {
        const projectId = 28270340
        const expectedProjectName = 'test-project'
        const project = await gitlabProjectCollector.fetchProject(projectId)
        console.log('project', project)
        assert.equal(project.name, expectedProjectName)
      })
    })
    describe('fetchAllProjects', () => {
      it('Gets more than 0 projects from gitlab', async () => {
        const projects = await gitlabProjectCollector.fetchAllProjects()
        assert.equal(projects.length > 0, true)
      })
    })
    describe('fetchProjectRepoFiles', () => {
      it('Gets project repo files', async () => {
        const projectId = 28270340

        const repository = await gitlabProjectCollector.fetchProjectRepoCommits(projectId)
        console.log('repository', repository)
      })
    })
    describe('fetchProjectRepoTree', () => {
      it('Gets project repo tree', async () => {
        const projectId = 28270340

        const tree = await gitlabProjectCollector.fetchProjectRepoTree(projectId)
        console.log('tree', tree)
      })
    })
    describe.only('fetchProjectRepoFiles', () => {
      it('Gets project repo files', async () => {
        const projectId = 28270340
        const project = await gitlabProjectCollector.fetchProject(projectId)
        const tree = await gitlabProjectCollector.fetchProjectRepoTree(projectId)
        const files = await gitlabProjectCollector.fetchProjectFiles(projectId, tree, project.default_branch)
        console.log('files', files)
      })
    })
  })
})
