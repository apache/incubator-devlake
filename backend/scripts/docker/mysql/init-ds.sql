CREATE
DATABASE IF NOT EXISTS bitbucket;
alter
database bitbucket character set utf8 collate utf8_bin;
GRANT ALL PRIVILEGES ON bitbucket.* TO
'merico';

CREATE
DATABASE IF NOT EXISTS jira;
alter
database jira character set utf8 collate utf8_bin;
GRANT ALL PRIVILEGES ON jira.* TO
'merico';