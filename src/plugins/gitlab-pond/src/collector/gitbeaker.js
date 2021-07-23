const { Gitlab } = require('@gitbeaker/node')
const { gitlab: { token } } = require('../../../../../../config/resolveConfig')

const gitbeaker = async () => {
  return new Gitlab({
    token
  })
}
module.exports = gitbeaker
