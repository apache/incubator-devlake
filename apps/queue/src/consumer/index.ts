import { DynamicModule } from '@nestjs/common';
import { BullQueueModule } from '../bull/queue.module';
import PluginModule from 'plugins/core/src/plugins.module';
import { ConsumerService } from './service';
import { plugins } from '../../../../plugins';

export class ConsumerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ConsumerModule,
      imports: [
        BullQueueModule.forRoot(queue),
        PluginModule.forRootAsync(plugins),
        ConsumerService,
      ],
      providers: [],
      exports: [ConsumerService],
    };
  }
}
