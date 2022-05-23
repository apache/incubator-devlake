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

context('Sidebar', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('shows merico application logo', () => {
    cy.get('.sidebar-card')
      .should('have.class', 'bp3-card')
      .find('img.logo')
      .should('be.visible')
      .and('have.attr', 'src', '/logo.svg')
  })

  it('shows grafana dashboards access button', () => {
    cy.get('.sidebar-card')
      .find('.dashboardBtn')
      .should('have.class', 'bp3-button bp3-outlined')
      .contains('View Dashboards')
  })

  it('displays apache 2.0 license notice', () => {
    cy.get('.sidebar-card')
      .should('have.class', 'bp3-card')
      .contains('Apache 2.0 License')
  })
  
  describe('Navigation System', () => {

    it('displays primary navigation menu', () => {
      cy.get('.sidebar-card ')
        .should('have.class', 'bp3-card')
        .find('.sidebarMenu')
        .first()
        .should('match', 'ul')
    })

    it('provides access to data integrations', () => {
      cy.get('.sidebar-card ')
        .should('have.class', 'bp3-card')
        .find('.sidebarMenu')
        .first()
        .should('match', 'ul')
        .contains('li', /data integrations/i)
        .should('be.visible')
    })

    it('provides access to triggers', () => {
      cy.get('.sidebar-card ')
        .should('have.class', 'bp3-card')
        .find('.sidebarMenu')
        .first()
        .should('match', 'ul')
        .contains('li', /triggers/i)
        .should('be.visible')
    })

    it('provides access to pipelines', () => {
      cy.get('.sidebar-card ')
        .should('have.class', 'bp3-card')
        .find('.sidebarMenu')
        .first()
        .should('match', 'ul')
        .contains('li', /pipelines/i)
        .should('be.visible')
    })

    it('provides external access to lake github documentation', () => {
      cy.get('.sidebar-card ')
        .should('have.class', 'bp3-card')
        .find('.sidebarMenu')
        .first()
        .should('match', 'ul')
        .contains('li', /documentation/i)
        .should('be.visible')
    })

    it('provides external access to merico network links', () => {
      cy.get('.sidebar-card ')
        .should('have.class', 'bp3-card')
        .find('.sidebarMenu')
        .first()
        .should('match', 'ul')
        .contains('li', /merico network/i)
        .should('be.visible')
    })


  })

})