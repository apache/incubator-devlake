#!/usr/bin/env node

const amqp = require('amqplib/callback_api')
const _has = require('lodash/has')

const jira = require('./jira')

const amqpUrl = 'amqp://guest:guest@localhost:5672/rabbitmq'

amqp.connect(amqpUrl, function (error0, connection) {
  if (error0) {
    throw error0
  }
  connection.createChannel(function (error1, channel) {
    if (error1) {
      throw error1
    }
    const queue = 'collection'

    channel.assertQueue(queue)

    console.log(' [*] Waiting for messages in %s. To exit press CTRL+C', queue)

    channel.consume(queue, async function (msg) {
      const job = JSON.parse(msg.content.toString())

      if (_has(job, 'jira')) {
        await jira.collect(job.jira)
      }
    }, {
      noAck: true
    })
  })
})
