'use strict'

const {
  Model
} = require('sequelize')

module.exports = (sequelize, DataTypes) => {
  class JiraUser extends Model {
    
  }

  JiraUser.init({
    uuid: {
      primaryKey: true,
      type: DataTypes.UUID
    },
    self: DataTypes.STRING,
    account_id: DataTypes.STRING,
    name: DataTypes.STRING,
    key: DataTypes.STRING,
    email_address: DataTypes.STRING,
    avatar_urls_48x48: DataTypes.STRING,
    avatar_urls_32x32: DataTypes.STRING,
    avatar_urls_24x24: DataTypes.STRING,
    avatar_urls_16x16: DataTypes.STRING,
    display_name: DataTypes.STRING,
    active: DataTypes.BOOLEAN,
    timezone: DataTypes.STRING,
    createdAt: {
      allowNull: false,
      type: DataTypes.DATE,
      defaultValue: DataTypes.NOW
    },
    updatedAt: {
      allowNull: false,
      type: DataTypes.DATE,
      defaultValue: DataTypes.NOW
    },
  }, {
    sequelize,
    modelName: 'JiraJiraUser'
  })

  return JiraUser
}