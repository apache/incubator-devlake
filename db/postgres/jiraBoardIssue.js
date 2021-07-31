'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class JiraBoardIssue extends Model {
    /**
     * Helper method for defining associations.
     * This method is not a part of Sequelize lifecycle.
     * The `models/index` file will call this method automatically.
     */
    static associate(models) {
      // define association here
    }
  };

  JiraBoardIssue.init({
    boardId: {
      type: DataTypes.INTEGER,
      primaryKey: true
    },
    issueId: {
      type: DataTypes.INTEGER,
      primaryKey: true
    }
  }, {
    sequelize,
    modelName: 'JiraBoardIssue',
    underscored: true
  })

  return JiraBoardIssue
}