import Head from 'next/head'
import { useState } from 'react'
import dotenv from 'dotenv'
import path from 'path'
import * as fs from 'fs/promises'
import styles from '../styles/Home.module.css'
import Nav from '../components/Nav'
import Sidebar from '../components/Sidebar'
import Content from '../components/Content'

export default function Home(props) {
  const { env } = props

  const [dbUrl, setDbUrl] = useState('')
  const [port, setPort] = useState('')
  const [mode, setMode] = useState('')

  function updateEnv(key, value) {
    fetch(`http://localhost:4000/api/setenv/${key}/${value}`)
    alert('updated')
  }

  return (
    <div className={styles.container}>

      <Head>
        <title>Create Next App</title>
        <meta name="description" content="Lake: Config" />
        <link rel="icon" href="/favicon.ico" />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
        <link href="https://fonts.googleapis.com/css2?family=Source+Sans+Pro:wght@400;600&display=swap" rel="stylesheet" />
      </Head>

      <Nav />
      <Sidebar />
      <Content />

      {/* <main className={styles.main}>

        <img src="/logo.svg" className={styles.logo} />

        <p className={styles.description}>Configure your <code className={styles.code}>.env</code> file values</p>

        <div className={styles.formContainer}>
          <h3 className={styles.headline}>Main Database Connection</h3>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>DB_URL</label>
          <input className={styles.input} type="text" onChange={(e) => setDbUrl(e.target.value)} defaultValue={env.DB_URL} />
          <button className={styles.button} onClick={() => updateEnv('DB_URL', dbUrl)}>save</button>
        </div>

        <div className={styles.formContainer}>
          <h3 className={styles.headline}>REST Configuration</h3>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>PORT</label>
          <input className={styles.input} type="text" onChange={(e) => setPort(e.target.value)} defaultValue={env.PORT} />
          <button className={styles.button} onClick={() => updateEnv('PORT', port)}>save</button>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>MODE</label>
          <input className={styles.input} type="text" onChange={(e) => setMode(e.target.value)} defaultValue={env.MODE} />
          <button className={styles.button} onClick={() => updateEnv('MODE', mode)}>save</button>
        </div>

        <div className={styles.formContainer}>
          <h3 className={styles.headline}>Jira Configuration</h3>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>JIRA_ENDPOINT</label>
          <input className={styles.input} type="text" onChange={(e) => setJiraEndpoint(e.target.value)} defaultValue={env.JIRA_ENDPOINT} />
          <button className={styles.button} onClick={() => updateEnv('JIRA_ENDPOINT', jiraEndpoint)}>save</button>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>JIRA_BASIC_AUTH_ENCODED</label>
          <input className={styles.input} type="text" onChange={(e) => setJiraBasicAuthEncoded(e.target.value)} defaultValue={env.JIRA_BASIC_AUTH_ENCODED} />
          <button className={styles.button} onClick={() => updateEnv('JIRA_BASIC_AUTH_ENCODED', jiraBasicAuthEncoded)}>save</button>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>JIRA_ISSUE_EPIC_KEY_FIELD</label>
          <input className={styles.input} type="text" onChange={(e) => setJiraIssueEpicKeyField(e.target.value)} defaultValue={env.JIRA_ISSUE_EPIC_KEY_FIELD} />
          <button className={styles.button} onClick={() => updateEnv('JIRA_ISSUE_EPIC_KEY_FIELD', jiraIssueEpicKeyField)}>save</button>
        </div>

        <div className={styles.formContainer}>
          <h3 className={styles.headline}>Gitlab Configuration</h3>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>GITLAB_ENDPOINT</label>
          <input className={styles.input} type="text" onChange={(e) => setGitlabEndpoint(e.target.value)} defaultValue={env.GITLAB_ENDPOINT} />
          <button className={styles.button} onClick={() => updateEnv('GITLAB_ENDPOINT', gitlabEndpoint)}>save</button>
        </div>

        <div className={styles.formContainer}>
          <label className={styles.label}>GITLAB_AUTH</label>
          <input className={styles.input} type="text" onChange={(e) => setGitlabAuth(e.target.value)} defaultValue={env.GITLAB_AUTH} />
          <button className={styles.button} onClick={() => updateEnv('GITLAB_AUTH', gitlabAuth)}>save</button>
        </div>

      </main> */}
    </div>
  )
}

export async function getStaticProps() {
  // const fs = require('fs').promises

  const filePath = path.join(process.cwd(), 'data', '../../config-ui/.env')
  const fileData = await fs.readFile(filePath)
  const env = dotenv.parse(fileData)

  return {
    props: {
      env
    },
  }
}
