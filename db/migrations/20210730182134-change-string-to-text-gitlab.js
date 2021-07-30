'use strict'

module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.changeColumn('gitlab_commits', 'message', {
      type: Sequelize.TEXT,
      allowNull: true,
    })
    await queryInterface.changeColumn('gitlab_commits', 'title', {
      type: Sequelize.TEXT,
      allowNull: true,
    })
    await queryInterface.changeColumn('gitlab_merge_requests', 'description', {
      type: Sequelize.TEXT,
      allowNull: true,
    })
    await queryInterface.changeColumn('gitlab_merge_requests', 'title', {
      type: Sequelize.TEXT,
      allowNull: true,
    })
  },

  down: async (queryInterface, Sequelize) => {
    
  }
};