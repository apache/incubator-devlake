BEGIN;

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
