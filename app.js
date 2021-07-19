require('module-alias/register')

const { JiraUser } = require('@db/postgres')


let main = async ()=>{
  let jiraUser = await JiraUser.create({})

}

main()