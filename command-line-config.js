const prompt = require("prompt-sync")({ sigint: true });
var fs = require('fs');
const { dirname } = require("path");

const confirmationString = '(y/n)'

module.exports = {
  async replaceInFile(path, findRegex, replacementString){
    // console.log("JON >>> replaceInFile", path, findRegex, replacementString);
    try {
      await fs.readFile(path, 'utf8', async function (err,data) {
        if (err) {
          return console.log(err);
        }
        var result = data.replace(findRegex, replacementString);
  
        await fs.writeFile(path, result, "utf8", async function (err) {
          if (err) return console.log(err);
        });
      });
    } catch (error) {
      console.error('ERROR: could not replace in file', error)  
      throw error    
    }
  },

  copyFile(path, newPath){
    try {
      fs.copyFile(path, newPath, (err) => {
        if (err) throw err;
        return
      });
    } catch (error) {
      console.error('ERROR: could not copy file: ', error)      
    }
  },

  copyConfigFilesIfNotExists(){
    module.exports.copyFile(
      `${__dirname}/config/local.sample.js`,
      `${__dirname}/config/local.js`
    );
    module.exports.copyFile(
      `${__dirname}/config/plugins.sample.js`,
      `${__dirname}/config/plugins.js`
    );
    module.exports.copyFile(
      `${__dirname}/config/docker.sample.js`,
      `${__dirname}/config/docker.js`
    );
  },

  welcomePrompt(){
    console.log("*******************************************************");
    console.log(
    `██       █████  ██   ██ ███████
██      ██   ██ ██  ██  ██ 
██      ███████ █████   █████
██      ██   ██ ██  ██  ██   
███████ ██   ██ ██   ██ ███████`
    );

    console.log("... by Merico (https://meri.co/)");
    console.log("*******************************************************");
  },

  async main () {
    module.exports.welcomePrompt()
    
    let createConfig = prompt(`➤➤➤ Would you like us to create config files for you? ${confirmationString}`);
    if (createConfig === 'y') { module.exports.copyConfigFilesIfNotExists();}
    
    console.log('\x1b[36m%s\x1b[0m',
      "❕ TIP: You can read more on how to get your jira board id here: https://github.com/merico-dev/lake/tree/main/src/plugins/jira-pond#find-board-id"
    );
    let jiraBoardid = prompt(`➤➤➤ What is your jira board id? `);
    await module.exports.setJiraBoardId(jiraBoardid);

    console.log('\x1b[36m%s\x1b[0m', 
      "❕ TIP: You can read more on how to get your github project id here: https://github.com/merico-dev/lake/tree/main/src/plugins/gitlab-pond#finding-project-id"
    );
    let githubProjectId = prompt(`➤➤➤ What is your github project id? `);
    await module.exports.setGitlabProjectId(githubProjectId);
    
    // console.log(
    //   "\x1b[36m%s\x1b[0m",
    //   "❕ TIP: You can read more on how to get jira token here: https://github.com/merico-dev/lake/tree/main/src/plugins/jira-pond#generating-api-token"
    // );
    // let jiraToken = prompt(`➤➤➤ What is your jira token? `);
    // module.exports.setJiraToken(jiraToken);
    
    // console.log(
    //   "\x1b[36m%s\x1b[0m",
    //   "❕ TIP: You can read more on how to get your github token here: https://github.com/merico-dev/lake/tree/main/src/plugins/gitlab-pond#create-a-gitlab-api-token"
    // );
    // let gitlabToken = prompt(`➤➤➤ What is your github token? `);
    // module.exports.setGitlabToken(gitlabToken);
  },

  async setJiraBoardId(value) {
    await module.exports.replaceInFile(
      `${__dirname}/config/local.js`,
      /"<your-board-id>"/i,
      value
    );
  },

  async setGitlabProjectId(value) {
    await module.exports.replaceInFile(
      `${__dirname}/config/local.js`,
      /"<your-gitlab-project-id>"/i,
      value
    );
  },

  async setJiraToken(value) {
    await module.exports.replaceInFile(
      `${__dirname}/config/plugins.js`,
      /<your-jira-token>/i,
      value
    );
  },

  async setGitlabToken(value) {
    await module.exports.replaceInFile(
      `${__dirname}/config/plugins.js`,
      /<your-gitlab-token>/i,
      value
    );
  },
};

module.exports.main()


