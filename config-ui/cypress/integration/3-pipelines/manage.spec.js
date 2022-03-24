/// <reference types="cypress" />

context('Manage Pipelines & Runs', () => {
  beforeEach(() => {
    cy.visit('/pipelines')
  })

  it('provides access/management of pipeline runs', () => {
    cy.get('ul.bp3-breadcrumbs')
      .find('a.bp3-breadcrumb-current')
      .contains(/manage pipeline runs/i)
      .should('be.visible')
      .should('have.attr', 'href', '/pipelines')

    cy.get('.headlineContainer')
      .find('h1')
      .contains(/pipeline runs/i)
      .should('be.visible')
  })

  it('displays pipeline runs data table', () => {
    cy.get('.pipelines-table')
      .should('have.class', 'bp3-html-table')
      .should('be.visible')
      .find('thead')
      .find('th')
      .should('contain', 'ID')
      .should('contain', 'Pipeline Name')
      .should('contain', 'Duration')
      .should('contain', 'Status')
  })

  it('displays pipelines filter controls panel', () => {
    cy.get('.filter-status-group')
      .should('have.class', 'bp3-button-group')
      .should('be.visible')
      .children().should('have.length', 4)
  })

  it('displays data table pagination controls', () => {
    cy.get('.operations.panel')
      .find('.pagination-controls')
      .should('be.visible')

    cy.get('.btn-prev-page').should('be.visible')
    cy.get('.btn-next-page').should('be.visible')
    cy.get('.btn-select-page-size').should('be.visible')
  })

  it('has action to create new pipeline run', () => {
    cy.get('.bp3-button')
      .should('have.class', 'bp3-intent-primary')
      .contains(/create run/i)
      .should('be.visible')
      .as('createRunBtn')

    cy.get('@createRunBtn').click()
    cy.url().should('include', `${Cypress.config().baseUrl}/pipelines/create`)
  })

})