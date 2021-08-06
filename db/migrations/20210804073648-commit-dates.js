'use strict'

module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.changeColumn('gitlab_commits', 'authored_date', {
      type: 'timestamptz using  cast(authored_date as timestamptz)'
    })
    await queryInterface.changeColumn('gitlab_commits', 'authored_date', {
      type: Sequelize.DATE
    })
    await queryInterface.changeColumn('gitlab_commits', 'committed_date', {
      type: 'timestamptz using cast(committed_date as timestamptz)'
    })
    await queryInterface.changeColumn('gitlab_commits', 'committed_date', {
      type: Sequelize.DATE
    })
  },

  down: async (queryInterface, Sequelize) => {
    await queryInterface.changeColumn('gitlab_commits', 'authored_date', {
      type: Sequelize.STRING
    })
    await queryInterface.changeColumn('gitlab_commits', 'committed_date', {
      type: Sequelize.STRING
    })
  }
}
