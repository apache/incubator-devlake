import { ClassProvider, DynamicModule, Type } from '@nestjs/common';
import DependencyResolver from './dependency.resolver';
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
      imports: [],
      providers: [...providers, DependencyResolver],
      exports: [...providers, DependencyResolver],
    };
  }
}
