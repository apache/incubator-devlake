import Head from 'next/head'
import { useState } from 'react'
import dotenv from 'dotenv'
import path from 'path'
import * as fs from 'fs/promises'
import { existsSync } from 'fs';
import styles from '../styles/Home.module.css'
import { FormGroup, InputGroup, Button } from "@blueprintjs/core"
import Nav from '../components/Nav'
import Sidebar from '../components/Sidebar'
import Content from '../components/Content'

export default function Home(props) {
  const { env } = props

  const [dbUrl, setDbUrl] = useState(env.DB_URL)
  const [port, setPort] = useState(env.PORT)
  const [mode, setMode] = useState(env.MODE)
  const [jiraEndpoint, setJiraEndpoint] = useState(env.JIRA_ENDPOINT)
  const [jiraBasicAuthEncoded, setJiraBasicAuthEncoded] = useState(env.JIRA_BASIC_AUTH_ENCODED)
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState(env.JIRA_ISSUE_EPIC_KEY_FIELD)
  const [gitlabEndpoint, setGitlabEndpoint] = useState(env.GITLAB_ENDPOINT)
  const [gitlabAuth, setGitlabAuth] = useState(env.GITLAB_AUTH)

  function updateEnv(key, value) {
    fetch(`/api/setenv/${key}/${encodeURIComponent(value)}`)
  }

  function saveAll(e) {
    e.preventDefault()
    updateEnv('DB_URL', dbUrl)
    updateEnv('PORT', port)
    updateEnv('MODE', mode)
    updateEnv('JIRA_ENDPOINT', jiraEndpoint)
    updateEnv('JIRA_BASIC_AUTH_ENCODED', jiraBasicAuthEncoded)
    updateEnv('JIRA_ISSUE_EPIC_KEY_FIELD', jiraIssueEpicKeyField)
    updateEnv('GITLAB_ENDPOINT', gitlabEndpoint)
    updateEnv('GITLAB_AUTH', gitlabAuth)
    alert('Config file updated, please restart devlake to apply new configuration')
  }

  return (
    <div className={styles.container}>

      <Head>
        <title>Devlake Config-UI</title>
        <meta name="description" content="Lake: Config" />
        <link rel="icon" href="/favicon.ico" />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin />
        <link href="https://fonts.googleapis.com/css2?family=Source+Sans+Pro:wght@400;600&display=swap" rel="stylesheet" />
        <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500;600&display=swap" rel="stylesheet" />
      </Head>

      <Nav />
      <Sidebar />
      <Content>
        <main className={styles.main}>

          <div className={styles.headlineContainer}>
            <h1>Configuration</h1>
            <p className={styles.description}>Configure your <code className={styles.code}>.env</code> file values</p>
          </div>

          <form className={styles.form}>
            <div className={styles.headlineContainer}>
              <h2 className={styles.headline}>Main Database Connection</h2>
              <p className={styles.description}>Settings for the mySQL database</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="DB_URL"
                inline={true}
                labelFor="db-url"
                helperText="The URL Connection string to the database"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="db-url"
                  placeholder="Enter DB Connection String"
                  defaultValue={dbUrl}
                  onChange={(e) => setDbUrl(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <div className={styles.headlineContainer}>
              <h2 className={styles.headline}>REST Configuration</h2>
              <p className={styles.description}>Configure main REST Settings</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="PORT"
                inline={true}
                labelFor="port"
                helperText="The main port for the REST server"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="port"
                  placeholder="Enter Port eg. :8080"
                  defaultValue={port}
                  onChange={(e) => setPort(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="MODE"
                inline={true}
                labelFor="mode"
                helperText="The development mode for the server"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="mode"
                  placeholder="Enter Mode eg. debug"
                  defaultValue={mode}
                  onChange={(e) => setMode(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <div className={styles.headlineContainer}>
              <h2 className={styles.headline}>Jira Configuration</h2>
              <p className={styles.description}>Jira Account and config settings</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="JIRA_ENDPOINT"
                inline={true}
                labelFor="jira-endpoint"
                helperText="Your custom url endpoint for Jira"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="jira-endpoint"
                  placeholder="Enter Jira endpoint eg. https://merico.atlassian.net"
                  defaultValue={jiraEndpoint}
                  onChange={(e) => setJiraEndpoint(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="JIRA_BASIC_AUTH_ENCODED"
                inline={true}
                labelFor="jira-basic-auth"
                helperText="Your encoded Jira auth token"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="jira-basic-auth"
                  placeholder="Enter Jira Auth eg. EJrLG8DNeXADQcGOaaaX4B47"
                  defaultValue={jiraBasicAuthEncoded}
                  onChange={(e) => setJiraBasicAuthEncoded(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="JIRA_ISSUE_EPIC_KEY_FIELD"
                inline={true}
                labelFor="jira-epic-key"
                helperText="Your custom epic key field (optional)"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="jira-epic-key"
                  placeholder="Enter Jira epic key field"
                  defaultValue={jiraIssueEpicKeyField}
                  onChange={(e) => setJiraIssueEpicKeyField(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <div className={styles.headlineContainer}>
              <h2 className={styles.headline}>Gitlab Configuration</h2>
              <p className={styles.description}>Gitlab account and config settings</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="GITLAB_ENDPOINT"
                inline={true}
                labelFor="gitlab-endpoint"
                helperText="Gitlab API Endpoint"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="gitlab-endpoint"
                  placeholder="Enter Gitlab API endpoint"
                  defaultValue={gitlabEndpoint}
                  onChange={(e) => setGitlabEndpoint(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                label="GITLAB_AUTH"
                inline={true}
                labelFor="gitlab-auth"
                helperText="Gitlab Auth Token"
                className={styles.formGroup}
                contentClassName={styles.formGroup}
              >
                <InputGroup
                  id="gitlab-auth"
                  placeholder="Enter Gitlab Auth Token eg. uJVEDxabogHbfFyu2riz"
                  defaultValue={gitlabAuth}
                  onChange={(e) => setGitlabAuth(e.target.value)}
                  className={styles.input}
                />
              </FormGroup>
            </div>

            <Button type="submit" outlined={true} large={true} className={styles.saveBtn} onClick={saveAll}>Save Config</Button>
          </form>
        </main>
      </Content>
    </div>
  )
}

export async function getStaticProps() {
  // const fs = require('fs').promises

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
