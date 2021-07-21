const assert = require('assert')
const gitlabProjectCollector = require('../../../src/collection/collectors/gitlab/projects')

describe('gitlab collector', () => {
	describe('Project collector', () => {
		describe('fetchProject', () => {
			it('Gets a project from gitlab with the expected name', async () => {
				const projectId = 28270340
				const expectedProjectName = 'test-project'
				let project = await gitlabProjectCollector.fetchProject(projectId)
				assert.equal(project.name, expectedProjectName)
			})
		})
		describe('fetchAllProjects', () => {
			it('Gets more than 0 projects from gitlab', async () => {
				let projects = await gitlabProjectCollector.fetchAllProjects()
				assert.equal(projects.length > 0, true)
			})
		})
	})
})