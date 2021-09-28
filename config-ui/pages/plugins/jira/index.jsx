import Head from 'next/head'
import { useState, useEffect } from 'react'
import styles from '../../../styles/Home.module.css'
import {
  Tooltip, Position, FormGroup, InputGroup, Button, Label, Icon, Classes, Tab, Tabs, Overlay, Dialog
} from '@blueprintjs/core'
import dotenv from 'dotenv'
import path from 'path'
import * as fs from 'fs/promises'
import { existsSync } from 'fs'
import { findStrBetween } from '../../../utils/findStrBetween'
import { readAndSet } from '../../../utils/readAndSet'
import Nav from '../../../components/Nav'
import Sidebar from '../../../components/Sidebar'
import Content from '../../../components/Content'
import SaveAlert from '../../../components/SaveAlert'
import MappingTag from './MappingTag'
import MappingTagStatus from './MappingTagStatus'
import ClearButton from './ClearButton'

export default function Home(props) {
  const { env } = props

  const [alertOpen, setAlertOpen] = useState(false)
  const [jiraEndpoint, setJiraEndpoint] = useState(env.JIRA_ENDPOINT)
  const [jiraBasicAuthEncoded, setJiraBasicAuthEncoded] = useState(env.JIRA_BASIC_AUTH_ENCODED)
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState(env.JIRA_ISSUE_EPIC_KEY_FIELD)
  const [jiraIssueStoryCoefficient, setJiraIssueStoryCoefficient] = useState(env.JIRA_ISSUE_STORYPOINT_COEFFICIENT)
  const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState(env.JIRA_ISSUE_STORYPOINT_FIELD)
  const [jiraBoardGitlabeProjects, setJiraBoardGitlabeProjects] = useState(env.JIRA_BOARD_GITLAB_PROJECTS)

  // Type mappings state
  const [typeMappingBug, setTypeMappingBug] = useState([])
  const [typeMappingIncident, setTypeMappingIncident] = useState([])
  const [typeMappingRequirement, setTypeMappingRequirement] = useState([])
  const [typeMappingAll, setTypeMappingAll] = useState()

  // Status mappings state
  const [customStatusOverlay, setCustomStatusOverlay] = useState(false)
  const [statusTabId, setStatusTabId] = useState(0)
  const [statusMappingReqBug, setStatusMappingReqBug] = useState([])
  const [statusMappingResBug, setStatusMappingResBug] = useState([])
  const [statusMappingReqIncident, setStatusMappingReqIncident] = useState([])
  const [statusMappingResIncident, setStatusMappingResIncident] = useState([])
  const [statusMappingReqStory, setStatusMappingReqStory] = useState([])
  const [statusMappingResStory, setStatusMappingResStory] = useState([])
  const [customStatus, setCustomStatus] = useState([])
  const [customStatusName, setCustomStatusName] = useState('')


  function updateEnv(key, value) {
    fetch(`/api/setenv/${key}/${encodeURIComponent(value)}`)
  }

  function saveAll(e) {
    e.preventDefault()
    updateEnv('JIRA_ENDPOINT', jiraEndpoint)
    updateEnv('JIRA_BASIC_AUTH_ENCODED', jiraBasicAuthEncoded)
    updateEnv('JIRA_ISSUE_EPIC_KEY_FIELD', jiraIssueEpicKeyField)
    updateEnv('JIRA_ISSUE_TYPE_MAPPING', typeMappingAll)
    updateEnv('JIRA_ISSUE_BUG_STATUS_MAPPING', `Requirement:${statusMappingReqBug};Resolved:${statusMappingResBug};`)
    updateEnv('JIRA_ISSUE_INCIDENT_STATUS_MAPPING', `Requirement:${statusMappingReqIncident};Resolved:${statusMappingResIncident};`)
    updateEnv('JIRA_ISSUE_STORY_STATUS_MAPPING', `Requirement:${statusMappingReqStory};Resolved:${statusMappingResStory};`)
    updateEnv('JIRA_ISSUE_STORYPOINT_COEFFICIENT', jiraIssueStoryCoefficient)
    updateEnv('JIRA_ISSUE_STORYPOINT_FIELD', jiraIssueStoryPointField)
    updateEnv('JIRA_BOARD_GITLAB_PROJECTS', jiraBoardGitlabeProjects)

    // Save all custom status data
    customStatus.map(status => {
      const requirement = status.reqValue.toString()
      const resolved = status.resValue.toString()
      const name = `JIRA_ISSUE_${status.name.toUpperCase()}_STATUS_MAPPING`
      updateEnv(name, `Requirement:${requirement};Resolved:${resolved};`)
    })

    setAlertOpen(true)
  }

  useEffect(() => {
    const typeBug = 'Bug:' + typeMappingBug.toString() + ';'
    const typeIncident = 'Incident:' + typeMappingIncident.toString() + ';'
    const typeRequirement = 'Requirement:' + typeMappingRequirement.toString() + ';'
    const all = typeBug + typeIncident + typeRequirement
    setTypeMappingAll(all)
  }, [typeMappingBug, typeMappingIncident, typeMappingRequirement])

  useEffect(() => {
    // Load type & status mappings
    const envStr = [
      env.JIRA_ISSUE_TYPE_MAPPING,
      env.JIRA_ISSUE_BUG_STATUS_MAPPING,
      env.JIRA_ISSUE_INCIDENT_STATUS_MAPPING,
      env.JIRA_ISSUE_STORY_STATUS_MAPPING
    ]
    const fields = [
      {
        tagName: 'Bug:', tagLen: 4, isStatus: false, str: envStr[0],
        fn1: (arr) => setTypeMappingBug(arr)
      },
      {
        tagName: 'Incident:', tagLen: 9, isStatus: false, str: envStr[0],
        fn1: (arr) => setTypeMappingIncident(arr)
      },
      {
        tagName: 'Requirement:', tagLen: 12, isStatus: false, str: envStr[0],
        fn1: (arr) => setTypeMappingRequirement(arr)
      },
      {
        tagName: 'Bug:', tagLen: null, isStatus: true, str: envStr[1],
        fn1: (arr) => setStatusMappingReqBug(arr),
        fn2: (arr) => setStatusMappingResBug(arr)
      },
      {
        tagName: 'Incident:', tagLen: null, isStatus: true, str: envStr[2],
        fn1: (arr) => setStatusMappingReqIncident(arr),
        fn2: (arr) => setStatusMappingResIncident(arr)
      },
      {
        tagName: 'Story:', tagLen: null, isStatus: true, str: envStr[3],
        fn1: (arr) => setStatusMappingReqStory(arr),
        fn2: (arr) => setStatusMappingResStory(arr)
      },
    ]

    fields.map(field => {
      readAndSet(field.tagName, field.tagLen, field.isStatus, field.str, field.fn1, field.fn2)
    })

    //Load custom status mappings
    for (const field in env) {
      const bug = 'JIRA_ISSUE_BUG'
      const incident = 'JIRA_ISSUE_INCIDENT'
      const story = 'JIRA_ISSUE_STORY'
      const isStatusMapping = field.includes('_STATUS_MAPPING')
      const isNotDefault = (!field.includes(bug)) && (!field.includes(incident)) && (!field.includes(story))

      if (isStatusMapping && isNotDefault) {
        const strName = field.slice(11, -15)
        const strValuesReq = findStrBetween(env[field], 'Requirement:', ';')
        const strValuesRes = findStrBetween(env[field], 'Resolved:', ';')
        const req = strValuesReq[0].slice(12, -1).split(',')
        const res = strValuesRes[0].slice(9, -1).split(',')

        setCustomStatus(customStatus => [...customStatus, {name: strName, reqValue: req || '', resValue: res || ''}])
      }
    }
  }, [])

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
              <h3 className={styles.headline}>Issue Type Mappings</h3>
              <p className={styles.description}>Map your own issue types to Dev Lake's standard types</p>
            </div>

            <MappingTag
              labelName="Bug"
              labelIntent="danger"
              typeOrStatus="type"
              values={typeMappingBug}
              helperText="JIRA_ISSUE_TYPE_MAPPING"
              rightElement={<ClearButton onClick={() => setTypeMappingBug([])} />}
              onChange={(values) => setTypeMappingBug(values)}
            />

            <MappingTag
              labelName="Incident"
              labelIntent="warning"
              typeOrStatus="type"
              values={typeMappingIncident}
              helperText="JIRA_ISSUE_TYPE_MAPPING"
              rightElement={<ClearButton onClick={() => setTypeMappingIncident([])} />}
              onChange={(values) => setTypeMappingIncident(values)}
            />

            <MappingTag
              labelName="Requirement"
              labelIntent="primary"
              typeOrStatus="type"
              values={typeMappingRequirement}
              helperText="JIRA_ISSUE_TYPE_MAPPING"
              rightElement={<ClearButton onClick={() => setTypeMappingRequirement([])} />}
              onChange={(values) => setTypeMappingRequirement(values)}
            />

            <div className={styles.headlineContainer}>
              <h3 className={styles.headline}>Issue Status Mappings</h3>
              <p className={styles.description}>Map your own issue statuses to Dev Lake's standard statuses for every issue type</p>
            </div>

            <div className={styles.formContainer}>

              <Tabs id="StatusMappings" onChange={(id) => setStatusTabId(id)} selectedTabId={statusTabId} className={styles.statusTabs}>
                <Tab id={0} title="Bug" panel={
                  <MappingTagStatus
                    reqValue={statusMappingReqBug}
                    resValue={statusMappingResBug}
                    envName="JIRA_ISSUE_BUG_STATUS_MAPPING"
                    clearBtnReq={<ClearButton onClick={() => setStatusMappingReqBug([])} />}
                    clearBtnRes={<ClearButton onClick={() => setStatusMappingResBug([])} />}
                    onChangeReq={(values) => setStatusMappingReqBug(values)}
                    onChangeRes={(values) => setStatusMappingResBug(values)}
                  />
                } />
                <Tab id={1} title="Incident" panel={
                  <MappingTagStatus
                    reqValue={statusMappingReqIncident}
                    resValue={statusMappingResIncident}
                    envName="JIRA_ISSUE_INCIDENT_STATUS_MAPPING"
                    clearBtnReq={<ClearButton onClick={() => setStatusMappingReqIncident([])} />}
                    clearBtnRes={<ClearButton onClick={() => setStatusMappingResIncident([])} />}
                    onChangeReq={(values) => setStatusMappingReqIncident(values)}
                    onChangeRes={(values) => setStatusMappingResIncident(values)}
                  />
                } />
                <Tab id={2} title="Story" panel={
                  <MappingTagStatus
                    reqValue={statusMappingReqStory}
                    resValue={statusMappingResStory}
                    envName="JIRA_ISSUE_STORY_STATUS_MAPPING"
                    clearBtnReq={<ClearButton onClick={() => setStatusMappingReqStory([])} />}
                    clearBtnRes={<ClearButton onClick={() => setStatusMappingResStory([])} />}
                    onChangeReq={(values) => setStatusMappingResStory(values)}
                    onChangeRes={(values) => setStatusMappingResStory(values)}
                  />
                } />

                {customStatus.length > 0 && customStatus.map((status, i) => {
                  const statusAll = customStatus.filter(obj => obj.name != status.name)
                  const mapObj = customStatus.find(obj => obj.name === status.name)

                  return <Tab id={i + 3} key={i} title={status.name} panel={
                    <MappingTagStatus
                      reqValue={mapObj.reqValue.length > 0 ? mapObj.reqValue : []}
                      resValue={mapObj.resValue.length > 0 ? mapObj.resValue : []}
                      envName={`JIRA_ISSUE_${status.name.toUpperCase()}_STATUS_MAPPING`}
                      clearBtnReq={<ClearButton onClick={() => {
                        mapObj.reqValue = []
                        setCustomStatus([...statusAll, mapObj])
                      }} />}
                      clearBtnRes={<ClearButton onClick={() => {
                        mapObj.resValue = []
                        setCustomStatus([...statusAll, mapObj])
                      }} />}
                      onChangeReq={(values) => {
                        mapObj.reqValue = values
                        setCustomStatus([...statusAll, mapObj])
                      }}
                      onChangeRes={(values) => {
                        mapObj.resValue = values
                        setCustomStatus([...statusAll, mapObj])
                      }}
                    />
                  } />
                })}

                <Button icon="add" onClick={() => setCustomStatusOverlay(true)} className={styles.addNewStatusBtn}>Add New</Button>

                <Dialog
                  style={{ width: '100%', maxWidth: "664px", height: "auto" }}
                  icon="diagram-tree"
                  onClose={() => setCustomStatusOverlay(false)}
                  title="Add a New Custom Status"
                  isOpen={customStatusOverlay}
                  onOpened={() => setCustomStatusName('')}
                  autoFocus={false}
                  className={styles.customStatusDialog}
                >
                  <div className={Classes.DIALOG_BODY}>
                  <form onSubmit={(e) => {
                    e.preventDefault()
                    setCustomStatus([...customStatus, {name: customStatusName, reqValue: '', resValue: ''}])
                    setCustomStatusOverlay(false)
                  }}>
                    <FormGroup
                      className={styles.formGroup}
                      className={styles.customStatusFormGroup}
                    >
                      <InputGroup
                        id="custom-status"
                        placeholder="Enter custom status name"
                        onChange={(e) => setCustomStatusName(e.target.value)}
                        className={styles.customStatusInput}
                        autoFocus={true}
                      />
                      <Button icon="add" onClick={() => {
                          setCustomStatus([...customStatus, {name: customStatusName, reqValue: '', resValue: ''}])
                          setCustomStatusOverlay(false)
                        }}
                        className={styles.addNewStatusBtnDialog}
                        onSubmit={(e) => e.preventDefault()}>Add New</Button>
                    </FormGroup>
                  </form>
                  </div>
                </Dialog>

                <Tabs.Expander />
              </Tabs>
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
