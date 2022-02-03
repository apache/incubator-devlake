const HttpHelper = require('./util/http')
const apiUrls = require('./config/apiUrls');
const fs = require('fs')
/* Main should:
  1. Read plugins from config.
  2. Get paths to request from plugins.
  3. Request endpoints to plugin API.
  4. Handle errors.
  5. Write response JSON to file system.
*/
async function main(){
  await fetchDataFromApis()
}

async function fetchDataFromApis(){
  // Read plugins from config.

  dbJson = {}
  for(let apiUrl of apiUrls){
    await HttpHelper.get(apiUrl.url, (res, err) => {
      if(err){
        // Handle errors
        console.log('err', err);
      } else {
        // Write response data to file
        console.log('res.data[0]', res.data);
        dbJson[apiUrl.name] = res.data
      }
    })
  }

  fs.writeFileSync(`./db.json`, JSON.stringify(dbJson))
}

main()
