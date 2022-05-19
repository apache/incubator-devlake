/// <reference types="cypress" />

context('Navbar', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('shows merico github icon link', () => {
    cy.get('.navbar')
      .should('have.class', 'bp3-navbar')
      .find('a[href="https://github.com/apache/incubator-devlake"]')
      .should('be.visible')
      .and('have.class', 'navIconLink')
  })

  it('shows merico email icon link', () => {
    cy.get('.navbar')
      .should('have.class', 'bp3-navbar')
      .find('a[href="mailto:hello@merico.dev"]')
      .should('be.visible')
      .and('have.class', 'navIconLink')
  })

  it('shows merico discord icon link', () => {
    cy.get('.navbar')
      .should('have.class', 'bp3-navbar')
      .find('a[href="https://discord.com/invite/83rDG6ydVZ"]')
      .should('be.visible')
      .and('have.class', 'navIconLink')
  })
})