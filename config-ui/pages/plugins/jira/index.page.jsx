import Head from 'next/head'
import { useState, useEffect, useRef } from 'react'
import styles from '../../../styles/Home.module.css'
import {
  Tooltip, Position, FormGroup, InputGroup, Button, Label, Icon, Classes, Dialog
} from '@blueprintjs/core'
import dotenv from 'dotenv'
import path from 'path'
import * as fs from 'fs/promises'
import { existsSync } from 'fs'
import Nav from '../../../components/Nav'
import Sidebar from '../../../components/Sidebar'
import Content from '../../../components/Content'
import SaveAlert from '../../../components/SaveAlert'
import MappingTag from './MappingTag'
import MappingTagStatus from './MappingTagStatus'
import ClearButton from './ClearButton'

function parseMapping(mappingString) {
  const mapping = {}
  if (!mappingString.trim()) {
    return mapping
  }
  for (const item of mappingString.split(";")) {
    let [standard, customs] = item.split(":")
    standard = standard.trim()
    mapping[standard] = mapping[standard] || []
    if (!customs) {
      continue
    }
    for (const custom of customs.split(",")) {
      mapping[standard].push(custom.trim())
    }
  }
  return mapping
}

export default function Home(props) {
  const { env } = props

  const [alertOpen, setAlertOpen] = useState(false)
  const [jiraEndpoint, setJiraEndpoint] = useState(env.JIRA_ENDPOINT)
  const [jiraBasicAuthEncoded, setJiraBasicAuthEncoded] = useState(env.JIRA_BASIC_AUTH_ENCODED)
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState(env.JIRA_ISSUE_EPIC_KEY_FIELD)
  const [jiraIssueStoryCoefficient, setJiraIssueStoryCoefficient] = useState(env.JIRA_ISSUE_STORYPOINT_COEFFICIENT)
  const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState(env.JIRA_ISSUE_STORYPOINT_FIELD)

  // Type mappings state
  const defaultTypeMapping = parseMapping(env.JIRA_ISSUE_TYPE_MAPPING)
  const [typeMappingBug, setTypeMappingBug] = useState(defaultTypeMapping.Bug || [])
  const [typeMappingIncident, setTypeMappingIncident] = useState(defaultTypeMapping.Incident || [])
  const [typeMappingRequirement, setTypeMappingRequirement] = useState(defaultTypeMapping.Requirement || [])
  const [typeMappingAll, setTypeMappingAll] = useState()

  // status mapping
  const defaultStatusMappings = []
  for (const [key, value] of Object.entries(env)) {
    const m = /^JIRA_ISSUE_([A-Z]+)_STATUS_MAPPING$/.exec(key)
    if (!m) {
      continue
    }
    const type = m[1]
    defaultStatusMappings.push({
      type,
      key,
      mapping: parseMapping(value)
    })
  }
  const [statusMappings, setStatusMappings] = useState(defaultStatusMappings)
  const filterStatusMappingInput = useRef(null)
  function setStatusMapping(key, values, status) {
    setStatusMappings(statusMappings.map(mapping => {
      if (mapping.key === key) {
        mapping.mapping[status] = values
      }
      return mapping
    }))
  }

  const statusMappingsInput = useRef()
  const [customStatusOverlay, setCustomStatusOverlay] = useState(false)
  const [customStatusName, setCustomStatusName] = useState('')
  const [filteredStatusMappings, setFilteredStatusMappings] = useState(statusMappings)
  function addStatusMapping(e) {
    const type = customStatusName.trim().toUpperCase()
    if (statusMappings.find(e => e.type === type)) {
      return
    }
    const result = [
      ...statusMappings,
      {
        type,
        key: `JIRA_ISSUE_${type}_STATUS_MAPPING`,
        mapping: {
          Resolved: [],
          Rejected: [],
        }
      }
    ]
    statusMappingsInput.current = ''
    setStatusMappings(result)
    setFilteredStatusMappings(result)
    setCustomStatusOverlay(false)
    e.preventDefault()
  }

  function filterStatusMappings(e) {
    const query = e.target.value.toUpperCase()
    if (statusMappings) {
      const filtered = statusMappings.filter(item => {
        if (query === item.key.slice(11, -15)) {
          return true
        } else if (query === '') {
          return true
        }
      })
      setFilteredStatusMappings(filtered)
    }
  }

  function updateEnv(key, value) {
    fetch(`/api/setenv/${key}/${encodeURIComponent(value)}`)
  }

  function saveAll(e) {
    e.preventDefault()
    updateEnv('JIRA_ENDPOINT', jiraEndpoint)
    updateEnv('JIRA_BASIC_AUTH_ENCODED', jiraBasicAuthEncoded)
    updateEnv('JIRA_ISSUE_EPIC_KEY_FIELD', jiraIssueEpicKeyField)
    updateEnv('JIRA_ISSUE_TYPE_MAPPING', typeMappingAll)
    updateEnv('JIRA_ISSUE_STORYPOINT_COEFFICIENT', jiraIssueStoryCoefficient)
    updateEnv('JIRA_ISSUE_STORYPOINT_FIELD', jiraIssueStoryPointField)

    // Save all custom status data
    statusMappings.map(mapping => {
      const { Resolved, Rejected } = mapping.mapping
      updateEnv(mapping.key, `Rejected:${Rejected ? Rejected.join(',') : ''};Resolved:${Resolved ? Resolved.join(',') : ''};`)
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
              placeholderText="Add Issue Types..."
              values={typeMappingBug}
              helperText="JIRA_ISSUE_TYPE_MAPPING"
              rightElement={<ClearButton onClick={() => setTypeMappingBug([])} />}
              onChange={(values) => setTypeMappingBug(values)}
            />

            <MappingTag
              labelName="Incident"
              labelIntent="warning"
              typeOrStatus="type"
              placeholderText="Add Issue Types..."
              values={typeMappingIncident}
              helperText="JIRA_ISSUE_TYPE_MAPPING"
              rightElement={<ClearButton onClick={() => setTypeMappingIncident([])} />}
              onChange={(values) => setTypeMappingIncident(values)}
            />

            <MappingTag
              labelName="Requirement"
              labelIntent="primary"
              typeOrStatus="type"
              placeholderText="Add Issue Types..."
              values={typeMappingRequirement}
              helperText="JIRA_ISSUE_TYPE_MAPPING"
              rightElement={<ClearButton onClick={() => setTypeMappingRequirement([])} />}
              onChange={(values) => setTypeMappingRequirement(values)}
            />

            <div className={styles.headlineContainer}>
              <h3 className={styles.headline}>Issue Status Mappings</h3>
              <p className={styles.description}>Map your own issue statuses to Dev Lake's standard statuses for every issue type</p>
            </div>

            <div className={styles.jiraFormContainer}>

              <FormGroup
                  label=""
                  inline={true}
                  labelFor="filter-status-mappings"
                  className={styles.statusMappingsFilterFormGroup}
                  contentClassName={styles.formGroup}
                >
                <Label>
                  <b>Filter&nbsp;status&nbsp;mappings</b><br/><br/>
                  <InputGroup
                    id="filter-status-mappings"
                    leftIcon="search"
                    placeholder="Search by name"
                    ref={statusMappingsInput}
                    onChange={(e) => filterStatusMappings(e)}
                    className={styles.input}
                  />
                </Label>
              </FormGroup>

              {(statusMappings && statusMappings.length > 0) && filteredStatusMappings.map((statusMapping, i) =>
                <div
                key={statusMapping.key} className={styles.jiraFormContainerItem}>
                  <p>Mapping {statusMapping.type} </p>
                  <div>
                    <MappingTagStatus
                      reqValue={statusMapping.mapping.Rejected || []}
                      resValue={statusMapping.mapping.Resolved || []}
                      envName={statusMapping.key}
                      clearBtnReq={<ClearButton onClick={() => setStatusMapping(statusMapping.key, [], 'Rejected')} />}
                      clearBtnRes={<ClearButton onClick={() => setStatusMapping(statusMapping.key, [], 'Resolved')} />}
                      onChangeReq={values => setStatusMapping(statusMapping.key, values, 'Rejected')}
                      onChangeRes={values => setStatusMapping(statusMapping.key, values, 'Resolved')}
                      className={styles.mappingTagStatus}
                    />
                  </div>
                </div>
              )}
              <Button icon="add" onClick={() => setCustomStatusOverlay(true)} className={styles.addNewStatusBtn}>Add New</Button>

              <Dialog
                icon="diagram-tree"
                onClose={() => setCustomStatusOverlay(false)}
                title="Add a New Status Mapping"
                isOpen={customStatusOverlay}
                onOpened={() => setCustomStatusName('')}
                autoFocus={false}
                className={styles.customStatusDialog}
              >
                <div className={Classes.DIALOG_BODY}>
                <form onSubmit={addStatusMapping}>
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
                    <Button icon="add" onClick={addStatusMapping}
                      className={styles.addNewStatusBtnDialog}
                      onSubmit={(e) => e.preventDefault()}>Add New</Button>
                  </FormGroup>
                </form>
                </div>
              </Dialog>

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
