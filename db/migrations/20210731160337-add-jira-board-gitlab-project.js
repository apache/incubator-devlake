'use strict'
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('jira_board_gitlab_projects', {
      boardId: {
        allowNull: false,
        primaryKey: true,
        type: Sequelize.INTEGER,
        field: 'board_id'
      },
      projectId: {
        allowNull: false,
        primaryKey: true,
        type: Sequelize.INTEGER,
        field: 'project_id'
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
    await queryInterface.dropTable('jira_board_gitlab_projects')
  }
}
