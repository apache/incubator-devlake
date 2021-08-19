import { ClassProvider, DynamicModule } from '@nestjs/common';
import { plugins } from 'plugins';

export const providers = plugins.map(
  (plugin) =>
    <ClassProvider<Plugin>>{
      provide: plugin.name,
      useClass: plugin,
    },
);

export class PluginModule {
  static async forRootAsync(): Promise<DynamicModule> {
    return {
      module: PluginModule,
      providers,
      exports: [...providers],
    };
  }
}
