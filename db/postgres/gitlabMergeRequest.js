'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class GitlabMergeRequest extends Model {

  }

  GitlabMergeRequest.init({
    id: {
      primaryKey: true,
      type: DataTypes.INTEGER
    },
    title: {
      type: DataTypes.STRING
    },
    projectId: {
      type: DataTypes.INTEGER
    },
    numberOfReviewers: {
      type: DataTypes.INTEGER
    },
    state: {
      type: DataTypes.STRING
    },
    webUrl: {
      type: DataTypes.STRING
    },
    userNotesCount: {
      type: DataTypes.INTEGER
    },
    workInProgress: {
      type: DataTypes.BOOLEAN
    },
    sourceBranch: {
      type: DataTypes.STRING
    },
    mergedAt: {
      type: DataTypes.DATE
    },
    gitlabCreatedAt: {
      type: DataTypes.DATE
    },
    closedAt: {
      type: DataTypes.DATE
    },
    mergedByUsername: {
      type: DataTypes.STRING
    },
    description: {
      type: DataTypes.STRING
    },
    authorUsername: {
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
    modelName: 'GitlabMergeRequest',
    underscored: true
  })

  return GitlabMergeRequest
}
