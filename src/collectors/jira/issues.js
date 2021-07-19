const axios = require('axios')
const config = require('@config/resolveConfig').jira

module.exports = {

  async collectIssues(project) {
    try {

      console.log(config.jira)

      const response = await axios.get(`${config.host}/rest/api/3/search?jql=project="${project}"`, {
        headers: {
          'Accept': 'application/json',
          'Authorization': `Basic ${config.basicAuth}`
        }
      })

      return response.data.issues

    } catch (error) {
      console.error(error)
    }
  }
}
