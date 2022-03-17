/// <reference types="cypress" />

context('Data Integration Providers', () => {
  beforeEach(() => {
    cy.visit('/integrations')
  })

  describe('JIRA Data Provider', () => {
    it('provides access to jira integration', () => {
      cy.visit('/integrations/jira')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/jira integration/i)
    })

    it('displays connection sources data table', () => {
      cy.visit('/integrations/jira')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'ID')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })

    it('displays add connection button', () => {
      cy.visit('/integrations/jira')
      cy.get('button.bp3-button').contains('Add Connection')
        .should('be.visible')
    })

    it('displays refresh connections button', () => {
      cy.visit('/integrations/jira')
      cy.get('button.bp3-button').contains('Refresh Connections')
        .should('be.visible')
    })

  })

  describe('GitLab Data Provider', () => {
    it('provides access to gitlab integration', () => {
      cy.visit('/integrations/gitlab')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/gitlab integration/i)
    })
    it('displays connection sources data table', () => {
      cy.visit('/integrations/gitlab')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })
    it('limited to one (1) connection source', () => {
      cy.visit('/integrations/gitlab')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('cannot add a new connection', () => {
      cy.visit('/integrations/gitlab')
      cy.get('button.bp3-button').contains('Add Connection')
        .parent()
        .should('have.class', 'bp3-disabled')
        .should('have.attr', 'disabled')
    })
  })

  describe('GitHub Data Provider', () => {
    it('provides access to github integration', () => {
      cy.visit('/integrations/github')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/github integration/i)
    })
    it('displays connection sources data table', () => {
      cy.visit('/integrations/github')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })
    it('limited to one (1) connection source', () => {
      cy.visit('/integrations/github')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('cannot add a new connection', () => {
      cy.visit('/integrations/github')
      cy.get('button.bp3-button').contains('Add Connection')
        .parent()
        .should('have.class', 'bp3-disabled')
        .should('have.attr', 'disabled')
    })
  })

  describe('Jenkins Data Provider', () => {
    it('provides access to jenkins integration', () => {
      cy.visit('/integrations/jenkins')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/jenkins integration/i)
    })
    it('displays connection sources data table', () => {
      cy.visit('/integrations/jenkins')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })
    it('limited to one (1) connection source', () => {
      cy.visit('/integrations/jenkins')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('cannot add a new connection', () => {
      cy.visit('/integrations/jenkins')
      cy.get('button.bp3-button').contains('Add Connection')
        .parent()
        .should('have.class', 'bp3-disabled')
        .should('have.attr', 'disabled')
    })
  })


})