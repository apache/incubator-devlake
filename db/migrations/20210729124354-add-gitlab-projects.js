'use strict'
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('gitlab_projects', {
      name: {
        type: Sequelize.STRING
      },
      id: {
        type: Sequelize.INTEGER,
        primaryKey: true,
      },
      pathWithNamespace: {
        type: Sequelize.STRING,
        field: 'path_with_namespace'
      },
      webUrl: {
        type: Sequelize.STRING,
        field: 'web_url'
      },
      visibility: {
        type: Sequelize.STRING,
      },
      openIssuesCount: {
        type: Sequelize.INTEGER,
        field: 'open_issues_count'
      },
      starCount: {
        type: Sequelize.INTEGER,
        field: 'star_count'
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
    await queryInterface.dropTable('gitlab_projects')
  }
}
