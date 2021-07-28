require('module-alias/register')

const issueCollector = require('../collector/issues')
const changelogCollector = require('../collector/changelogs')
const constants = require('@config/constants.json').jira

module.exports = {
  async enrich (rawDb, enrichedDb, projectId) {
    console.log('INFO: Starting Jira Enrichment for projectId: ', projectId)
    await module.exports.enrichLeadTimeOnIssues(
      rawDb,
      enrichedDb,
      projectId
    )
    console.log('INFO: Done enriching Jira issues')
  },

  /**
   * 
   * @param {String} jiraIssueTypeString 
   * @returns key
   * 
   * Sometimes, users in Jira can set their own labels for issue statuses like 'Done'
   * This method allows users to set their own labels from their Jira system in the /config/constants.json file.
   */
  mapIssueTypeFromConfiguration (jiraIssueTypeString, issueTypes) {
    if(!jiraIssueTypeString || jiraIssueTypeString === ''){
      return ''
    } else {
      for (const key in issueTypes) {
        if (Object.hasOwnProperty.call(issueTypes, key)) {
          const element = issueTypes[key];
          if(element.toLowerCase() === jiraIssueTypeString.toLowerCase()){
            return key
          }
        }
      }
      // If no mapping is found, return the original value from the Jira API
      return jiraIssueTypeString
    }
  },

  async enrichLeadTimeOnIssues (rawDb, enrichedDb, projectId) {
    const { JiraIssue } = enrichedDb

    const issues = await issueCollector.findIssues({
      'fields.project.id': `${projectId}`
    }, rawDb)

    const creationPromises = []
    const leadTimePromises = []
    const issuesToCreate = []
    issues.forEach(async issue => {
      leadTimePromises.push(module.exports.calculateLeadTime(issue, rawDb))
      issuesToCreate.push({
        id: issue?.id,
        url: issue?.self,
        title: issue?.fields?.summary,
        projectId: issue?.fields?.project?.id,
        issueType: module.exports.mapIssueTypeFromConfiguration(issue?.fields?.issueType?.name, constants.issueTypes)
        // description: issue.fields.description
      })
    })

    const leadTimes = await Promise.all(leadTimePromises)

    leadTimes.forEach((leadTime, index) => {
      let issue = issuesToCreate[index]
      console.log('INFO: issueId & leadTime', issue.id, leadTime)
      issue = {
        leadTime,
        ...issue
      }
      creationPromises.push(JiraIssue.findOrCreate({
        where: {
          id: issue.id
        },
        defaults: issue
      }))
    })

    await Promise.all(creationPromises)
  },

  async calculateLeadTime (issue, db) {
    const changelogs = await changelogCollector.findChangelogs({
      issueId: `${issue.id}`
    }, db)

    let leadTime = 0
    let lastTime = new Date(issue.fields.created).getTime()
    let isDone = false

    for (const change of changelogs) {
      for (const item of change.items) {
        if (item.field === 'status') {
          const changeTime = new Date(change.created).getTime()

          if (!constants.closedStatuses.includes(item.fromString)) {
            const elapsedTime = changeTime - lastTime

            leadTime += elapsedTime
          }

          lastTime = changeTime
          isDone = closedStatuses.includes(item.toString)
        }
      }
    }

    return isDone
      ? Math.round(leadTime / 1000)
      : 0
  }
}
