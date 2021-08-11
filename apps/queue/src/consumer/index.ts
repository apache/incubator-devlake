import { DynamicModule } from '@nestjs/common';
import { BullQueueModule } from '../bull/queue.module';
import PluginModule from 'plugins/core/src/plugin';
import { ConsumerService } from './service';
import Jira from 'plugins/jira/src';

export class ConsumerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ConsumerModule,
      imports: [
        BullQueueModule.forRoot(queue),
        PluginModule.Register([{ name: 'Jira', schedule: Jira }]),
        ConsumerService,
      ],
      providers: [],
      exports: [ConsumerService],
    };
  }
}
