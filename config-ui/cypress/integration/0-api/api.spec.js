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
        expect(response.body[0]).to.have.property('id').and.to.eq(1)
        expect(response.body[0]).to.have.property('name').and.to.eq('Jenkins')
        expect(response.body[0]).to.have.property('endpoint')
        expect(response.body[0]).to.have.property('username')
        expect(response.body[0]).to.have.property('password')
        expect(response.body[0]).to.have.property('proxy')
      })
  })

  it('provides gitlab connection resources', () => {
    cy.request(`${Cypress.env('apiUrl')}plugins/gitlab/connections`)
      .should((response) => {
        expect(response.status).to.eq(200)
        expect(response.headers).to.have.property('content-type').and.to.eq('application/json; charset=utf-8')
        expect(response.body).to.be.an('array')
        expect(response.body[0]).to.have.property('id').and.to.eq(1)
        expect(response.body[0]).to.have.property('name').and.to.eq('Gitlab')
        expect(response.body[0]).to.have.property('endpoint')
        expect(response.body[0]).to.have.property('auth')
        expect(response.body[0]).to.have.property('proxy')
      })
  })

  it('provides github connection resources', () => {
    cy.request(`${Cypress.env('apiUrl')}plugins/github/connections`)
      .should((response) => {
        expect(response.status).to.eq(200)
        expect(response.headers).to.have.property('content-type').and.to.eq('application/json; charset=utf-8')
        expect(response.body).to.be.an('array')
        expect(response.body[0]).to.have.property('id').and.to.eq(1)
        expect(response.body[0]).to.have.property('name').and.to.eq('Github')
        expect(response.body[0]).to.have.property('endpoint')
        expect(response.body[0]).to.have.property('auth')
        expect(response.body[0]).to.have.property('proxy')
        expect(response.body[0]).to.have.property('prType')
        expect(response.body[0]).to.have.property('prComponent')
        expect(response.body[0]).to.have.property('issueSeverity')
        expect(response.body[0]).to.have.property('issuePriority')
        expect(response.body[0]).to.have.property('issueComponent')
        expect(response.body[0]).to.have.property('issueTypeBug')
        expect(response.body[0]).to.have.property('issueTypeIncident')
        expect(response.body[0]).to.have.property('issueTypeRequirement')
      })
  })
})