import { Injectable } from '@nestjs/common';
import Plugin from '../../core/src/plugin';
import { Cron, CronExpression } from '@nestjs/schedule';
import IssueCollector from './runners/IssueCollector';
import { Connection } from 'typeorm';
import { InjectConnection } from '@nestjs/typeorm';

export type JiraCollector =
  | 'ISSUE'
  | 'CHANGELOG'
  | 'COMMENTS'
  | 'REMOTELINK'
  | 'BOARD';

export type JiraOptions = {
  collectors: JiraCollector[];
};

@Injectable()
class JiraPlugin implements Plugin {
  constructor(
    @InjectConnection('jiraModuleDb')
    private connection: Connection,
    private issueCollector: IssueCollector,
  ) {}

  version(): number {
    return 1;
  }

  async migrateDown(): Promise<void> {
    throw new Error('not support');
  }

  async migrateUp(pluginPrev: string): Promise<void> {
    // TODO
    // await this.connection.runMigrations({ transaction: 'each' });
    await this.connection.synchronize();
    console.info(pluginPrev);
  }

  async execute(options: JiraOptions): Promise<void> {
    await this.issueCollector.collectData({}, null);
    //TODO: Add jira collector and enrichment
    console.info('Execute Jira', options);
    return;
  }
}

export default JiraPlugin;
