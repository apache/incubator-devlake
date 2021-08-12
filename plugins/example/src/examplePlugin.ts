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
class ExamplePlugin implements Plugin {
  constructor(
    @InjectConnection('exampleModuleDb')
    private connection: Connection,
    private issueCollector: IssueCollector,
  ) {}

  version(): number {
    return 1;
  }

  async migrateDown(currentVersion: string): Promise<void> {
    console.info(currentVersion);
    return;
  }

  async migrateUp(pluginPrev: string): Promise<void> {
    // 如果不在意migrate，直接用synchronize自动更新表结构
    // If you don't care about migrate, use synchronize to automatically update the table structure
    // await this.connection.synchronize();

    // 如果希望执行完整的typeorm migrate
    // If you want to perform a complete typeorm migrate
    await this.connection.runMigrations({ transaction: 'each' });

    // 或者按照其他任何插件作者换的方式
    // Or follow any other way as plugin author want
  }

  async execute(options: JiraOptions): Promise<void> {
    await this.issueCollector.collectData({}, null);
    //TODO: Add jira collector and enrichment
    console.info('Execute Example', options);
    return;
  }
}

export default ExamplePlugin;
