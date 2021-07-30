"use strict";
module.exports = {
  up: async (queryInterface, Sequelize) => {
    await queryInterface.createTable("gitlab_merge_requests", {
      id: {
        primaryKey: true,
        type: Sequelize.INTEGER,
      },
      title: {
        type: Sequelize.STRING,
      },
      projectId: {
        type: Sequelize.INTEGER,
        field: "project_id",
      },
      numberOfReviewers: {
        type: Sequelize.INTEGER,
        field: "number_of_reviewers",
      },
      state: {
        type: Sequelize.STRING
      },
      title: {
        type: Sequelize.STRING
      },
      webUrl: {
        type: Sequelize.STRING,
        field: "web_url",
      },
      userNotesCount: {
        type: Sequelize.INTEGER,
        field: "user_notes_count",
      },
      workInProgress: {
        type: Sequelize.BOOLEAN,
        field: "work_in_progress",
      },
      sourceBranch: {
        type: Sequelize.STRING,
        field: "source_branch",
      },
      mergedAt: {
        type: Sequelize.DATE,
        field: "merged_at",
      },
      gitlabCreatedAt: {
        type: Sequelize.DATE,
        field: "gitlab_created_at",
      },
      closedAt: {
        type: Sequelize.DATE,
        field: "closed_at",
      },
      mergedByUsername: {
        type: Sequelize.STRING,
        field: "merged_by_username",
      },
      description: {
        type: Sequelize.STRING
      },
      authorUsername: {
        type: Sequelize.STRING,
        field: "author_username",
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
    await queryInterface.dropTable("gitlab_merge_requests");
  },
};
