import { DynamicModule, Provider, Type } from '@nestjs/common';
import { CollectorMap, COLLECTORS_METADATA } from './collector.decorate';
import CollectorRef from './collectorref';
import Scheduler from './scheculer';

export type ScheduleProvider = {
  name: string;
  schedule: Type<Scheduler<any, any>>;
};

class PluginModule {
  static Register(schedules: ScheduleProvider[]): DynamicModule {
    const providers: Provider[] = [];
    schedules.forEach((schedule) => {
      providers.push({
        provide: schedule.name,
        useClass: schedule.schedule,
      });
      const collecotrs: CollectorMap = Reflect.getMetadata(
        COLLECTORS_METADATA,
        schedule.schedule,
      );
      Object.keys(collecotrs).forEach((collectorName) => {
        providers.push({
          provide: `${schedule.name}/collector/${collectorName}`,
          useClass: collecotrs[collectorName],
        });
      });
    });
    return {
      module: PluginModule,
      imports: [],
      providers: [CollectorRef, ...providers],
      exports: [...providers, CollectorRef],
    };
  }
}

export default PluginModule;
