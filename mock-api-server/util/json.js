const JsonHelper = {
  getPathsFromPlugin(plugin){
    // This is just sample code.
    // TODO: We need to actually handle params in relative paths.
    let paths = []
    let baseUrl = plugin.baseUrl
    plugin.collections.forEach(collection => {
      paths.push(`${baseUrl}${collection.relativeUrl}`)
    })
    return paths
  }
}
module.exports = JsonHelper