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

context('Manage Pipelines & Runs', () => {
  beforeEach(() => {
    cy.visit('/pipelines')
  })

  it('provides access/management of pipeline runs', () => {
    cy.get('ul.bp3-breadcrumbs')
      .find('a.bp3-breadcrumb-current')
      .contains(/all pipeline runs/i)
      .should('be.visible')
      .should('have.attr', 'href', '/pipelines')

    cy.get('.headlineContainer')
      .find('h1')
      .contains(/all pipeline runs/i)
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
      .contains(/create pipeline run/i)
      .should('be.visible')
      .as('createRunBtn')

    cy.get('@createRunBtn').click()
    cy.url().should('include', `${Cypress.config().baseUrl}/pipelines/create`)
  })

})