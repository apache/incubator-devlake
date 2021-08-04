import { DynamicModule } from '@nestjs/common';
import { BullQueueModule } from '../bull/queue.module';
import { providers } from './providers';
import { ProducerService } from './service';

export class ProducerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ProducerModule,
      imports: [BullQueueModule.forRoot(queue), ProducerService],
      providers,
      exports: [ProducerService],
    };
  }
}
