const env = process.env.NODE_ENV

const loadConfigFile = (filePath) => {
  try {
    let file = require(filePath)
    console.info(`INFO: loading config`, filePath)
    return file
  } catch (e) {
    console.info(`INFO: loading local config by default`)
    return require('./local.js')
  }
}

module.exports = loadConfigFile(`./${env}.js`)
