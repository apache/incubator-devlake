const prompt = require('prompt-sync')({ sigint: true })
const fs = require('fs')
const path = require('path')
const readline = require('readline')

const groupConfig = (configArr, isSample) => {
  const arr = []
  configArr.forEach(name => {
    arr.push(path.join(__dirname, `../config/${name}${isSample ? '.sample' : ''}.js`))
  })
  return arr
}

// Main Config Vars
const confirmationString = '(y/n) '
const configNames = ['docker', 'local', 'plugins']
const allConfig = groupConfig(configNames, false)
const allConfigSample = groupConfig(configNames, true)

module.exports = {
  welcomePrompt () {
    console.log('*******************************************************')
    console.log(
    `██       █████  ██   ██ ███████
██      ██   ██ ██  ██  ██
██      ███████ █████   █████
██      ██   ██ ██  ██  ██
███████ ██   ██ ██   ██ ███████`
    )

    console.log('')
    console.log('... by Merico (https://meri.co)')
    console.log('*******************************************************')
  },

  checkAndCreateConfig () {
    const createConfig = prompt(`➤➤➤ Would you like us to create config files for you? ${confirmationString}`)

    if (createConfig === 'y') {
      allConfig.forEach((path, i) => {
        if (fs.existsSync(path)) {
          // We have a config file, skip
          console.log(`Config detected for ${path}`)
        } else {
          // We need to create a new config file from example
          console.log(`➤➤➤ Creating new Sample ${allConfig[i]}`)
          fs.copyFile(allConfigSample[i], allConfig[i], err => {
            if (err) throw err
          })
        }
      })
    } else {
      process.exit(1)
    }
  },

  openReadSyncLocal (pathSample, path) {
    const readInterface = readline.createInterface({
      input: fs.createReadStream(pathSample),
      output: false,
      console: false
    })

    const writeStream = fs.createWriteStream(path)

    // Collect replacement vars in local.js
    console.log('')
    console.log('\x1b[36m%s\x1b[0m',
    `❕ TIP: You can read more on how to get your jira board id here:
    \nhttps://github.com/merico-dev/lake/tree/main/src/plugins/jira-pond#find-board-id \n`)
    const jiraBoardId = prompt('➤➤➤ What is your jira board id?  ')

    console.log('')
    console.log('\x1b[36m%s\x1b[0m',
    `❕ TIP: You can read more on how to get your gitlab project id here:
    \nhttps://github.com/merico-dev/lake/tree/main/src/plugins/gitlab-pond#finding-project-id \n`)
    const gitlabProjectId = prompt('➤➤➤ What is your gitlab project id?  ')

    // Replace lines in local.js
    readInterface.on('line', (line) => {
      // Jira board ID
      if (line.match('"<your-board-id>"')) {
        writeStream.write(line.replace('"<your-board-id>"', jiraBoardId) + '\n')
      } else if (line.match('"<your-gitlab-project-id>"')) {
        writeStream.write(line.replace('"<your-gitlab-project-id>"', gitlabProjectId) + '\n')
      } else {
        writeStream.write(line + '\n')
      }
    })
  },

  openReadSyncPlugins (pathSample, path) {
    const readInterface = readline.createInterface({
      input: fs.createReadStream(pathSample),
      output: false,
      console: false
    })

    const writeStream = fs.createWriteStream(path)

    // Collect replacement vars in plugins.js
    console.log('')
    console.log('\x1b[36m%s\x1b[0m',
    `❕ TIP: You can read more on how to get jira token here:
    \nhttps://github.com/merico-dev/lake/tree/main/src/plugins/jira-pond#generating-api-token \n`)
    const jiraToken = prompt('➤➤➤ What is your jira token?  ')

    console.log('\x1b[36m%s\x1b[0m',
      '❕ TIP: This is the email you use to login to Jira')
    const jiraEmail = prompt('➤➤➤ What is your jira user email?  ')

    console.log(
      '\x1b[36m%s\x1b[0m',
      '❕ TIP: This is the base url for jira that you use. IE: for this url: https://merico.atlassian.net/secure/RapidBoard.jspa?rapidView=8&projectKey=EE, you would use https://merico.atlassian.net'
    )
    const jiraHost = prompt('➤➤➤ What is your jira host url?  ')

    console.log('')
    console.log('\x1b[36m%s\x1b[0m',
    `❕ TIP: You can read more on how to get your github token here:
    \nhttps://github.com/merico-dev/lake/tree/main/src/plugins/gitlab-pond#create-a-gitlab-api-token \n`)
    const gitlabToken = prompt('➤➤➤ What is your gitlab token?  ')

    // Replace lines in plugins.js
    readInterface.on('line', (line) => {
      // Jira board ID
      if (line.match('"<your-jira-token>"')) {
        writeStream.write(line.replace('<your-jira-token>', jiraToken) + '\n')
      } else if (line.match('"<your-gitlab-token>"')) {
        writeStream.write(line.replace('<your-gitlab-token>', gitlabToken) + '\n')
      } else if (line.match('"<your-jira-email>"')) {
        writeStream.write(line.replace('<your-jira-email>', jiraEmail) + '\n')
      } else if (line.match('"<your-jira-host>"')) {
        writeStream.write(line.replace('<your-jira-host>', jiraHost) + '\n')
      } else {
        writeStream.write(line + '\n\n')
      }
    })
  },

  main () {
    module.exports.welcomePrompt()
    module.exports.checkAndCreateConfig()
    module.exports.openReadSyncLocal(allConfigSample[1], allConfig[1]) // local.js
    module.exports.openReadSyncPlugins(allConfigSample[2], allConfig[2]) // plugins.js
  }
}

module.exports.main()
