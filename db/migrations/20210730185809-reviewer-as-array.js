'use strict'

module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.addColumn('gitlab_merge_requests', 'reviewers', {
      type: Sequelize.ARRAY(Sequelize.STRING),
      allowNull: true
    })
  },

  down: async (queryInterface, Sequelize) => {

  }
}
