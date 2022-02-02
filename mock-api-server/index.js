const HttpHelper = require('./util/http')
const plugins = require('./config/plugins.json');
const { getPathsFromPlugin } = require('./util/json');

/* Main should:
  1. Read plugins from config.
  2. Get paths to request from plugins.
  3. Request endpoints to plugin API.
  4. Handle errors.
  5. Write response JSON to file system.
*/
async function main(){
  // Read plugins from config.
  for(let plugin of plugins){
    // Get paths to request from plugins.
    paths = getPathsFromPlugin(plugin)
    for(let path of paths){
      // Request endpoints to plugin API.
      await HttpHelper.get(path, (res, err) => {
        if(err){
          // Handle errors
          console.log('err', err);
        } else {
          // Write response data to file
          console.log('res.data[0]', res.data[0]);
        }
      })
    }
  }
}

main()