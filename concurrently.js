const concurrently = require('concurrently');

concurrently([
  "node src/collection/main.js",
  "node src/collection/worker.js",
  "node src/enrichment/main.js",
  "node src/enrichment/worker.js"
], {
  killOthers: ['failure', 'success']
}).then(
    function onSuccess(exitInfo) {
      // This code is necessary to make sure the parent terminates 
      // when the application is closed successfully.
      process.exit();
    },
    function onFailure(exitInfo) {
      // This code is necessary to make sure the parent terminates 
      // when the application is closed because of a failure.
      process.exit();
    }
  );