import { ClassProvider, DynamicModule, Type } from '@nestjs/common';
import PluginInterface from './plugin.interface';

export default class PluginModule {
  static async forRootAsync(
    plugins: Type<PluginInterface>[],
  ): Promise<DynamicModule> {
    const providers: ClassProvider[] = plugins.map((p) => ({
      provide: p.name,
      useClass: p,
    }));
    return {
      module: PluginModule,
      providers,
      exports: [...providers],
    };
  }
}
