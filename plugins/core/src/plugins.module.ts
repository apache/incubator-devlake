import { DynamicModule, Provider, Type } from '@nestjs/common';
import DependencyResolver from './dependency.resolver';
import { EXPORTS_META_KEY } from './exports.decorator';
import PluginInterface from './plugin.interface';

export default class PluginModule {
  static async forRootAsync(
    plugins: Type<PluginInterface>[],
  ): Promise<DynamicModule> {
    const providers: Provider[] = [DependencyResolver];
    const entities = [];
    for (const plugin of plugins) {
      providers.push({
        provide: plugin.name,
        useClass: plugin,
      });
      const i = new plugin();
      const plgunExports = i.exports();
      for (const e of plgunExports) {
        const tasks = DependencyResolver.resolveEntity(e);
        for (const t of tasks) {
          providers.push({
            provide: t.name,
            useClass: t,
          });
          entities.push(Reflect.getMetadata(EXPORTS_META_KEY, t));
        }
      }
    }
    return {
      module: PluginModule,
      imports: [],
      providers: [...providers],
      exports: [...providers],
    };
  }
}
