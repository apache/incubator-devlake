/// <reference types="cypress" />

context('API Network Requests', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('listens for network ping request', () => {
    cy.request(`${Cypress.env('apiUrl')}ping`)
      .should((response) => {
        expect(response.status).to.eq(200)
      })
  })

  it('provides jira connection resources', () => {
    cy.request(`${Cypress.env('apiUrl')}plugins/jira/connections`)
      .should((response) => {
        expect(response.status).to.eq(200)
        expect(response.headers).to.have.property('content-type').and.to.eq('application/json; charset=utf-8')
        expect(response.body).to.be.an('array')
        expect(response.body[0]).to.have.property('createdAt')
        expect(response.body[0]).to.have.property('updatedAt')
        expect(response.body[0]).to.have.property('id')
        expect(response.body[0]).to.have.property('name')
        expect(response.body[0]).to.have.property('endpoint')
        expect(response.body[0]).to.have.property('basicAuthEncoded')
        expect(response.body[0]).to.have.property('epicKeyField')
        expect(response.body[0]).to.have.property('storyPointField')
        expect(response.body[0]).to.have.property('remotelinkCommitShaPattern')
        expect(response.body[0]).to.have.property('proxy')
      })
  })

  it('provides jenkins connection resources', () => {
    cy.request(`${Cypress.env('apiUrl')}plugins/jenkins/connections`)
      .should((response) => {
        expect(response.status).to.eq(200)
        expect(response.headers).to.have.property('content-type').and.to.eq('application/json; charset=utf-8')
        expect(response.body).to.be.an('array')
        expect(response.body[0]).to.have.property('ID').and.to.eq(1)
        expect(response.body[0]).to.have.property('Name').and.to.eq('Jenkins')
        expect(response.body[0]).to.have.property('Endpoint')
        expect(response.body[0]).to.have.property('Username')
        expect(response.body[0]).to.have.property('Password') 
        expect(response.body[0]).to.have.property('Proxy')
      })
  })

  it('provides gitlab connection resources', () => {
    cy.request(`${Cypress.env('apiUrl')}plugins/gitlab/connections`)
      .should((response) => {
        expect(response.status).to.eq(200)
        expect(response.headers).to.have.property('content-type').and.to.eq('application/json; charset=utf-8')
        expect(response.body).to.be.an('array')
        expect(response.body[0]).to.have.property('ID').and.to.eq(1)
        expect(response.body[0]).to.have.property('Name').and.to.eq('Gitlab')
        expect(response.body[0]).to.have.property('Endpoint')
        expect(response.body[0]).to.have.property('Auth')
        expect(response.body[0]).to.have.property('Proxy')
      })
  })

  it('provides github connection resources', () => {
    cy.request(`${Cypress.env('apiUrl')}plugins/github/connections`)
      .should((response) => {
        expect(response.status).to.eq(200)
        expect(response.headers).to.have.property('content-type').and.to.eq('application/json; charset=utf-8')
        expect(response.body).to.be.an('array')
        expect(response.body[0]).to.have.property('ID').and.to.eq(1)
        expect(response.body[0]).to.have.property('Name').and.to.eq('Github')
        expect(response.body[0]).to.have.property('Endpoint')
        expect(response.body[0]).to.have.property('Auth')
        expect(response.body[0]).to.have.property('Proxy')
        expect(response.body[0]).to.have.property('PrType')
        expect(response.body[0]).to.have.property('PrComponent')
        expect(response.body[0]).to.have.property('IssueSeverity')
        expect(response.body[0]).to.have.property('IssuePriority')
        expect(response.body[0]).to.have.property('IssueComponent')
        expect(response.body[0]).to.have.property('IssueTypeBug')
        expect(response.body[0]).to.have.property('IssueTypeIncident')
        expect(response.body[0]).to.have.property('IssueTypeRequirement')
      })
  })
})