'use strict'
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('jira_users', {
      self: {
        type: Sequelize.STRING
      },
      accountId: {
        type: Sequelize.STRING,
        field: 'account_id'
      },
      name: {
        type: Sequelize.STRING
      },
      key: {
        type: Sequelize.STRING
      },
      emailAddress: {
        type: Sequelize.STRING,
        primaryKey: true,
        field: 'email_address'
      },
      displayName: {
        type: Sequelize.STRING,
        field: 'display_name'
      },
      active: {
        type: Sequelize.BOOLEAN
      },
      timezone: {
        type: Sequelize.STRING
      },
      createdAt: {
        allowNull: false,
        type: Sequelize.DATE,
        field: 'created_at'
      },
      updatedAt: {
        allowNull: false,
        type: Sequelize.DATE,
        field: 'updated_at'
      }
    })
  },
  down: async (queryInterface, Sequelize) => {
    await queryInterface.dropTable('jira_users')
  }
}
