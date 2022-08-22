/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

    it('can create a new connection', () => {
      cy.fixture('new-jira-connection').as('JIRAConnectionConnectionJSON')
      cy.intercept('POST', '/api/plugins/jira/connections', { statusCode: 201, body: '@JIRAConnectionConnectionJSON' }).as('createJIRAConnection')
      cy.visit('/integrations/jira')
      cy.get('button#btn-add-new-connection').click()
      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.disabled')

      cy.get('input#connection-name').type('TEST JIRA INSTANCE')
      cy.get('input#connection-endpoint').type('https://test-46f2c29a-2955-4fa6-8de8-fffeaf8cf8e0.atlassian.net/rest/')
      cy.get('input#connection-username').type('some username')
      cy.get('input#connection-password').type('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')
      cy.get('input#connection-proxy').type('http://proxy.localhost:8800')

      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.enabled')
        .click()

      cy.wait('@createJIRAConnection').its('response.statusCode').should('eq', 201)
      cy.url().should('include', '/integrations/jira')
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
    it('limited to one (1) connection', () => {
      cy.visit('/integrations/gitlab')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('can create a new connection', () => {
      cy.fixture('new-gitlab-connection').as('GitlabConnectionConnectionJSON')
      cy.intercept('POST', '/api/plugins/gitlab/connections', { statusCode: 201, body: '@GitlabConnectionConnectionJSON' }).as('createGitlabConnection')
      cy.visit('/integrations/gitlab')
      cy.get('button#btn-add-new-connection').click()
      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.disabled')

      cy.get('input#connection-name').type('TEST JIRA INSTANCE')
      cy.get('input#connection-endpoint').type('https://gitlab.com/api/')
      cy.get('input#connection-token').type('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')
      cy.get('input#connection-proxy').type('http://proxy.localhost:8800')

      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.enabled')
        .click()

      cy.wait('@createGitlabConnection').its('response.statusCode').should('eq', 201)
      cy.url().should('include', '/integrations/gitlab')
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
    it('limited to one (1) connection', () => {
      cy.visit('/integrations/github')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('can create a new connection', () => {
      cy.fixture('new-github-connection').as('GithubConnectionConnectionJSON')
      cy.intercept('POST', '/api/plugins/github/connections', { statusCode: 201, body: '@GithubConnectionConnectionJSON' }).as('createGithubConnection')
      cy.visit('/integrations/github')
      cy.get('button#btn-add-new-connection').click()
      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.disabled')

      cy.get('input#connection-name').type('TEST JIRA INSTANCE')
      cy.get('input#connection-endpoint').type('https://github.com/api/')
      cy.get('input#pat-id-0').type('xxxxx0')
      cy.get('input#connection-proxy').type('http://proxy.localhost:8800')

      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.enabled')
        .click()

      cy.wait('@createGithubConnection').its('response.statusCode').should('eq', 201)
      cy.url().should('include', '/integrations/github')
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
    it('limited to one (1) connection', () => {
      cy.visit('/integrations/jenkins')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('can create a new connection', () => {
      cy.fixture('new-jenkins-connection').as('JenkinsConnectionConnectionJSON')
      cy.intercept('POST', '/api/plugins/jenkins/connections', { statusCode: 201, body: '@JenkinsConnectionConnectionJSON' }).as('createJenkinsConnection')
      cy.visit('/integrations/jenkins')
      cy.get('button#btn-add-new-connection').click()
      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.disabled')

      cy.get('input#connection-name').type('TEST JIRA INSTANCE')
      cy.get('input#connection-endpoint').type('https://jenkins.com/api/')
      cy.get('input#connection-username').type('username')
      cy.get('input#connection-password').type('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')

      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.enabled')
        .click()

      cy.wait('@createJenkinsConnection').its('response.statusCode').should('eq', 201)
      cy.url().should('include', '/integrations/jenkins')
    })
  })
})
