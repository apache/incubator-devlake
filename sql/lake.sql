/* count of issue type  */
WITH
    jira_types as (SELECT DISTINCT `type` FROM jira_issues)
SELECT t.type,(SELECT COUNT(*) FROM jira_issues WHERE `type`=t.type) FROM jira_types t

/* Lead-time of Epic */