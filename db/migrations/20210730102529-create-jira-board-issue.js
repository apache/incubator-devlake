'use strict'

module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('jira_board_issues', {
      boardId: {
        type: Sequelize.INTEGER,
        allowNull: false,
        primaryKey: true,
        field: 'board_id'
      },
      issueId: {
        type: Sequelize.INTEGER,
        allowNull: false,
        primaryKey: true,
        field: 'issue_id'
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
    await queryInterface.dropTable('jira_board_issues')
  }
}
