BEGIN;

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

COMMIT;
