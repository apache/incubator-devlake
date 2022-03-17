/// <reference types="cypress" />

context('Pipeline RUN Activity', () => {
  beforeEach(() => {
    cy.fixture('pipelines').then((pipelinesJSON) => {
      cy.intercept('GET', '/api/pipelines', { body: pipelinesJSON }).as('getPipelines')
    })
    cy.fixture('pipeline-activity').then((pipelinesActivityJSON) => {
      cy.intercept('GET', '/api/pipelines/100', { body: pipelinesActivityJSON }).as('getPipelineActivity')
    })    
    cy.fixture('pipeline-tasks').then((pipelinesTasksJSON) => {
      cy.intercept('GET', '/api/pipelines/100/tasks', { body: pipelinesTasksJSON }).as('getPipelineActivityTasks')
    })    
    cy.visit('/pipelines')
  })

  it('provides access to monitor pipeline activity', () => {
    cy.wait('@getPipelines').then(({ response }) => {
      expect(response.body.count).to.eq(10)
      expect(response.body.pipelines.length).to.eq(response.body.count)
      let run = response.body.pipelines[0]
      cy.visit(`/pipelines/activity/${run.ID}`)
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/pipeline activity/i)
        .should('be.visible')
      cy.get('.stage-panel-card')
        .should('be.visible')
    })
  })

  it('shows pipeline name', () => {
    cy.wait('@getPipelines').then(({ response }) => {
      const Run = response.body.pipelines[0]
      cy.visit(`/pipelines/activity/${Run.ID}`)
    })
    cy.wait('@getPipelineActivity').then(({ response }) => {
      const Activity = response.body
      cy.get('.pipeline-activity')
        .find('.pipeline-name')
        .contains(Activity.name)
        .should('be.visible')
    })
  })

  it('shows pipeline status', () => {
    cy.wait('@getPipelines').then(({ response }) => {
      const Run = response.body.pipelines[0]
      cy.visit(`/pipelines/activity/${Run.ID}`)
    })
    cy.wait('@getPipelineActivity').then(({ response }) => {
      const Activity = response.body
      cy.get('.pipeline-activity')
        .find('.pipeline-status')
        .contains(Activity.status.split('_')[1])
        .should('be.visible')

    })
  })

  it('shows pipeline duration', () => {
    cy.wait('@getPipelines').then(({ response }) => {
      const Run = response.body.pipelines[0]
      cy.visit(`/pipelines/activity/${Run.ID}`)
    })
    cy.wait('@getPipelineActivity').then(({ response }) => {
      const Activity = response.body
      cy.get('.pipeline-activity')
        .find('.pipeline-duration')
        .should('be.visible')

    })
  })

  it('shows pipeline task activity list', () => {
    cy.wait('@getPipelines').then(({ response }) => {
      const Run = response.body.pipelines[0]
      cy.visit(`/pipelines/activity/${Run.ID}`)
    })
    cy.wait('@getPipelineActivity').then(({ response }) => {
      const Activity = response.body
      cy.get('.pipeline-task-activity')
        .children()
        .should('have.length', Activity.tasks[0].length)
    })
  })

  it('displays provider run settings and configuration', () => {
    cy.wait('@getPipelines').then(({ response }) => {
      const Run = response.body.pipelines[0]
      cy.visit(`/pipelines/activity/${Run.ID}`)
    })
    cy.wait('@getPipelineActivity').then(({ response }) => {
      const Activity = response.body
      cy.get('.run-settings')
        .should('be.visible')
      Activity.tasks[0].forEach(task => {
        cy.get('.run-settings')
          .find(`.${task.plugin.toLowerCase()}-settings`)
          .should('be.visible')
      })
    })
  })

  it('has pipeline code inspector', () => {
    cy.wait('@getPipelines').then(({ response }) => {
      const Run = response.body.pipelines[0]
      cy.visit(`/pipelines/activity/${Run.ID}`)
    })
    cy.get('.btn-inspect-pipeline')
      .should('be.visible')
      .click()

    cy.get('.drawer-json-inspector')
      .should('be.visible')
      .find('.bp3-drawer-header')
      .contains(/inspect run/i)
  })
})