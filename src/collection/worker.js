#!/usr/bin/env node

require('module-alias/register')
const axios = require('axios')
const _has = require('lodash/has')

const jira = require('../collectors/jira')
const consumer = require('../queue/consumer')
const enrichmentApiUrl = require('@config/resolveConfig').enrichment.connectionString

const queue = 'collection'

const jobHandler = async (job) => {
  if (_has(job, 'jira')) {
    await jira.collect(job.jira)
  }

  await axios.post(enrichmentApiUrl, job)
}

consumer(queue, jobHandler)
