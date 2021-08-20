import { DynamicModule, Provider, Type } from '@nestjs/common';
import CustomTypeOrmModule from 'apps/rest/src/providers/typeorm.module';
import DependencyResolver from './dependency.resolver';
import { PRODUCER_META_KEY } from './exports.decorator';
import { IMPORTS_META_KEY } from './imports.decorator';
import PluginInterface from './plugin.interface';
import Task from './task.interface';

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
      const entities = i.exports();
      const registTask = (task: Type<Task>) => {
        providers.push({
          provide: task.name,
          useClass: task,
        });
        const importEntities = Reflect.getMetadata(IMPORTS_META_KEY, task);
        for (const e of importEntities) {
          entities.push(e);
          const DepTask = Reflect.getMetadata(PRODUCER_META_KEY, e);
          if (DepTask) {
            registTask(DepTask);
          }
        }
      };
      for (const entity of entities) {
        const Task = Reflect.getMetadata(PRODUCER_META_KEY, entity);
        entities.push(entity);
        if (Task) {
          registTask(Task);
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
