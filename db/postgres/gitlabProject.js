'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class GitlabProject extends Model {

  }

  GitlabProject.init({
    name: {
      type: DataTypes.STRING
    },
    id: {
      primaryKey: true,
      type: DataTypes.INTEGER
    },
    pathWithNamespace: {
      type: DataTypes.STRING
    },
    webUrl: {
      type: DataTypes.STRING
    },
    visibility: {
      type: DataTypes.STRING
    },
    openIssuesCount: {
      type: DataTypes.INTEGER
    },
    starCount: {
      type: DataTypes.INTEGER
    },
    createdAt: {
      allowNull: false,
      type: DataTypes.DATE,
      defaultValue: DataTypes.NOW
    },
    updatedAt: {
      allowNull: false,
      type: DataTypes.DATE,
      defaultValue: DataTypes.NOW
    }
  }, {
    sequelize,
    modelName: 'GitlabProject',
    underscored: true
  })

  return GitlabProject
}
