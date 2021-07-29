'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class GitlabCommit extends Model {

  }

  GitlabCommit.init({
    uuid: {
      primaryKey: true,
      type: DataTypes.UUID,
      defaultValue: DataTypes.UUIDV4
    },
    id: {
      type: DataTypes.STRING
    },
    shortId: {
      type: DataTypes.STRING,
    },
    title: {
      type: DataTypes.STRING
    },
    message: {
      type: DataTypes.STRING
    },
    authorName: {
      type: DataTypes.STRING,
    },
    authorEmail: {
      type: DataTypes.STRING,
    },
    authoredDate: {
      type: DataTypes.STRING,
    },
    committerName: {
      type: DataTypes.STRING,
    },
    committerEmail: {
      type: DataTypes.STRING,
    },
    committedDate: {
      type: DataTypes.STRING,
    },
    webUrl: {
      type: DataTypes.STRING,
    },
    additions: {
      type: DataTypes.INTEGER,
    },
    deletions: {
      type: DataTypes.INTEGER,
    },
    total: {
      type: DataTypes.INTEGER,
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
