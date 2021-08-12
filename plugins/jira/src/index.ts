import { Injectable, Scope } from '@nestjs/common';
import Collector, {
  CollectorMap,
  COLLECTORS_METADATA,
} from 'plugins/core/src/collector.decorate';
import Scheduler from 'plugins/core/src/scheculer';
import IssueCollector from './collector/issue';

export enum JiraCollector {
  ISSUE = 'ISSUE',
  CHANGELOG = 'CHANGELOG',
  COMMENTS = 'COMMENTS',
  REMOTELINK = 'REMOTELINK',
  BOARD = 'BOARD',
}

export type JiraOptions = {
  collectors: JiraCollector[];
};

@Injectable({
  scope: Scope.TRANSIENT,
})
@Collector({
  [JiraCollector.ISSUE]: IssueCollector,
})
class Jira extends Scheduler<void, JiraOptions> {
  async execute(options: JiraOptions): Promise<void> {
    //TODO: Add jira collector and enrichment
    const collectorMaps: CollectorMap = Reflect.getMetadata(
      COLLECTORS_METADATA,
      Jira,
    );
    const collector = new collectorMaps[JiraCollector.ISSUE]();
    await collector.execute(options)
    return;
  }
}

export default Jira;
