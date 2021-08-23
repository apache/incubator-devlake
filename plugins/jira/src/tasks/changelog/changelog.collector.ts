import { Inject } from '@nestjs/common';
import Exports from 'plugins/core/src/exports.decorator';
import Task from 'plugins/core/src/task.interface';
import { Repository } from 'typeorm';
import ChangelogEntity from './changelog.entity';

export type JiraSource = {
  host: string;
  username: string;
  token: string;
};

@Exports(ChangelogEntity)
export default class ChangelogCollector implements Task {
  // @Inject(IssueEntity) private IssueRepository: Repository<IssueEntity>;

  name(): string {
    return 'IssueChangelog';
  }

  async execute(source: JiraSource): Promise<void> {
    //TODO: do collector
    console.info('Excute Issue Changelog Collector', source);
    return;
  }
}