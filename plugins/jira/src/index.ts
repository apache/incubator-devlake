import { Injectable, Scope } from '@nestjs/common';
import Scheduler from 'plugins/core/src/scheculer';

export type JiraCollector =
  | 'ISSUE'
  | 'CHANGELOG'
  | 'COMMENTS'
  | 'REMOTELINK'
  | 'BOARD';

export type JiraOptions = {
  collectors: JiraCollector[];
};

@Injectable({
  scope: Scope.TRANSIENT,
})
class Jira extends Scheduler<void, JiraOptions> {
  async execute(options: JiraOptions): Promise<void> {
    //TODO: Add jira collector and enrichment
    console.info('Excute Jira', options);
    return;
  }
}

export default Jira;
