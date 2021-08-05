import { BullModule } from '@nestjs/bull';
import { DynamicModule } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import Bull from 'bull';
import { URL } from 'url';

export class BullQueueModule {
  static forRoot(queue = 'default'): DynamicModule {
    const defaultQueueModule = BullModule.registerQueueAsync({
      name: queue,
      useFactory: (config: ConfigService): Bull.QueueOptions => {
        const redis = config.get<string>('REDIS_URL');
        const redisUrl = new URL(redis);
        return {
          redis: {
            host: redisUrl.hostname,
            port: parseInt(redisUrl.port),
            db: parseInt(redisUrl.pathname.slice(1)),
            username: redisUrl.username,
            password: redisUrl.pathname,
          },
          defaultJobOptions: {
            attempts: 3,
          },
        };
      },
      imports: [ConfigModule],
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
