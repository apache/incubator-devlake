import Head from 'next/head'
import { useState } from 'react'
import dotenv from 'dotenv'
import path from 'path'
import * as fs from 'fs/promises'
import { existsSync } from 'fs';
import styles from '../styles/Home.module.css'
import { FormGroup, InputGroup, Button, Alert } from '@blueprintjs/core'
import Nav from '../components/Nav'
import Sidebar from '../components/Sidebar'
import Content from '../components/Content'

export default function Home(props) {
  const { env } = props

  const [dbUrl, setDbUrl] = useState(env.DB_URL)
  const [port, setPort] = useState(env.PORT)
  const [mode, setMode] = useState(env.MODE)
  const [alertOpen, setAlertOpen] = useState(false)

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

          <div className={styles.headlineContainer}>
            <h1>Configuration</h1>
            <p className={styles.description}>Configure your <code className={styles.code}>.env</code> file values</p>
          </div>

          <form className={styles.form}>
            <div className={styles.headlineContainer}>
              <h2 className={styles.headline}>Devlake</h2>
              <p className={styles.description}>Settings for the Devlake framework</p>
            </div>

            <Card className={styles.formSection}>
              <h3>Basic (DO NOT CHANGE THIS SECTION UNLESS YOUR ARE DEVELOPER)</h3>
              <div className={styles.formContainer}>
                <FormGroup
                  label="*DB_URL"
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
              <div className={styles.formContainer}>
                <FormGroup
                  label="*PORT"
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
            </Card>


            <Button type="submit" outlined={true} large={true} className={styles.saveBtn} onClick={saveAll}>Save Config</Button>
          </form>
        </main>
      </Content>
    </div>

    <Alert
      canEscapeKeyCancel={true}
      canOutsideClickCancel={true}
      confirmButtonText="Ok"
      isOpen={alertOpen}
      onClose={() => setAlertOpen(false)}>
      <h4>Config File Updated</h4>
      <p>To apply new configuration, restart devlake by running: <br/><br/><code>docker-compose up -d</code></p>
    </Alert>
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
