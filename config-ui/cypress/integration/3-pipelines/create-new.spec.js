/// <reference types="cypress" />

context('Create New Pipelines Interface', () => {
  beforeEach(() => {
    cy.visit('/pipelines/create')
  })

  it('provides access to creating a new pipeline', () => {
    cy.get('ul.bp3-breadcrumbs')
      .find('a.bp3-breadcrumb-current')
      .contains(/run pipeline/i)
      .should('be.visible')
      .should('have.attr', 'href', '/pipelines/create')

    cy.get('.headlineContainer')
      .find('h1')
      .contains(/run new pipeline/i)
      .should('be.visible')
  })  

  it('has form control for pipeline name', () => {
    cy.get('h2')
    .contains(/pipeline name/i)
    .should('be.visible')

    cy.get('input#pipeline-name')
    .should('be.visible')
  })

  it('has plugin support for gitlab data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-gitlab')
      .should('be.visible')
  })

  it('has plugin support for github data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-github')
      .should('be.visible')
  })

  it('has plugin support for jenkins data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-jenkins')
      .should('be.visible')
  })
  
  it('has plugin support for jira data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-jira')
      .should('be.visible')
  })
    
  it('has plugin support for refdiff plugin provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-refdiff')
      .should('be.visible')
  })
  
  it('has plugin support for gitextractor plugin provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-gitextractor')
      .should('be.visible')
  })

  it('has form button control for running pipeline', () => {
    cy.get('.btn-run-pipeline')
      .should('be.visible')
  })

  it('has form button control for resetting pipeline', () => {
    cy.get('.btn-reset-pipeline')
      .should('be.visible')
  })

  it('has form button control for viewing all pipelines (manage)', () => {
    cy.get('.btn-view-jobs')
      .should('be.visible')
  })

  it('supports advanced-mode user interface options', () => {
    cy.get('.advanced-mode-toggleswitch')
      .should('be.visible')
      .find('.bp3-control-indicator')
      .click()

    cy.get('h2')
      .contains(/pipeline name \(advanced\)/i)
      .should('be.visible')
  })


})

context('RUN / Trigger New Pipelines', () => {
  beforeEach(() => {
    cy.visit('/pipelines/create')
  })

  it('can run a jenkins pipeline', () => {
    cy.fixture('new-jenkins-pipeline').then((newJenkinsPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newJenkinsPipelineJSON.ID}`, { body: newJenkinsPipelineJSON }).as('JenkinsPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newJenkinsPipelineJSON }).as('runJenkinsPipeline')
      cy.fixture('new-jenkins-pipeline-tasks').then((newJenkinsPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newJenkinsPipelineJSON.ID}/tasks`, { body: newJenkinsPipelineTasksJSON }).as('JenkinsPipelineTasks')
      })
    })

    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT JENKINS ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-jenkins')
      .should('be.visible')
      .click()
    
    cy.get('button#btn-run-pipeline').click()
    cy.wait('@JenkinsPipeline')
    cy.wait('@JenkinsPipelineTasks')
    cy.wait('@runJenkinsPipeline').then(({ response }) => {
      const NewJenkinsRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewJenkinsRun.ID}`)
    })
  })

  it('can run a gitlab pipeline', () => {
    cy.fixture('new-gitlab-pipeline').then((newGitlabPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newGitlabPipelineJSON.ID}`, { body: newGitlabPipelineJSON }).as('GitlabPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newGitlabPipelineJSON }).as('runGitlabPipeline')
      cy.fixture('new-gitlab-pipeline-tasks').then((newGitlabPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newGitlabPipelineJSON.ID}/tasks`, { body: newGitlabPipelineTasksJSON }).as('GitlabPipelineTasks')
      })
    })
    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT GITLAB ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-gitlab')
      .should('be.visible')
      .click()
    
    cy.get('.input-project-id').find('input').type('278964{enter}')
    
    cy.get('button#btn-run-pipeline').click()
    cy.wait('@GitlabPipeline')
    cy.wait('@GitlabPipelineTasks')
    cy.wait('@runGitlabPipeline').then(({ response }) => {
      const NewGitlabRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewGitlabRun.ID}`)
    })
  })

  it('can run a github pipeline', () => {
    cy.fixture('new-github-pipeline').then((newGithubPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newGithubPipelineJSON.ID}`, { body: newGithubPipelineJSON }).as('GithubPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newGithubPipelineJSON }).as('runGithubPipeline')
      cy.fixture('new-github-pipeline-tasks').then((newGithubPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newGithubPipelineJSON.ID}/tasks`, { body: newGithubPipelineTasksJSON }).as('GithubPipelineTasks')
      })
    })
    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT GITHUB ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-github')
      .should('be.visible')
      .click()
      .trigger('mouseleave')
    
    cy.get('input#owner').click().type('merico-dev', {force: true})
    cy.get('input#repository-name').type('lake')
    
    cy.get('button#btn-run-pipeline').click()
    cy.wait('@GithubPipeline')
    cy.wait('@GithubPipelineTasks')
    cy.wait('@runGithubPipeline').then(({ response }) => {
      const NewGithubRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewGithubRun.ID}`)
    })
  })


})