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
class Jira extends Scheduler<void> {
  name(): string {
    return 'jira';
  }

  version(): number {
    return 1;
  }

  async migrateDown(currentVersion: string): Promise<void> {
    console.info(currentVersion);
    return;
  }

  async migrateUp(pluginPrev: string, oldVersion: string): Promise<string> {
    console.info(pluginPrev, oldVersion);
    return 'hx8f23r1';
  }

  async execute(options: JiraOptions): Promise<void> {
    //TODO: Add jira collector and enrichment
    console.info('Excute Jira', options);
    return;
  }
}

export default Jira;
