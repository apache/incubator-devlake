import { DynamicModule, Provider, Type } from '@nestjs/common';
import Scheduler from './scheculer';

export type ScheduleProvider = {
  name: string;
  schedule: Type<Scheduler<any, any>>;
};

class PluginModule {
  static Register(schedules: ScheduleProvider[]): DynamicModule {
    const providers: Provider[] = schedules.map((schedule) => ({
      provide: schedule.name,
      useClass: schedule.schedule,
    }));
    return {
      module: PluginModule,
      providers,
      exports: [...providers],
    };
  }
}

export default PluginModule;
