require('module-alias/register')

const { JiraUser } = require('@db/postgres')

let main = async ()=>{
  // get data from Jira
    // get users
    // get issues
    // get changelogs

  // store all raw data in mongodb

  // store data in postgress
  let jiraUser = await JiraUser.create({})

  // enhance data in postgres db

}

main()