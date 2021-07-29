'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class GitlabProject extends Model {

  }

  GitlabProject.init({
    uuid: {
      primaryKey: true,
      type: DataTypes.UUID,
      defaultValue: DataTypes.UUIDV4
    },
    name: {
      type: DataTypes.STRING
    },
    id: {
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
