"use strict";
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable("gitlab_merge_request_notes", {
      id: {
        primaryKey: true,
        type: Sequelize.INTEGER,
      },
      noteableId: {
        type: Sequelize.INTEGER,
        field: 'noteable_id'
      },
      noteableIid: {
        type: Sequelize.INTEGER,
        field: 'noteable_iid'
      },
      authorUsername: {
        type: Sequelize.STRING,
        field: 'author_username'
      },
      body: {
        type: Sequelize.TEXT
      },
      gitlabCreatedAt: {
        type: Sequelize.STRING,
        field: 'gitlab_created_at'
      },
      noteableType: {
        type: Sequelize.STRING,
        field: 'noteable_type'
      },
      confidential: {
        type: Sequelize.STRING
      },
      createdAt: {
        allowNull: false,
        type: Sequelize.DATE,
        field: "created_at",
      },
      updatedAt: {
        allowNull: false,
        type: Sequelize.DATE,
        field: "updated_at",
      },
    });
  },
  down: async (queryInterface, Sequelize) => {
    await queryInterface.dropTable("gitlab_merge_request_notes");
  },
};
