require('module-alias/register')

const jira = require('@collectors/jira')
const { JiraIssue } = require('@db/postgres')

const { MongoRaw } = require('@db/mongodb/connection')

const main = async () => {
  // get data from Jira
  // get users
  // get issues
  const issues = await jira.issues.collectIssues('test-api')
  console.log(issues)
  // get changelogs

  // store data in mongodb
  const { mongoClient, mongoDb } = await MongoRaw.connect()
  const issuesCollection = MongoRaw.createCollection(mongodb, 'issues')
	// issuesCollection.insertMany([issues])
	// MongoRaw.disconnect(mongoClient)  

  // store data in postgress
  await JiraIssue.create({})

  // enhance data in postgres db
}

main()
