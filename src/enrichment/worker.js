#!/usr/bin/env node

require('module-alias/register')
const _has = require('lodash/has')

const dbConnector = require('@mongo/connection')
const { enrichment } = require('../plugins')
const consumer = require('../queue/consumer')
const enrichedDb = require('@db/postgres')

const queue = 'enrichment'

const jobHandler = async (job) => {
  const {
    db: rawDb, client
  } = await dbConnector.connect()

  console.log("INFO >>> recieve enriche job")
  try {
    await Promise.all(
      Object.keys(job)
        .filter(key => _has(enrichment, key))
        .map(pluginName => enrichment[pluginName](rawDb, enrichedDb, job[pluginName]))
    )
  } catch (error) {
    console.log('Failed to enrich', error)
  } finally {
    dbConnector.disconnect(client)
  }
}

consumer(queue, jobHandler)
