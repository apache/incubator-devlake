'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class JiraBoardGitlabProject extends Model {
    /**
     * Helper method for defining associations.
     * This method is not a part of Sequelize lifecycle.
     * The `models/index` file will call this method automatically.
     */
    static associate (models) {
      // define association here
    }
  };

  JiraBoardGitlabProject.init({
    boardId: {
      type: DataTypes.INTEGER,
      primaryKey: true
    },
    projectId: {
      type: DataTypes.INTEGER,
      primaryKey: true
    }
  }, {
    sequelize,
    modelName: 'JiraBoardGitlabProject',
    underscored: true
  })
  return JiraBoardGitlabProject
}
