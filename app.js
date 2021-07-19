require('module-alias/register')

const { JiraUser } = require('@db/postgres')
const jira = require('@collectors/jira')

let main = async ()=>{
  // get data from Jira
    // get users
    // get issues
    const issues = await jira.issues.collectIssues('test-api')
    console.log(issues) 
    // get changelogs

  // store data in postgress
  let jiraUser = await JiraUser.create({})

  // enhance data in postgres db

}

main()
