import { DynamicModule } from '@nestjs/common';
import { BullQueueModule } from '../bull/queue.module';
import { PluginModule } from './plugins';
import { ConsumerService } from './service';

export class ConsumerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ConsumerModule,
      imports: [
        BullQueueModule.forRoot(queue),
        PluginModule.forRootAsync(),
        ConsumerService,
      ],
      providers: [],
      exports: [ConsumerService],
    };
  }
}
