'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class GitlabCommit extends Model {

  }

  GitlabCommit.init({
    projectId: {
      type: DataTypes.INTEGER
    },
    id: {
      primaryKey: true,
      type: DataTypes.STRING
    },
    shortId: {
      type: DataTypes.STRING
    },
    title: {
      type: DataTypes.TEXT
    },
    message: {
      type: DataTypes.TEXT
    },
    authorName: {
      type: DataTypes.STRING
    },
    authorEmail: {
      type: DataTypes.STRING
    },
    authoredDate: {
      type: DataTypes.STRING
    },
    committerName: {
      type: DataTypes.STRING
    },
    committerEmail: {
      type: DataTypes.STRING
    },
    committedDate: {
      type: DataTypes.STRING
    },
    webUrl: {
      type: DataTypes.STRING
    },
    additions: {
      type: DataTypes.INTEGER
    },
    deletions: {
      type: DataTypes.INTEGER
    },
    total: {
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
    modelName: 'GitlabCommit',
    underscored: true
  })

  return GitlabCommit
}
