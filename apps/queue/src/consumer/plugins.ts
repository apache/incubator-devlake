import { DynamicModule } from '@nestjs/common';
import JiraModule from '../../../../plugins/jira/src';
import JiraPlugin from '../../../../plugins/jira/src/JiraPlugin';
import ExamplePlugin from '../../../../plugins/example/src/examplePlugin';
import ExampleModule from '../../../../plugins/example/src';

export const modules = [JiraModule, ExampleModule];
export const plugins = [JiraPlugin, ExamplePlugin];

export class PluginModule {
  static async forRootAsync(): Promise<DynamicModule> {
    return {
      module: PluginModule,
      imports: modules,
      exports: [...modules],
    };
  }
}
