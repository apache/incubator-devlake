const assert = require('assert')
const gitlabProjectCollector = require('../../../src/collection/collectors/gitlab/projects')

describe('gitlab collector', () => {
	describe('Project collector', () => {
		describe('fetchProject', () => {
			it('Gets a project from gitlab with the expected name', async () => {
				const projectId = 28270340
				const expectedProjectName = 'test-project'
				let project = await gitlabProjectCollector.fetchProject(projectId)
				console.log('project', project);
				assert.equal(project.name, expectedProjectName)
			})
		})
		describe('fetchAllProjects', () => {
			it('Gets more than 0 projects from gitlab', async () => {
				let projects = await gitlabProjectCollector.fetchAllProjects()
				assert.equal(projects.length > 0, true)
			})
		})
		describe('fetchProjectRepoFiles', () => {
			it('Gets project repo files', async () => {
				const projectId = 28270340

				let repository = await gitlabProjectCollector.fetchProjectRepoCommits(projectId)
				console.log('repository', repository);
			})
		})
		describe('fetchProjectRepoTree', () => {
			it('Gets project repo tree', async () => {
				const projectId = 28270340

				let tree = await gitlabProjectCollector.fetchProjectRepoTree(projectId)
				console.log('tree', tree);
			})
		})
		describe.only('fetchProjectRepoFiles', () => {
			it('Gets project repo files', async () => {
				const projectId = 28270340
				let project = await gitlabProjectCollector.fetchProject(projectId)
				let tree = await gitlabProjectCollector.fetchProjectRepoTree(projectId)
				let files = await gitlabProjectCollector.fetchProjectFiles(projectId, tree, project.default_branch)
				console.log('files', files);
			})
		})
	})
})