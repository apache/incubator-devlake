const dbConnector = require('@mongo/connection');
const md5 = require('../../util/md5');

class Schema {
  get valiator () {
    return {
      primary: ['id'],
      properties: {
        id: {
          type: 'int',
          description: 'data unique id',
          required: true
        }
      }
    }
  }

  get mongoValidator () {
    const _validator = this.valiator
    const required = _validator.primary
    const properties = {}
    const convertProperty = (props) => {
      const mongoSchema = {
        bsonType: props.type,
        description: props.description
      }
      if (props.type === 'object') {
        mongoSchema.properties = Object.keys(props.properties).reduce((r, p) => {
          r[p] = convertProperty(props.properties[p])
          return r
        }, {})
      }
      return mongoSchema
    }
    Object.keys(_validator.properties).forEach(property => {
      const props = _validator.properties[property]
      properties[property] = convertProperty(props)
      if (props.required) {
        required.push(property)
      }
    })
    return {
      $jsonSchema: {
        required,
        properties
      }
    }
  }

  constructor (document) {
    this.document = document
  }

  async save () {
    const { db } = await dbConnector.connect()
    const collections = await db.collections()
    const name = this.constructor.name
    let coll = collections.find(c => c.name === name)
    if (!coll) {
      coll = await db.createCollection(name, { validator: this.mongoValidator })
    }
    const { primary } = this.valiator
    const primaryValue = primary.reduce((r, k) => {
      return `${r}${this.document[k]}`
    }, '')
    const documentId = md5(primaryValue)
    this.document._id = documentId
    await coll.findOneAndReplace({ _id: documentId }, this.document, { upsert: true })
  }
}

module.exports = Schema
