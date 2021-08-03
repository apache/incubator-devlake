'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class GitlabMergeRequestNote extends Model {

  }

  GitlabMergeRequestNote.init({
    id: {
      primaryKey: true,
      type: DataTypes.STRING
    },
    noteableId: {
      type: DataTypes.STRING
    },
    noteableIid: {
      type: DataTypes.STRING
    },
    authorUsername: {
      type: DataTypes.STRING
    },
    body: {
      type: DataTypes.TEXT
    },
    gitlabCreatedAt: {
      type: DataTypes.STRING
    },
    noteableType: {
      type: DataTypes.STRING
    },
    confidential: {
      type: DataTypes.STRING
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
    modelName: 'GitlabMergeRequestNote',
    underscored: true
  })

  return GitlabMergeRequestNote
}
