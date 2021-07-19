require('module-alias/register')

const { JiraUser } = require('@db/postgres')
const jira = require('@collectors/jira')

let main = async ()=>{
  // get data from Jira
  const issues = await jira.issues.collectIssues('test-api')
  console.log(issues)

  // store data in mongodb

  // store data in postgress
  let jiraUser = await JiraUser.create({})

  // enhance data in postgres db

}

main()
