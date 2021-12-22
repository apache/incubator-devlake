BEGIN;
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

COMMIT;
