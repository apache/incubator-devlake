module.exports = {
  /**
   *
   * @param {String} mapValues
   * @returns key
   *
   * Sometimes, users in Jira can set their own labels for issue statuses like 'Done'
   * This method allows users to set their own labels from their Jira system in the /config/constants.json file.
   */
  mapValue (input, config) {
    if (!input || input === '') {
      return ''
    }

    for (const key in config) {
      const value = config[key]
      if (Array.isArray(value)) {
        const matchFound = value.map(x => x.toLowerCase()).includes(input.toLowerCase())
        if (matchFound) {
          return key
        }
      } else {
        if (value.toLowerCase() === input.toLowerCase()) {
          return key
        }
      }
    }
    // If no mapping is found, return the original value from the Jira API
    return input
  }
}
