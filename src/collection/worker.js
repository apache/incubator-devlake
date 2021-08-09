#!/usr/bin/env node

require('module-alias/register')
const axios = require('axios')
const _has = require('lodash/has')

const dbConnector = require('@mongo/connection')
const { collection } = require('../plugins')
const consumer = require('../queue/consumer')
const config = require('@config/resolveConfig').lake || {}

const queue = 'collection'

const jobHandler = async (job) => {
  console.log('INFO: Collection worker received job: ', job)
  // const {
  //   db, client
  // } = await dbConnector.connect()

  // try {
  //   await Promise.all(
  //     Object.keys(job)
  //       .filter(key => _has(collection, key))
  //       .map(pluginName => collection[pluginName](db, job[pluginName]))
  //   )
  // } catch (error) {
  //   console.log('Failed to collect', error)
  // } finally {
  //   dbConnector.disconnect(client)
  // }

  // await axios.post(
  //   `http://localhost:${process.env.ENRICHMENT_PORT || 3000}`,
  //   job,
  //   { headers: { 'x-token': config.token || '' } }
  // )
}

consumer(queue, jobHandler)
