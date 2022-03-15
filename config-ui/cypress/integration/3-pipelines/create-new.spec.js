/// <reference types="cypress" />

context('Create New Pipelines', () => {
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
})