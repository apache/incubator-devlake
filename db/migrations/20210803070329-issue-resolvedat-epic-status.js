'use strict';

module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.addColumn('jira_issues', 'issue_resolved_at', {
      type: Sequelize.DataTypes.DATE,
    })
    await queryInterface.addColumn('jira_issues', 'epic_key', {
      type: Sequelize.DataTypes.TEXT,
    })
    await queryInterface.addColumn('jira_issues', 'status', {
      type: Sequelize.DataTypes.TEXT,
    })
  },

  down: async (queryInterface, Sequelize) => {
    await queryInterface.removeColumn('jira_issues', 'issue_resolved_at')
    await queryInterface.removeColumn('jira_issues', 'epic_key')
    await queryInterface.removeColumn('jira_issues', 'status')
  }
};
