// class ORM {
//   constructor (modelName) {
//     this.model = require('@db/models')[modelName]
//   }

//   async findOne (where) {
//     try {
//       return await this.model.findOne({
//         where
//       })
//     } catch (error) {
//       throw new Error(error)
//     }
//   }

//   async update (values, where) {
//     try {
//       return await this.model.findOne({
//         values,
//         where
//       })
//     } catch (error) {
//       throw new Error(error)
//     }
//   }

//   async findOrCreate (newRecord, whereClause) {
//     try {
//       const result = await this.model.findOrCreate({
//         where: whereClause,
//         defaults: newRecord
//       })
//       return result[0]
//     } catch (error) {
//       console.error(error)
//       throw error
//     }
//   }

//   async findAll (where) {
//     try {
//       return await this.model.findAll({
//         where
//       })
//     } catch (error) {
//       console.error(error)
//       throw error
//     }
//   }

//   async delete (where) {
//     try {
//       await this.model.destroy({
//         where
//       })
//     } catch (error) {
//       console.error('ERROR: asdfh8ehe: ', error)
//       throw error
//     }
//   }
// }
