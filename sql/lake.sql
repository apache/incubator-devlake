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

/* 开发环境 !!! 危险，删除今天的issue，重新测试拉取数据用。*/
DELETE FROM jira_issues 
WHERE 
  updated >= date(now())
  AND updated < DATE_ADD(date(now()),INTERVAL 1 DAY);



/* 迭代完成率 */
WITH
  sprint_type_issue_count as (
    SELECT 
      i.sprint_name,
      i.type,
      COUNT(1) as issue_count
    FROM jira_issues i
    GROUP BY i.sprint_name,i.type
  ),
  sprint_type_issue_done_count as (
    SELECT
      i.sprint_name,
      i.type,
      COUNT(1) as done_count
    FROM jira_issues i 
    WHERE i.status_name='已完成'
    GROUP BY i.sprint_name,i.type
  )
SELECT 
    ic.sprint_name,
    ic.type,
    ic.issue_count,
    idc.done_count,
    idc.done_count/ic.issue_count as done_ratio
FROM sprint_type_issue_count ic,
  sprint_type_issue_done_count idc
WHERE ic.type = idc.type 
  AND ic.sprint_name=idc.sprint_name
GROUP BY ic.sprint_name,ic.type;

/* 周完成率 */
WITH 
  jira_issues_weeks_count AS (
    SELECT 
      DATE_FORMAT(i.changelog_updated,'%Y,%u,1') AS weeks,
      COUNT(1) AS issue_count
    FROM jira_issues i
    GROUP BY DATE_FORMAT(i.changelog_updated,'%Y,%u,1')
  ),
  jira_issues_weeks_done_count AS (
    SELECT 
      DATE_FORMAT(i.changelog_updated,'%Y,%u,1') AS weeks,
      COUNT(1) AS done_count
    FROM jira_issues i
    WHERE i.status_name in ('已完成','已关闭')
    GROUP BY DATE_FORMAT(i.changelog_updated,'%Y,%u,1')
  )
SELECT 
  STR_TO_DATE(t1.weeks,'%Y,%u,%w') AS `time`,
  t1.weeks,
  t2.done_count,
  t1.issue_count,
  t2.done_count/t1.issue_count AS done_ratio
FROM jira_issues_weeks_count t1,
  jira_issues_weeks_done_count t2
WHERE t1.weeks=t2.weeks
GROUP BY t1.weeks
ORDER BY t1.weeks;


/* 周完成率 2 */
WITH 
  jira_issues_weeks AS (
    SELECT 
      *,
      STR_TO_DATE(DATE_FORMAT(i.changelog_updated,'%Y,%u,1'),'%Y,%u,%w') AS weeks_day
    FROM jira_issues i
  )
SELECT * FROM (
  WITH
    jira_issues_count AS (
      SELECT 
        t.weeks_day,
        COUNT(1) AS count_week
      FROM jira_issues_weeks t
      GROUP BY t.weeks_day
    ),
    jira_issues_count_done AS (
      SELECT 
        t.weeks_day,
        COUNT(1) AS count_done
      FROM jira_issues_weeks t
      WHERE t.status_name in ('已完成','已关闭')
      GROUP BY t.weeks_day
    )
  SELECT 
    t1.weeks_day,
    IFNULL(t2.count_done, 0) AS count_done,
    t1.count_week,
    IFNULL(t2.count_done, 0)/t1.count_week AS done_ratio
  FROM jira_issues_count t1
  LEFT JOIN jira_issues_count_done t2 
    ON t1.weeks_day=t2.weeks_day
  ORDER BY t1.weeks_day DESC
);