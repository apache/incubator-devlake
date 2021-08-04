import { DynamicModule } from '@nestjs/common';
import { BullQueueModule } from '../bull/queue.module';
import { providers } from './providers';
import { ConsumerService } from './service';

export class ConsumerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ConsumerModule,
      imports: [BullQueueModule.forRoot(queue), ConsumerService],
      providers,
      exports: [ConsumerService],
    };
  }
}
