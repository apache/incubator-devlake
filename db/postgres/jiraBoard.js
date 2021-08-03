'use strict';

const {
  Model
} = require('sequelize');

module.exports = (sequelize, DataTypes) => {
  class JiraBoard extends Model {
    /**
     * Helper method for defining associations.
     * This method is not a part of Sequelize lifecycle.
     * The `models/index` file will call this method automatically.
     */
    static associate(models) {
      // define association here
    }
  };

  JiraBoard.init({
    id: {
      type: DataTypes.INTEGER,
      primaryKey: true
    },
    projectId: {
      type: DataTypes.INTEGER,
      allowNull: false
    },
    name: {
      type: DataTypes.TEXT,
      allowNull: false
    },
    webUrl: {
      type: DataTypes.TEXT
    },
    type: {
      type: DataTypes.TEXT
    }
  }, {
    sequelize,
    modelName: 'JiraBoard',
    underscored: true
  })

  return JiraBoard
}