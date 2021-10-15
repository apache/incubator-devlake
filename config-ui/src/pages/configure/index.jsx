import React, { useState } from 'react'
import dotenv from 'dotenv'
import path from 'path'
import * as fs from 'fs/promises'
import { existsSync } from 'fs';
import styles from '../styles/Home.module.css'
import { FormGroup, InputGroup, Button, Alert, Tooltip, Position, Label } from '@blueprintjs/core'
import Nav from '../components/Nav'
import Sidebar from '../components/Sidebar'
import Content from '../components/Content'
import SaveAlert from '../components/SaveAlert'

export default function Home(props) {
  const { env } = props

  const [alertOpen, setAlertOpen] = useState(false)
  const [dbUrl, setDbUrl] = useState(env.DB_URL)
  const [port, setPort] = useState(env.PORT)
  const [mode, setMode] = useState(env.MODE)

  function updateEnv(key, value) {
    fetch(`/api/setenv/${key}/${encodeURIComponent(value)}`)
  }

  function saveAll(e) {
    e.preventDefault()
    updateEnv('DB_URL', dbUrl)
    updateEnv('PORT', port)
    updateEnv('MODE', mode)
    setAlertOpen(true)
  }

  return (
    <>
    <div className={styles.container}>

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
              <h2 className={styles.headline}>Dev Lake Database Connection</h2>
              <p className={styles.description}>Settings for the MySQL database</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="db-url"
                className={styles.formGroup}
                helperText="DB_URL"
                contentClassName={styles.formGroup}
              >
                <Tooltip content="The URL Connection string to the database" position={Position.TOP}>
                  <Label>
                    Database&nbsp;URL <span className={styles.requiredStar}>*</span>
                    <InputGroup
                      id="db-url"
                      placeholder="Enter DB Connection String"
                      defaultValue={dbUrl}
                      onChange={(e) => setDbUrl(e.target.value)}
                      className={styles.input}
                    />
                  </Label>
                </Tooltip>
              </FormGroup>
            </div>

            <div className={styles.headlineContainer}>
              <h2 className={styles.headline}>Dev Lake API Server</h2>
              <p className={styles.description}>Configure main REST Settings</p>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="port"
                className={styles.formGroup}
                helperText="PORT"
                contentClassName={styles.formGroup}
              >
                <Tooltip content="The main port for the REST server" position={Position.TOP}>
                  <Label>
                    Port <span className={styles.requiredStar}>*</span>
                    <InputGroup
                      id="port"
                      placeholder="Enter Port eg. :8080"
                      defaultValue={port}
                      onChange={(e) => setPort(e.target.value)}
                      className={styles.input}
                    />
                  </Label>
                </Tooltip>
              </FormGroup>
            </div>

            <div className={styles.formContainer}>
              <FormGroup
                inline={true}
                labelFor="mode"
                className={styles.formGroup}
                helperText="MODE"
                contentClassName={styles.formGroup}
              >
                <Tooltip content="The development mode for the server" position={Position.TOP}>
                  <Label>
                    Mode
                    <InputGroup
                      id="mode"
                      placeholder="Enter Mode eg. debug"
                      defaultValue={mode}
                      onChange={(e) => setMode(e.target.value)}
                      className={styles.input}
                    />
                  </Label>
                </Tooltip>
              </FormGroup>
            </div>

            <Button type="submit" outlined={true} large={true} className={styles.saveBtn} onClick={saveAll}>Save Config</Button>
          </form>
        </main>
      </Content>
    </div>

    <SaveAlert alertOpen={alertOpen} onClose={() => setAlertOpen(false)} />
    </>
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
