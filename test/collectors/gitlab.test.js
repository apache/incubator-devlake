const gitlabProjectCollector = require('../../src/collection/collectors/gitlab/projects')

describe('gitlab collector', () => {
	describe('Project collector', () => {
    it('Gets a project from gitlab', async () => {
			const projectId = '/28270340'
			let project = await gitlabProjectCollector.fetchProject(projectId)
			console.log('project', project);
    })
	})
})