const env = process.env.NODE_ENV

const loadConfigFile = (filePath) => {
  try {
    return require(filePath)
  } catch (e) {
    console.info(`resolveConfig:loadConfigFile: Could not load ${filePath} error=${e}`)
    console.info(`loading local config by default`)
    return require('./local.js')
  }
}

module.exports = loadConfigFile(`./${env}.js`)
