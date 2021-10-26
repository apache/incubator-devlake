/* count of issue type  */
WITH
    jira_types as (SELECT DISTINCT `type` FROM jira_issues)
SELECT t.type,(SELECT COUNT(*) FROM jira_issues WHERE `type`=t.type) FROM jira_types t;

/* Lead-time of Epic */
SELECT 
  i.key AS 'Jira Key',
  i.summary AS '项目概述',
  i.std_status AS '项目状态',
  i.lead_time DIV 1440 AS '需求交付周期',
  i.changelog_updated AS '最后更新时间'
FROM jira_issues i
WHERE 
  i.type = 'Epic'
ORDER BY i.updated_at DESC;

/* 迭代完成率 */
SELECT 
FROM jira_issues i
GROUP BY i.;

/* 开发环境 !!! 危险，删除今天的issue，重新测试拉取数据用。*/
DELETE FROM jira_issues 
WHERE 
  updated >= date(now())
  AND updated < DATE_ADD(date(now()),INTERVAL 1 DAY);