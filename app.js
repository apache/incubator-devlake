require('module-alias/register')

const jira = require('@collectors/jira')
const { JiraIssue } = require('@db/postgres')
const { MongoClient } = require('mongodb')
const connection = require('@config/resolveConfig').mongo.connectionString
const client = new MongoClient(connection)

const main = async () => {
  // get data from Jira
  // get users

  // get issues and store in mongodb
  try {
    const issues = await jira.issues.collectIssues('test-api')
    // console.log('Issues from Jira API: ', issues)

    // store data in mongodb
    await client.connect()

    let issueCollection
    const foundCollectionsCursor = await client.db().listCollections()
    const foundCollections = await foundCollectionsCursor.toArray()
    const collectionName = 'jira_issues'

    // check if Jira collection exists
    const collectionExists = foundCollections
      .some(collection => collection.name === collectionName)

    if (collectionExists === true) {
      issueCollection = client.db().collection(collectionName)
    } else {
      issueCollection = client.db().createCollection(collectionName)
    }

    // Insert issues into mongodb
    await issueCollection.insertMany(issues)
    const foundIssuesCursor = await issueCollection.find()
    const foundIssues = await foundIssuesCursor.toArray()

    // Insert data in postgress
    foundIssues.forEach(async issue => {
      await JiraIssue.create({
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id,
        description: issue.fields.description

        // TODO: additional jira issue fields
        // leadTime: issue.fields.timespent
      })
    })
  } catch (e) {
    console.error(e)
  } finally {
    await client.close()
  }

  // enhance data in postgres db
}

main()
