const { gitlab } = require('@config/resolveConfig')

module.exports = {
  async maybeSkip (promise, key) {
    if (gitlab.skip) {
      !gitlab.skip[key] && await promise
    } else {
      await promise
    }
  }
}
