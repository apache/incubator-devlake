import { Inject, Module, OnModuleDestroy } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import Redis from 'ioredis';
import { TasksService } from './tasks.services';
import { URL } from 'url';
import { ProducerModule } from 'apps/queue/src/producer';
import { EventsModule } from '../events/events.module';

@Module({
  imports: [ConfigModule.forRoot(), ProducerModule.forRoot(), EventsModule],
  providers: [
    {
      provide: 'REDIS_TASK_CLIENT',
      useFactory: (config: ConfigService): Redis.Redis => {
        const redisUrl = config.get<string>('REDIS_URL');
        const url = new URL(redisUrl);
        return new Redis({
          host: url.hostname,
          port: parseInt(url.port),
          db: parseInt(url.pathname.slice(1)),
          username: url.username,
          password: url.password,
        });
      },
      inject: [ConfigService],
    },
    {
      provide: TasksService,
      useClass: TasksService,
    },
  ],
  exports: [TasksService],
})
export class TasksModule implements OnModuleDestroy {
  constructor(@Inject('REDIS_TASK_CLIENT') private redis: Redis.Redis) {}

  async onModuleDestroy(): Promise<void> {
    await this.redis.disconnect();
  }
}
