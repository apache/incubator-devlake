# Jira Pond

## Metrics

Metric Name | Description
:------------ | :-------------
Requirement Count	| Number of issues with type "Requirement"
Requirement Lead Time	| Lead time of issues with type "Requirement"
Requirement Delivery Rate |	Ratio of delivered requirements to all requirements
Bug Count	| Number of issues with type "Bug"<br><i>bugs are found during testing</i>
Bug Age	| Lead time of issues with type "Bug"<br><i>both new and deleted lines count</i>
Bugs Count per 1k Lines of Code |	Amount of bugs per 1000 lines of code
Incident Count | Number of issues with type "Incident"<br><i>incidents are found when running in production</i>
Incident Age | Lead time of issues with type "Incident"
Incident Count per 1k Lines of Code | Amount of incidents per 1000 lines of code

## Find Board Id
1. Navigate to the Jira board in the browser
2. in the URL bar, get the board id from the parameter `?rapidView=`

**Example:**
`https://<your_jira_url>/secure/RapidBoard.jspa?rapidView=51`

## Generating API token
1. Once logged into Jira, visit the url `https://id.atlassian.com/manage-profile/security/api-tokens`
2. Click the **Create API Token** button, and give it any label name
3. Copy and save the API token somewhere
4. In a terminal run the following command, with **user email** and **API token** string

    `echo -n user@example.com:api_token_string | base64`
5. Copy the encoded API token string into the `lake` plugin config file `config/plugins.js`

## Jira Specific String Configuration

Adjust what is considered "Bug", "Incident" or "Requirement". This can be modified in `config/plugins.js`.



```js
{
  package: 'jira-pond',
  name: 'jira',
  configuration: {
    enrichment: {
      issue: {
        mapping: {
          // This maps issue types in your Jira system to the standard issue type in dev lake
          // In lake, we define bugs as issues found in development process whereas
          // incidents are issues found in production environment
          // Format: <Standard Type>: [<Jira Type>]
          type: {
            // This mapping powers the metrics like Bug Count, But Age, and etc
            // Replace 'Bug' with your own issue types for bugs.
            Bug: ['Bug'],
            // This mapping powers the metrics like Incident Count, Incident Age, and etc
            // Replace 'Incident' with your own issue types for incidents
            Incident: ['Incident']
          }
        },
        // Enables lake to track which epic an issue belongs to
        // Replace 'customfiled_10014' with your own field ID for the epic key
        epicKeyField: 'customfield_10014'
      }
    }
  }
}
```
