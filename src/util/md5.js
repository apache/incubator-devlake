const crypto = require('crypto')

module.exports = function md5 (value) {
  return crypto.createHash('md5').update(value).digest('hex')
}
