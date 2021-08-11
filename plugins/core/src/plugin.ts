import { DynamicModule, Provider } from '@nestjs/common';
import Scheduler from './scheculer';

export type ScheduleProvider = {
  name: string;
  schedule: typeof Scheduler;
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
