#!/usr/bin/env node

require('module-alias/register')
const _has = require('lodash/has')

const jira = require('./jira')
const consumer = require('../queue/consumer')

const queue = 'enrichment'

const jobHandler = async (job) => {
  if (_has(job, 'jira')) {
    await jira.enrich(job.jira)
  }
}

consumer(queue, jobHandler)
