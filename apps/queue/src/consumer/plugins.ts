import { DynamicModule } from '@nestjs/common';
import { modules } from 'plugins';

export class PluginModule {
  static async forRootAsync(): Promise<DynamicModule> {
    return {
      module: PluginModule,
      imports: modules,
      exports: [...modules],
    };
  }
}
