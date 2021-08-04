const fetcher = require('./fetcher')
const axios = require('axios')
const sinon = require('sinon')
const chai = require('chai')
const sinonChai = require('sinon-chai')
const {
  expect
} = chai
chai.use(sinonChai)

describe('fetcher', () => {
  afterEach(() => {
    sinon.restore()
  })

  describe('fetchPaged', () => {
    it('if only one page, only fetch once', async () => {
      sinon.stub(axios, 'get').resolves({
        data: [],
        headers: {
          'x-next-page': null
        }
      })

      fetcher.fetchPaged('https://fake.com/projects')

      for await (const item of fetcher.fetchPaged('https://fake.com/projects')) {
        return null
      }
      expect(axios.get).to.have.been.calledOnce
    })

    it('if more than one page, only fetch more than once', async () => {
      const stub = sinon.stub(axios, 'get')

      stub.onCall(0).resolves({
        data: [],
        headers: {
          'x-next-page': 1
        }
      })

      stub.onCall(1).resolves({
        data: [],
        headers: {
          'x-next-page': null
        }
      })

      for await (const item of fetcher.fetchPaged('https://fake.com/projects')) {
        return null
      }
      expect(axios.get).to.have.been.calledTwice
    })
  })
})
