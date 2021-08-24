import { DynamicModule } from '@nestjs/common';
import { BullQueueModule } from '../bull/queue.module';
import PluginModule from 'plugins/core/src/plugins.module';
import { ConsumerService } from './service';
import Jira from 'plugins/jira/src';
import { EventsModule } from 'apps/core/events/events.module';
import { TasksModule } from 'apps/core/tasks/tasks.module';

export class ConsumerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ConsumerModule,
      imports: [
        BullQueueModule.forRoot(queue),
        PluginModule.forRootAsync([Jira]),
        EventsModule.forRoot(),
        TasksModule,
      ],
      providers: [{ provide: ConsumerService, useClass: ConsumerService }],
      exports: [ConsumerService],
    };
  }
}
