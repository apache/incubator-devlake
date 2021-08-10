import { Injectable } from '@nestjs/common';
import Plugin from 'plugins/core/src/plugin';

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
class Jira implements Plugin {
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

  async register(): Promise<void> {
    //TODO: Add jira collector and enrichment
    console.info('Excute Jira');
    return;
  }

  async unregister(): Promise<void> {
    return;
  }
}

export default Jira;
