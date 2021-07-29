'use strict'
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable('gitlab_commits', {
      id: {
        primaryKey: true,
        type: Sequelize.STRING
      },
      shortId: {
        type: Sequelize.STRING,
        field: 'short_id'
      },
      title: {
        type: Sequelize.STRING
      },
      message: {
        type: Sequelize.STRING
      },
      authorName: {
        type: Sequelize.STRING,
        field: 'author_name'
      },
      authorEmail: {
        type: Sequelize.STRING,
        field: 'author_email'
      },
      authoredDate: {
        type: Sequelize.STRING,
        field: 'authored_date'
      },
      committerName: {
        type: Sequelize.STRING,
        field: 'committer_name'
      },
      committerEmail: {
        type: Sequelize.STRING,
        field: 'committer_email'
      },
      committedDate: {
        type: Sequelize.STRING,
        field: 'committed_date'
      },
      webUrl: {
        type: Sequelize.STRING,
        field: 'web_url'
      },
      additions: {
        type: Sequelize.INTEGER,
      },
      deletions: {
        type: Sequelize.INTEGER,
      },
      total: {
        type: Sequelize.INTEGER,
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
    await queryInterface.dropTable('gitlab_commits')
  }
}
