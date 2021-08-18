import { ClassProvider, DynamicModule } from '@nestjs/common';
import Jira from 'plugins/jira/src';

const JiraProvider: ClassProvider<Jira> = {
  provide: 'Jira',
  useClass: Jira,
};

export const providers = [JiraProvider];

export default class PluginModule {
  static async forRootAsync(): Promise<DynamicModule> {
    return {
      module: PluginModule,
      providers,
      exports: [...providers],
    };
  }
}
