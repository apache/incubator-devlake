import { Inject } from '@nestjs/common';
import Exports from 'plugins/core/src/exports.decorator';
import Imports from 'plugins/core/src/imports.decorator';
import Task from 'plugins/core/src/task.interface';
import { Repository } from 'typeorm';
import IssueEntity from '../issue/issue.entity';
import IssueLeadTimeEntity from './leadtime.entity';

export type JiraSource = {
  host: string;
  username: string;
  token: string;
};

@Imports([IssueEntity])
@Exports(IssueLeadTimeEntity)
export default class LeadTimeEnricher implements Task {
  @Inject(IssueEntity) private IssueRepository: Repository<IssueEntity>;
  @Inject(IssueLeadTimeEntity) private IssueLeadTimeRepository: Repository<IssueLeadTimeEntity>;

  name(): string {
    return 'JiraIssueLeadTime';
  }

  async execute(source: JiraSource): Promise<void> {
    //TODO: do collector
    console.info('Excute Jira Issue Collector', source);
    return;
  }
}