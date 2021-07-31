#!/usr/bin/env node

require('module-alias/register')
const _has = require('lodash/has')

const dbConnector = require('@mongo/connection')
const { buildPluginRegistry } = require('../plugins')
const consumer = require('../queue/consumer')
const enrichedDb = require('@db/postgres')

const queue = 'enrichment'

const jobHandler = async (job) => {
  const enrichment = await buildPluginRegistry('enricher')
  const {
    db: rawDb, client
  } = await dbConnector.connect()

  try {
    if (_has(job, 'jira')) {
      await enrichment.plugins[job.jira.enricher](rawDb, enrichedDb, job.jira)
    }
    if (_has(job, 'gitlab')) {
      await enrichment.plugins[job.gitlab.enricher](rawDb, enrichedDb, job.gitlab)
    }
  } catch (error) {
    console.log('Failed to enrich', error)
  } finally {
    dbConnector.disconnect(client)
  }
}

consumer(queue, jobHandler)
