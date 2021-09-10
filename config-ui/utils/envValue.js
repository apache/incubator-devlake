import path from 'path'
import os from 'os'
import fs from 'fs'

const envFilePath = path.join(process.cwd(), 'data', '../../.env')

// read .env file
const readEnvVars = () => fs.readFileSync(envFilePath, "utf-8").split(os.EOL)

// get .env values
export const getEnvValue = (key) => {
  const matchedLine = readEnvVars().find((line) => line.split("=")[0] === key)

  return matchedLine !== undefined ? matchedLine.split("=")[1] : null
}

// create/override .env values
export const setEnvValue = (key, value) => {
  const envVars = readEnvVars()
  const targetLine = envVars.find((line) => line.split("=")[0] === key)

  if (targetLine !== undefined) {
    // update existing line
    const targetLineIndex = envVars.indexOf(targetLine)
    envVars.splice(targetLineIndex, 1, `${key}=${value}`)
  } else {
    // create new key value
    envVars.push(`${key}=${value}`)
  }

  fs.writeFileSync(envFilePath, envVars.join(os.EOL))

  // required for next
  return null
}

//* Usage
// console.log(getEnvValue('value'))
// setEnvValue('ENV_KEY', 'value')
