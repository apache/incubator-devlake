const JsonHelper = {
  getEndpointsFromPluginJson(plugin){
    // This is just sample code.
    // TODO: We need to actually handle params in relative urls.
    let endpoints = []
    let baseUrl = plugin.baseUrl
    plugin.collections.forEach(collection => {
      endpoints.push(`${baseUrl}${collection.relativeUrl}`)
    })
    return endpoints
  },
  getEndpointWithParams(params)
}
module.exports = JsonHelper