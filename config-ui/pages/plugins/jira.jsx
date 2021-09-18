import Head from 'next/head'
import { useState, useEffect } from 'react'
import styles from '../../styles/Home.module.css'
import { Tooltip, Position, FormGroup, InputGroup, Button, Label, Icon } from '@blueprintjs/core'
import dotenv from 'dotenv'
import path from 'path'
import * as fs from 'fs/promises'
import { existsSync } from 'fs'
import Nav from '../../components/Nav'
import Sidebar from '../../components/Sidebar'
import Content from '../../components/Content'
import SaveAlert from '../../components/SaveAlert'

export default function Home(props) {
  const { env } = props

  const [alertOpen, setAlertOpen] = useState(false)
  const [jiraEndpoint, setJiraEndpoint] = useState(env.JIRA_ENDPOINT)
  const [jiraBasicAuthEncoded, setJiraBasicAuthEncoded] = useState(env.JIRA_BASIC_AUTH_ENCODED)
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState(env.JIRA_ISSUE_EPIC_KEY_FIELD)
  const [jiraIssueTypeMapping, setJiraIssueTypeMapping] = useState(env.JIRA_ISSUE_TYPE_MAPPING)
  const [jiraIssueBugStatusMapping, setJiraIssueBugStatusMapping] = useState(env.JIRA_ISSUE_BUG_STATUS_MAPPING)
  const [jiraIssueIncidentStatusMapping, setJiraIssueIncidentStatusMapping] = useState(env.JIRA_ISSUE_INCIDENT_STATUS_MAPPING)
  const [jiraIssueStoryStatusMapping, setJiraIssueStoryStatusMapping] = useState(env.JIRA_ISSUE_STORY_STATUS_MAPPING)
  const [jiraIssueStoryCoefficient, setJiraIssueStoryCoefficient] = useState(env.JIRA_ISSUE_STORYPOINT_COEFFICIENT)
  const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState(env.JIRA_ISSUE_STORYPOINT_FIELD)
  const [jiraBoardGitlabeProjects, setJiraBoardGitlabeProjects] = useState(env.JIRA_BOARD_GITLAB_PROJECTS)

  function updateEnv(key, value) {
    fetch(`/api/setenv/${key}/${encodeURIComponent(value)}`)
  }

  function saveAll(e) {
    e.preventDefault()
    updateEnv('JIRA_ENDPOINT', jiraEndpoint)
    updateEnv('JIRA_BASIC_AUTH_ENCODED', jiraBasicAuthEncoded)
    updateEnv('JIRA_ISSUE_EPIC_KEY_FIELD', jiraIssueEpicKeyField)
    updateEnv('JIRA_ISSUE_TYPE_MAPPING', jiraIssueTypeMapping)
    updateEnv('JIRA_ISSUE_BUG_STATUS_MAPPING', jiraIssueBugStatusMapping)
    updateEnv('JIRA_ISSUE_INCIDENT_STATUS_MAPPING', jiraIssueIncidentStatusMapping)
    updateEnv('JIRA_ISSUE_STORY_STATUS_MAPPING', jiraIssueStoryStatusMapping)
    updateEnv('JIRA_ISSUE_STORYPOINT_COEFFICIENT', jiraIssueStoryCoefficient)
    updateEnv('JIRA_ISSUE_STORYPOINT_FIELD', jiraIssueStoryPointField)
    updateEnv('JIRA_BOARD_GITLAB_PROJECTS', jiraBoardGitlabeProjects)
    setAlertOpen(true)
  }

  return (
    <div className={styles.container}>

      <Head>
        <title>Devlake Config-UI</title>
        <meta name="description" content="Lake: Config" />
        <link rel="icon" href="/favicon.ico" />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link href="https://fonts.googleapis.com/css2?family=Source+Sans+Pro:wght@400;600&display=swap" rel="stylesheet" />
        <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500;600&display=swap" rel="stylesheet" />
      </Head>

      <Nav />
      <Sidebar />
      <Content>
        <main className={styles.main}>

          <form className={styles.form}>

            <div className={styles.headlineContainer}>
              <h2 className={styles.headline}>Jira Plugin</h2>
              <p className={styles.description}>Jira Account and config settings</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label=""
                inline={true}
                labelFor="jira-endpoint"
                helperText="JIRA_ENDPOINT"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <Label>
                  Endpoint&nbsp;URL <span className={styles.requiredStar}>*</span>
                  <InputGroup
                    id="jira-endpoint"
                    placeholder="Enter Jira endpoint eg. https://merico.atlassian.net/rest"
                    defaultValue={jiraEndpoint}
                    onChange={(e) => setJiraEndpoint(e.target.value)}
                    className={styles.input}
                  />
                </Label>
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="jira-basic-auth"
                helperText="JIRA_BASIC_AUTH_ENCODED"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <Label>
                  Basic&nbsp;Auth&nbsp;Token <span className={styles.requiredStar}>*</span>
                  <InputGroup
                    id="jira-basic-auth"
                    placeholder="Enter Jira Auth eg. EJrLG8DNeXADQcGOaaaX4B47"
                    defaultValue={jiraBasicAuthEncoded}
                    onChange={(e) => setJiraBasicAuthEncoded(e.target.value)}
                    className={styles.input}
                  />
                </Label>
              </FormGroup>
            </div>

            <div className={styles.headlineContainer}>
              <h3 className={styles.headline}>Status Mappings</h3>
              <p className={styles.description}>Map your custom Jira status to the correct values</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="jira-bug-status-mapping"
                helperText="JIRA_ISSUE_BUG_STATUS_MAPPING"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <Tooltip content="Map your custom bug status to Devlake standard status" position={Position.TOP}>
                  <Label>
                    Issue&nbsp;Bug
                    <InputGroup
                      id="jira-bug-status-mapping"
                      placeholder="<STANDARD_STATUS_1>:<ORIGIN_STATUS_1>,<ORIGIN_STATUS_2>;<STANDARD_STATUS_2>"
                      defaultValue={jiraIssueBugStatusMapping}
                      onChange={(e) => setJiraIssueBugStatusMapping(e.target.value)}
                      className={styles.input}
                    />
                  </Label>
                </Tooltip>
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="jira-incident-status-mapping"
                helperText="JIRA_ISSUE_INCIDENT_STATUS_MAPPING"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <Tooltip content="Map your custom incident status to Devlake standard status" position={Position.TOP}>
                  <Label>
                    Issue&nbsp;Incident
                    <InputGroup
                      id="jira-bug-status-mapping"
                      placeholder="<STANDARD_STATUS_1>:<YOUR_STATUS_1>,<YOUR_STATUS_2>;<STANDARD_STATUS_2>"
                      defaultValue={jiraIssueIncidentStatusMapping}
                      onChange={(e) => setJiraIssueIncidentStatusMapping(e.target.value)}
                      className={styles.input}
                    />
                  </Label>
                </Tooltip>
              </FormGroup>
            </div>

          <div className={styles.formContainer}>
            <FormGroup
              inline={true}
              labelFor="jira-story-status-mapping"
              helperText="JIRA_ISSUE_STORY_STATUS_MAPPING"
              className={styles.formGroup}
              contentClassName={styles.formGroup}
            >
              <Tooltip content="Map your custom story status to Devlake standard status" position={Position.TOP}>
                <Label>
                Issue&nbsp;Story
                <InputGroup
                  id="jira-story-status-mapping"
                  placeholder="<STANDARD_STATUS_1>:<YOUR_STATUS_1>,<YOUR_STATUS_2>;<STANDARD_STATUS_2>"
                  defaultValue={jiraIssueStoryStatusMapping}
                  onChange={(e) => setJiraIssueStoryStatusMapping(e.target.value)}
                  className={styles.input}
                />
                </Label>
              </Tooltip>
            </FormGroup>
          </div>

          <div className={styles.headlineContainer}>
            <h3 className={styles.headline}>Type Mappings</h3>
            <p className={styles.description}>Map your custom Jira issue type to the correct values</p>
          </div>

          <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="jira-issue-type-mapping"
                helperText="JIRA_ISSUE_TYPE_MAPPING"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <Tooltip content="Map your custom type to Devlake standard type" position={Position.TOP}>
                  <Label>
                  Issue&nbsp;Type
                  <InputGroup
                    id="jira-issue-type-mapping"
                    placeholder="STANDARD_TYPE_1:ORIGIN_TYPE_1,ORIGIN_TYPE_2;STANDARD_TYPE_2:...."
                    defaultValue={jiraIssueTypeMapping}
                    onChange={(e) => setJiraIssueTypeMapping(e.target.value)}
                    className={styles.input}
                  />
                  </Label>
                </Tooltip>
              </FormGroup>
          </div>

          <div className={styles.headlineContainer}>
            <h3 className={styles.headline}>Jira / Gitlab Connection</h3>
            <p className={styles.description}>Connect jira board to gitlab projects</p>
            </div>

            <div className={styles.formContainer}>
            <FormGroup
              inline={true}
              labelFor="jira-board-projects"
              helperText="JIRA_BOARD_GITLAB_PROJECTS"
              className={styles.formGroup}
              contentClassName={styles.formGroup}
            >
              <Tooltip content="Jira board and Gitlab projects relationship" position={Position.TOP}>
                <Label>
                  Jira&nbsp;Board&nbsp;Gitlab&nbsp;Projects
                  <InputGroup
                    id="jira-storypoint-field"
                    placeholder="<JIRA_BOARD>:<GITLAB_PROJECT_ID>,...; eg. 8:8967944,8967945;9:8967946,8967947"
                    defaultValue={jiraBoardGitlabeProjects}
                    onChange={(e) => setJiraBoardGitlabeProjects(e.target.value)}
                    className={styles.input}
                  />
                </Label>
              </Tooltip>
            </FormGroup>
          </div>

          <div className={styles.headlineContainer}>
            <h3 className={styles.headline}>Additional Customization Settings</h3>
            <p className={styles.description}>Additional Jira settings</p>
          </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="jira-epic-key"
                helperText="JIRA_ISSUE_EPIC_KEY_FIELD"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                  <Label>
                    Issue&nbsp;Epic&nbsp;Key&nbsp;Field

                    <div>
                      <Tooltip content="Get help with Issue Epic Key Field" position={Position.TOP}>
                        <a href="https://github.com/merico-dev/lake/tree/main/plugins/jira#set-jira-custom-fields"
                          target="_blank"
                          className={styles.helpIcon}>
                            <Icon icon="help" size={15} />
                        </a>
                      </Tooltip>
                    </div>

                    <Tooltip content="Your custom epic key field" position={Position.TOP}>
                      <InputGroup
                        id="jira-epic-key"
                        placeholder="Enter Jira epic key field"
                        defaultValue={jiraIssueEpicKeyField}
                        onChange={(e) => setJiraIssueEpicKeyField(e.target.value)}
                        className={styles.helperInput}
                      />
                    </Tooltip>
                  </Label>
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="jira-storypoint-field"
                helperText="JIRA_ISSUE_STORYPOINT_FIELD"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <Tooltip content="Your custom story point key field" position={Position.TOP}>
                  <Label>
                    Issue&nbsp;Storypoint&nbsp;Field

                    <div>
                      <Tooltip content="Get help with Issue Story Point Field" position={Position.TOP}>
                        <a href="https://github.com/merico-dev/lake/tree/main/plugins/jira#set-jira-custom-fields"
                          target="_blank"
                          className={styles.helpIcon}>
                            <Icon icon="help" size={15} />
                        </a>
                      </Tooltip>
                    </div>
                    <InputGroup
                      id="jira-storypoint-field"
                      placeholder="Enter Jira Story Point Field"
                      defaultValue={jiraIssueStoryPointField}
                      onChange={(e) => setJiraIssueStoryPointField(e.target.value)}
                      className={styles.input}
                    />
                  </Label>
                </Tooltip>
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="jira-storypoint-coef"
                helperText="JIRA_ISSUE_STORYPOINT_COEFFICIENT"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <Tooltip content="Your custom story point coefficent (optional)" position={Position.TOP}>
                  <Label>
                    Issue&nbsp;Storypoint&nbsp;Coefficient <span className={styles.requiredStar}>*</span>
                    <InputGroup
                      id="jira-storypoint-coef"
                      placeholder="Enter Jira Story Point Coefficient"
                      defaultValue={jiraIssueStoryCoefficient}
                      onChange={(e) => setJiraIssueStoryCoefficient(e.target.value)}
                      className={styles.input}
                    />
                  </Label>
                </Tooltip>
              </FormGroup>
            </div>

            <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
            <Button type="submit" outlined={true} large={true} className={styles.saveBtn} onClick={saveAll}>Save Config</Button>
          </form>
        </main>
      </Content>
    </div>
  )
}

export async function getStaticProps() {

  const filePath = process.env.ENV_FILEPATH || path.join(process.cwd(), 'data', '../../.env')
  const exist = existsSync(filePath);
  if (!exist) {
    return {
      props: {
        env: {},
      }
    }
  }
  const fileData = await fs.readFile(filePath)
  const env = dotenv.parse(fileData)

  return {
    props: {
      env
    },
  }
}
