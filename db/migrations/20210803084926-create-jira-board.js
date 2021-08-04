'use strict'
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('jira_boards', {
      id: {
        allowNull: false,
        primaryKey: true,
        type: Sequelize.INTEGER,
        field: 'id'
      },
      projectId: {
        type: Sequelize.INTEGER,
        field: 'project_id'
      },
      name: {
        type: Sequelize.TEXT,
        field: 'name'
      },
      webUrl: {
        type: Sequelize.TEXT,
        field: 'web_url'
      },
      type: {
        type: Sequelize.TEXT,
        field: 'type'
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
    await queryInterface.dropTable('jira_boards')
  }
}
