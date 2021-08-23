import { Inject } from '@nestjs/common';
import Exports from 'plugins/core/src/exports.decorator';
import Imports from 'plugins/core/src/imports.decorator';
import Task from 'plugins/core/src/task.interface';
import { Repository } from 'typeorm';
import ChangelogEntity from '../changelog/changelog.entity';
import IssueLeadTimeEntity from './leadtime.entity';

export type JiraSource = {
  host: string;
  username: string;
  token: string;
};

@Imports([ChangelogEntity])
@Exports(IssueLeadTimeEntity)
export default class LeadTimeEnricher implements Task {
  // @Inject(IssueEntity) private IssueRepository: Repository<IssueEntity>;
  // @Inject(IssueLeadTimeEntity) private IssueLeadTimeRepository: Repository<IssueLeadTimeEntity>;

  name(): string {
    return 'JiraIssueLeadTime';
  }

  async execute(source: JiraSource): Promise<void> {
    //TODO: do enrichment
    console.info('Excute Jira Lead Time Enricher', source);
    await new Promise((resolve) => setTimeout(() => resolve(true), 100));
    return;
  }
}