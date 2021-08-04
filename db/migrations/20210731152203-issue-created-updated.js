'use strict'

module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.addColumn('jira_issues', 'issue_created_at', {
      type: Sequelize.DataTypes.DATE
    })
    await queryInterface.addColumn('jira_issues', 'issue_updated_at', {
      type: Sequelize.DataTypes.DATE
    })
  },

  down: async (queryInterface, Sequelize) => {
    await queryInterface.removeColumn('jira_issues', 'issue_created_at')
    await queryInterface.removeColumn('jira_issues', 'issue_updated_at')
  }
}
