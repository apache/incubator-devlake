import { Module } from '@nestjs/common';
import JiraPlugin from './JiraPlugin';
import { ScheduleModule } from '@nestjs/schedule';
import IssueCollector from './runners/IssueCollector';
import CustomTypeOrmModule from '../../../apps/rest/src/customTypeOrmModule';

@Module({
  imports: [
    CustomTypeOrmModule.forRootAsync('jiraModuleDb', {
      entityPrefix: 'plugin_jira_',
      synchronize: true,
      entitiesFunc: () => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        return require.context('./entities', true, /\.ts/);
      },
    }),
    ScheduleModule.forRoot(),
  ],
  providers: [JiraPlugin, IssueCollector],
  exports: [JiraPlugin],
})
class JiraModule {}

export default JiraModule;
