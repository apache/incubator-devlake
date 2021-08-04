import { BullModule } from '@nestjs/bull';
import { DynamicModule } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import Bull from 'bull';

export class BullQueueModule {
  static forRoot(queue = 'default'): DynamicModule {
    const defaultQueueModule = BullModule.registerQueueAsync({
      name: queue,
      useFactory: (config: ConfigService): Bull.QueueOptions => {
        const redis = config.get<string>('REDIS_URL');
        return {
          redis,
          defaultJobOptions: {
            attempts: 3,
          },
        };
      },
      inject: [ConfigService],
    });
    return {
      module: BullQueueModule,
      imports: [defaultQueueModule],
      exports: [defaultQueueModule],
      global: true,
    };
  }
}
