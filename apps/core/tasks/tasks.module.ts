import { Inject, Module, OnModuleDestroy } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import Redis from 'ioredis';
import { TasksService } from './tasks.services';
import { URL } from 'url';

@Module({
  imports: [ConfigModule.forRoot()],
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
  constructor(
    @Inject('REDIS_PUB_CLIENT') private pub: Redis.Redis,
    @Inject('REDIS_SUB_CLIENT') private sub: Redis.Redis,
  ) {}
  async onModuleDestroy(): Promise<void> {
    await this.pub.disconnect();
    await this.sub.disconnect();
  }
}