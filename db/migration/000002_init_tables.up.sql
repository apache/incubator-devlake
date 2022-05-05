BEGIN;

-- This line is required
USE lake_test;

-- MySQL dump 10.13  Distrib 8.0.22, for macos10.15 (x86_64)
--
-- Host: localhost    Database: lake
-- ------------------------------------------------------
-- Server version	8.0.26

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ae_commits`
--

DROP TABLE IF EXISTS `ae_commits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ae_commits` (
  `hex_sha` varchar(191) NOT NULL,
  `analysis_id` longtext,
  `author_email` longtext,
  `dev_eq` bigint DEFAULT NULL,
  `ae_project_id` bigint DEFAULT NULL,
  PRIMARY KEY (`hex_sha`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ae_projects`
--

DROP TABLE IF EXISTS `ae_projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ae_projects` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `git_url` longtext COMMENT 'url of the repo in github',
  `priority` bigint DEFAULT NULL,
  `ae_create_time` datetime(3) DEFAULT NULL,
  `ae_update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `boards`
--

DROP TABLE IF EXISTS `boards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `boards` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `url` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `builds`
--

DROP TABLE IF EXISTS `builds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `builds` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `job_id` varchar(191) DEFAULT NULL,
  `name` longtext,
  `commit_sha` longtext,
  `duration_sec` bigint unsigned DEFAULT NULL,
  `status` longtext,
  `started_date` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_builds_job_id` (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `changelogs`
--

DROP TABLE IF EXISTS `changelogs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `changelogs` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `issue_id` varchar(191) DEFAULT NULL,
  `author_id` longtext,
  `author_name` longtext,
  `field_id` longtext,
  `field_name` longtext,
  `from` longtext,
  `to` longtext,
  `created_date` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_changelogs_issue_id` (`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `commits`
--

DROP TABLE IF EXISTS `commits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `commits` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `repo_id` bigint unsigned DEFAULT NULL COMMENT 'References the repo the commit belongs to.',
  `sha` longtext COMMENT 'commit hash',
  `additions` bigint DEFAULT NULL COMMENT 'Added lines of code',
  `deletions` bigint DEFAULT NULL COMMENT 'Deleted lines of code',
  `dev_eq` bigint DEFAULT NULL COMMENT 'Merico developer equivalent from analysis engine',
  `message` longtext,
  `author_name` longtext,
  `author_email` longtext,
  `authored_date` datetime(3) DEFAULT NULL,
  `committer_name` longtext,
  `committer_email` longtext,
  `committed_date` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_commits_repo_id` (`repo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_commits`
--

DROP TABLE IF EXISTS `github_commits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_commits` (
  `sha` varchar(191) NOT NULL,
  `repository_id` bigint DEFAULT NULL,
  `author_name` longtext,
  `author_email` longtext,
  `authored_date` datetime(3) DEFAULT NULL,
  `committer_name` longtext,
  `committer_email` longtext,
  `committed_date` datetime(3) DEFAULT NULL,
  `message` longtext,
  `url` longtext,
  `additions` bigint DEFAULT NULL COMMENT 'Added lines of code',
  `deletions` bigint DEFAULT NULL COMMENT 'Deleted lines of code',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`sha`),
  KEY `idx_github_commits_repository_id` (`repository_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_issue_comments`
--

DROP TABLE IF EXISTS `github_issue_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_issue_comments` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `issue_id` bigint DEFAULT NULL COMMENT 'References the Pull Request',
  `body` longtext,
  `author_username` longtext,
  `github_created_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`),
  KEY `idx_github_issue_comments_issue_id` (`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_issue_events`
--

DROP TABLE IF EXISTS `github_issue_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_issue_events` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `issue_id` bigint DEFAULT NULL COMMENT 'References the Pull Request',
  `type` longtext COMMENT 'Events that can occur to an issue, ex. assigned, closed, labeled, etc.',
  `author_username` longtext,
  `github_created_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`),
  KEY `idx_github_issue_events_issue_id` (`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_issue_label_issues`
--

DROP TABLE IF EXISTS `github_issue_label_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_issue_label_issues` (
  `issue_label_id` bigint NOT NULL,
  `issue_id` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`issue_label_id`,`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_issue_labels`
--

DROP TABLE IF EXISTS `github_issue_labels`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_issue_labels` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `description` longtext,
  `color` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_issues`
--

DROP TABLE IF EXISTS `github_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_issues` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `number` bigint DEFAULT NULL COMMENT 'Used in API requests ex. api/repo/1/issue/<THIS_NUMBER>',
  `state` longtext,
  `title` longtext,
  `body` longtext,
  `priority` longtext,
  `type` longtext,
  `status` longtext,
  `assignee` longtext,
  `lead_time_minutes` bigint unsigned DEFAULT NULL,
  `closed_at` datetime(3) DEFAULT NULL,
  `github_created_at` datetime(3) DEFAULT NULL,
  `github_updated_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`),
  KEY `idx_github_issues_number` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_pull_request_comments`
--

DROP TABLE IF EXISTS `github_pull_request_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_pull_request_comments` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `pull_request_id` bigint DEFAULT NULL,
  `body` longtext,
  `author_username` longtext,
  `github_created_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`),
  KEY `idx_github_pull_request_comments_pull_request_id` (`pull_request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_pull_request_commit_pull_requests`
--

DROP TABLE IF EXISTS `github_pull_request_commit_pull_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_pull_request_commit_pull_requests` (
  `pull_request_commit_sha` varchar(191) NOT NULL,
  `pull_request_id` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`pull_request_commit_sha`,`pull_request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_pull_request_commits`
--

DROP TABLE IF EXISTS `github_pull_request_commits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_pull_request_commits` (
  `sha` varchar(191) NOT NULL,
  `pull_request_id` bigint DEFAULT NULL,
  `author_name` longtext,
  `author_email` longtext,
  `authored_date` datetime(3) DEFAULT NULL,
  `committer_name` longtext,
  `committer_email` longtext,
  `committed_date` datetime(3) DEFAULT NULL,
  `message` longtext,
  `url` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`sha`),
  KEY `idx_github_pull_request_commits_pull_request_id` (`pull_request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_pull_requests`
--

DROP TABLE IF EXISTS `github_pull_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_pull_requests` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `repository_id` bigint DEFAULT NULL,
  `number` bigint DEFAULT NULL,
  `state` longtext,
  `title` longtext,
  `github_created_at` datetime(3) DEFAULT NULL,
  `closed_at` datetime(3) DEFAULT NULL,
  `additions` bigint DEFAULT NULL,
  `deletions` bigint DEFAULT NULL,
  `comments` bigint DEFAULT NULL,
  `commits` bigint DEFAULT NULL,
  `review_comments` bigint DEFAULT NULL,
  `merged` tinyint(1) DEFAULT NULL,
  `merged_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`),
  KEY `idx_github_pull_requests_repository_id` (`repository_id`),
  KEY `idx_github_pull_requests_number` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_repositories`
--

DROP TABLE IF EXISTS `github_repositories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_repositories` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `html_url` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `github_reviewers`
--

DROP TABLE IF EXISTS `github_reviewers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `github_reviewers` (
  `github_id` bigint NOT NULL AUTO_INCREMENT,
  `login` longtext,
  `pull_request_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`github_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_commits`
--

DROP TABLE IF EXISTS `gitlab_commits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_commits` (
  `gitlab_id` varchar(191) NOT NULL,
  `project_id` bigint DEFAULT NULL,
  `title` longtext,
  `message` longtext,
  `short_id` longtext,
  `author_name` longtext,
  `author_email` longtext,
  `authored_date` datetime(3) DEFAULT NULL,
  `committer_name` longtext,
  `committer_email` longtext,
  `committed_date` datetime(3) DEFAULT NULL,
  `web_url` longtext,
  `additions` bigint DEFAULT NULL COMMENT 'Added lines of code',
  `deletions` bigint DEFAULT NULL COMMENT 'Deleted lines of code',
  `total` bigint DEFAULT NULL COMMENT 'Sum of added/deleted lines of code',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`gitlab_id`),
  KEY `idx_gitlab_commits_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_merge_request_commit_merge_requests`
--

DROP TABLE IF EXISTS `gitlab_merge_request_commit_merge_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_merge_request_commit_merge_requests` (
  `merge_request_commit_id` varchar(191) NOT NULL,
  `merge_request_id` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`merge_request_commit_id`,`merge_request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_merge_request_commits`
--

DROP TABLE IF EXISTS `gitlab_merge_request_commits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_merge_request_commits` (
  `commit_id` varchar(191) NOT NULL,
  `title` longtext,
  `message` longtext,
  `short_id` longtext,
  `author_name` longtext,
  `author_email` longtext,
  `authored_date` datetime(3) DEFAULT NULL,
  `committer_name` longtext,
  `committer_email` longtext,
  `committed_date` datetime(3) DEFAULT NULL,
  `web_url` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`commit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_merge_request_notes`
--

DROP TABLE IF EXISTS `gitlab_merge_request_notes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_merge_request_notes` (
  `gitlab_id` bigint NOT NULL AUTO_INCREMENT,
  `merge_request_id` bigint DEFAULT NULL,
  `merge_request_iid` bigint DEFAULT NULL COMMENT 'Used in API requests ex. /api/merge_requests/<THIS_IID>',
  `noteable_type` longtext,
  `author_username` longtext,
  `body` longtext,
  `gitlab_created_at` datetime(3) DEFAULT NULL,
  `confidential` tinyint(1) DEFAULT NULL,
  `resolvable` tinyint(1) DEFAULT NULL COMMENT 'Is or is not review comment',
  `system` tinyint(1) DEFAULT NULL COMMENT 'Is or is not auto-generated vs. human generated',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`gitlab_id`),
  KEY `idx_gitlab_merge_request_notes_merge_request_id` (`merge_request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_merge_requests`
--

DROP TABLE IF EXISTS `gitlab_merge_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_merge_requests` (
  `gitlab_id` bigint NOT NULL AUTO_INCREMENT,
  `iid` bigint DEFAULT NULL,
  `project_id` bigint DEFAULT NULL,
  `state` longtext,
  `title` longtext,
  `web_url` longtext,
  `user_notes_count` bigint DEFAULT NULL,
  `work_in_progress` tinyint(1) DEFAULT NULL,
  `source_branch` longtext,
  `merged_at` datetime(3) DEFAULT NULL,
  `gitlab_created_at` datetime(3) DEFAULT NULL,
  `closed_at` datetime(3) DEFAULT NULL,
  `merged_by_username` longtext,
  `description` longtext,
  `author_username` longtext,
  `first_comment_time` datetime(3) DEFAULT NULL COMMENT 'Time when the first comment occurred',
  `review_rounds` bigint DEFAULT NULL COMMENT 'How many rounds of review this MR went through',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`gitlab_id`),
  KEY `idx_gitlab_merge_requests_iid` (`iid`),
  KEY `idx_gitlab_merge_requests_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_pipelines`
--

DROP TABLE IF EXISTS `gitlab_pipelines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_pipelines` (
  `gitlab_id` bigint NOT NULL AUTO_INCREMENT,
  `project_id` bigint DEFAULT NULL,
  `gitlab_created_at` datetime(3) DEFAULT NULL,
  `status` longtext,
  `ref` longtext,
  `sha` longtext,
  `web_url` longtext,
  `duration` bigint DEFAULT NULL,
  `started_at` datetime(3) DEFAULT NULL,
  `finished_at` datetime(3) DEFAULT NULL,
  `coverage` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`gitlab_id`),
  KEY `idx_gitlab_pipelines_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_projects`
--

DROP TABLE IF EXISTS `gitlab_projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_projects` (
  `gitlab_id` bigint NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `path_with_namespace` longtext,
  `web_url` longtext,
  `visibility` longtext,
  `open_issues_count` bigint DEFAULT NULL,
  `star_count` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`gitlab_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gitlab_reviewers`
--

DROP TABLE IF EXISTS `gitlab_reviewers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gitlab_reviewers` (
  `gitlab_id` bigint NOT NULL AUTO_INCREMENT,
  `merge_request_id` bigint DEFAULT NULL,
  `project_id` bigint DEFAULT NULL,
  `name` longtext,
  `username` longtext,
  `state` longtext,
  `avatar_url` longtext,
  `web_url` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`gitlab_id`),
  KEY `idx_gitlab_reviewers_merge_request_id` (`merge_request_id`),
  KEY `idx_gitlab_reviewers_project_id` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `issues`
--

DROP TABLE IF EXISTS `issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `issues` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `board_id` varchar(191) DEFAULT NULL,
  `url` longtext,
  `key` longtext,
  `title` longtext,
  `summary` longtext,
  `epic_key` longtext,
  `type` longtext,
  `status` longtext,
  `story_point` bigint unsigned DEFAULT NULL,
  `original_estimate_minutes` bigint DEFAULT NULL,
  `aggregate_estimate_minutes` bigint DEFAULT NULL,
  `remaining_estimate_minutes` bigint DEFAULT NULL,
  `creator_id` longtext,
  `assignee_id` longtext,
  `resolution_date` datetime(3) DEFAULT NULL,
  `priority` longtext,
  `parent_id` longtext,
  `sprint_id` longtext,
  `created_date` datetime(3) DEFAULT NULL,
  `updated_date` datetime(3) DEFAULT NULL,
  `spent_minutes` bigint DEFAULT NULL,
  `lead_time_minutes` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_issues_board_id` (`board_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jenkins_builds`
--

DROP TABLE IF EXISTS `jenkins_builds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jenkins_builds` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `duration` double DEFAULT NULL,
  `display_name` longtext,
  `estimated_duration` double DEFAULT NULL,
  `number` bigint DEFAULT NULL,
  `result` longtext,
  `timestamp` bigint DEFAULT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `job_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jenkins_jobs`
--

DROP TABLE IF EXISTS `jenkins_jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jenkins_jobs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `class` longtext,
  `color` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_board_gitlab_projects`
--

DROP TABLE IF EXISTS `jira_board_gitlab_projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_board_gitlab_projects` (
  `jira_board_id` bigint unsigned NOT NULL,
  `gitlab_project_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`jira_board_id`,`gitlab_project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_board_issues`
--

DROP TABLE IF EXISTS `jira_board_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_board_issues` (
  `connection_id` bigint unsigned NOT NULL,
  `board_id` bigint unsigned NOT NULL,
  `issue_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`connection_id`,`board_id`,`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_board_sprints`
--

DROP TABLE IF EXISTS `jira_board_sprints`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_board_sprints` (
  `connection_id` bigint unsigned NOT NULL,
  `board_id` bigint unsigned NOT NULL,
  `sprint_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`connection_id`,`board_id`,`sprint_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_boards`
--

DROP TABLE IF EXISTS `jira_boards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_boards` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `board_id` bigint unsigned NOT NULL,
  `project_id` bigint unsigned DEFAULT NULL,
  `name` longtext,
  `self` longtext,
  `type` longtext,
  PRIMARY KEY (`connection_id`,`board_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_changelog_items`
--

DROP TABLE IF EXISTS `jira_changelog_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_changelog_items` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `changelog_id` bigint unsigned NOT NULL,
  `field` varchar(191) NOT NULL,
  `field_type` longtext,
  `field_id` longtext,
  `from` longtext,
  `from_string` longtext,
  `to` longtext,
  `to_string` longtext,
  PRIMARY KEY (`connection_id`,`changelog_id`,`field`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_changelogs`
--

DROP TABLE IF EXISTS `jira_changelogs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_changelogs` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `changelog_id` bigint unsigned NOT NULL,
  `issue_id` bigint unsigned DEFAULT NULL,
  `author_account_id` longtext,
  `author_display_name` longtext,
  `author_active` tinyint(1) DEFAULT NULL,
  `created` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`connection_id`,`changelog_id`),
  KEY `idx_jira_changelogs_issue_id` (`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_issue_status_mappings`
--

DROP TABLE IF EXISTS `jira_issue_status_mappings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_issue_status_mappings` (
  `connection_id` bigint unsigned NOT NULL,
  `user_type` varchar(50) NOT NULL,
  `user_status` varchar(50) NOT NULL,
  `standard_status` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`connection_id`,`user_type`,`user_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_issue_type_mappings`
--

DROP TABLE IF EXISTS `jira_issue_type_mappings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_issue_type_mappings` (
  `connection_id` bigint unsigned NOT NULL,
  `user_type` varchar(50) NOT NULL,
  `standard_type` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`connection_id`,`user_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_issues`
--

DROP TABLE IF EXISTS `jira_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_issues` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `issue_id` bigint unsigned NOT NULL,
  `project_id` bigint unsigned DEFAULT NULL,
  `self` longtext,
  `key` longtext,
  `summary` longtext,
  `type` longtext,
  `epic_key` longtext,
  `status_name` longtext,
  `status_key` longtext,
  `status_category` longtext,
  `story_point` double DEFAULT NULL,
  `original_estimate_minutes` bigint DEFAULT NULL,
  `aggregate_estimate_minutes` bigint DEFAULT NULL,
  `remaining_estimate_minutes` bigint DEFAULT NULL,
  `creator_account_id` longtext,
  `creator_account_type` longtext,
  `creator_display_name` longtext,
  `assignee_account_id` longtext COMMENT 'latest assignee',
  `assignee_account_type` longtext,
  `assignee_display_name` longtext,
  `priority_id` bigint unsigned DEFAULT NULL,
  `priority_name` longtext,
  `parent_id` bigint unsigned DEFAULT NULL,
  `parent_key` longtext,
  `sprint_id` bigint unsigned DEFAULT NULL,
  `sprint_name` longtext,
  `resolution_date` datetime(3) DEFAULT NULL,
  `created` datetime(3) DEFAULT NULL,
  `updated` datetime(3) DEFAULT NULL,
  `spent_minutes` bigint DEFAULT NULL,
  `lead_time_minutes` bigint unsigned DEFAULT NULL,
  `std_story_point` bigint unsigned DEFAULT NULL,
  `std_type` longtext,
  `std_status` longtext,
  `all_fields` json DEFAULT NULL,
  `changelog_updated` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`connection_id`,`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_projects`
--

DROP TABLE IF EXISTS `jira_projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_projects` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `id` varchar(191) NOT NULL,
  `key` longtext,
  `name` longtext,
  PRIMARY KEY (`connection_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_connections`
--

DROP TABLE IF EXISTS `jira_connections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_connections` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `endpoint` longtext,
  `basic_auth_encoded` longtext,
  `epic_key_field` varchar(50) DEFAULT NULL,
  `story_point_field` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_jira_connections_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_sprint_issues`
--

DROP TABLE IF EXISTS `jira_sprint_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_sprint_issues` (
  `connection_id` bigint unsigned NOT NULL,
  `sprint_id` bigint unsigned NOT NULL,
  `issue_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`connection_id`,`sprint_id`,`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_sprints`
--

DROP TABLE IF EXISTS `jira_sprints`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_sprints` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `sprint_id` bigint unsigned NOT NULL,
  `self` longtext,
  `state` longtext,
  `name` longtext,
  `start_date` datetime(3) DEFAULT NULL,
  `end_date` datetime(3) DEFAULT NULL,
  `complete_date` datetime(3) DEFAULT NULL,
  `origin_board_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`connection_id`,`sprint_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_users`
--

DROP TABLE IF EXISTS `jira_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_users` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `account_id` varchar(191) NOT NULL,
  `account_type` longtext,
  `name` longtext,
  `email` longtext,
  `avatar_url` longtext,
  `timezone` longtext,
  PRIMARY KEY (`connection_id`,`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jira_worklogs`
--

DROP TABLE IF EXISTS `jira_worklogs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jira_worklogs` (
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `connection_id` bigint unsigned NOT NULL,
  `issue_id` bigint unsigned NOT NULL,
  `worklog_id` varchar(191) NOT NULL,
  `author_id` longtext,
  `update_author_id` longtext,
  `time_spent` longtext,
  `time_spent_seconds` bigint DEFAULT NULL,
  `updated` datetime(3) DEFAULT NULL,
  `started` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`connection_id`,`issue_id`,`worklog_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `jobs`
--

DROP TABLE IF EXISTS `jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jobs` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `notes`
--

DROP TABLE IF EXISTS `notes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notes` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `pr_id` bigint unsigned DEFAULT NULL COMMENT 'References the pull request for this note',
  `type` longtext,
  `author` longtext,
  `body` longtext,
  `resolvable` tinyint(1) DEFAULT NULL COMMENT 'Is or is not a review comment',
  `system` tinyint(1) DEFAULT NULL COMMENT 'Is or is not auto-generated vs. human generated',
  `created_date` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_notes_pr_id` (`pr_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `type` longtext,
  `endpoint` longtext,
  `nonce` longtext,
  `response_code` bigint DEFAULT NULL,
  `response` longtext,
  `data` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pipelines`
--

DROP TABLE IF EXISTS `pipelines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pipelines` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) DEFAULT NULL,
  `tasks` json DEFAULT NULL,
  `total_tasks` bigint DEFAULT NULL,
  `finished_tasks` bigint DEFAULT NULL,
  `began_at` datetime(3) DEFAULT NULL,
  `finished_at` datetime(3) DEFAULT NULL,
  `status` longtext,
  `message` longtext,
  `spent_seconds` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pipelines_name` (`name`),
  KEY `idx_pipelines_finished_at` (`finished_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `prs`
--

DROP TABLE IF EXISTS `prs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prs` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `repo_id` bigint unsigned DEFAULT NULL,
  `state` longtext COMMENT 'open/closed or other',
  `title` longtext,
  `url` longtext,
  `created_date` datetime(3) DEFAULT NULL,
  `merged_date` datetime(3) DEFAULT NULL,
  `closed_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_prs_repo_id` (`repo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `repos`
--

DROP TABLE IF EXISTS `repos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `repos` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `url` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sprint_issues`
--

DROP TABLE IF EXISTS `sprint_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sprint_issues` (
  `sprint_id` varchar(191) NOT NULL,
  `issue_id` varchar(191) NOT NULL,
  PRIMARY KEY (`sprint_id`,`issue_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sprints`
--

DROP TABLE IF EXISTS `sprints`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sprints` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `board_id` varchar(191) DEFAULT NULL,
  `url` longtext,
  `state` longtext,
  `name` longtext,
  `start_date` datetime(3) DEFAULT NULL,
  `end_date` datetime(3) DEFAULT NULL,
  `complete_date` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sprints_board_id` (`board_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `plugin` varchar(191) DEFAULT NULL,
  `options` json DEFAULT NULL,
  `status` longtext,
  `message` longtext,
  `progress` float DEFAULT NULL,
  `pipeline_id` bigint unsigned DEFAULT NULL,
  `pipeline_row` bigint DEFAULT NULL,
  `pipeline_col` bigint DEFAULT NULL,
  `began_at` datetime(3) DEFAULT NULL,
  `finished_at` datetime(3) DEFAULT NULL,
  `spent_seconds` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tasks_pipeline_id` (`pipeline_id`),
  KEY `idx_tasks_finished_at` (`finished_at`),
  KEY `idx_tasks_plugin` (`plugin`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `email` longtext,
  `avatar_url` longtext,
  `timezone` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `worklogs`
--

DROP TABLE IF EXISTS `worklogs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `worklogs` (
  `id` varchar(255) NOT NULL COMMENT 'This key is generated based on details from the original plugin',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `issue_id` varchar(191) DEFAULT NULL,
  `board_id` varchar(191) DEFAULT NULL,
  `author_id` longtext,
  `update_author_id` longtext,
  `time_spent` longtext,
  `time_spent_seconds` bigint DEFAULT NULL,
  `updated` datetime(3) DEFAULT NULL,
  `started` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_worklogs_issue_id` (`issue_id`),
  KEY `idx_worklogs_board_id` (`board_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping events for database 'lake'
--

--
-- Dumping routines for database 'lake'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-12-17 11:30:44


SELECT now();

COMMIT;
