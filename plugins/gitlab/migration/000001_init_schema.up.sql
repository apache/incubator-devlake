BEGIN;

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

COMMIT;
