const { getEndpointsFromPluginJson } = require('../util/json');
const plugins = require('../config/plugins.json');

describe('getEndpointsFromPluginJson', () => {
  it('gets full url with params', () => {
    let expected = 'https://api.github.com/repos/'
    endpoints = getEndpointsFromPluginJson(plugins[0])
    console.log('endpoints', endpoints);
  })
})