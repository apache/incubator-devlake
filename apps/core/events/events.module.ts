import { DynamicModule, Inject, OnModuleDestroy } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import Redis from 'ioredis';
import { EventsService } from './events.service';
import { URL } from 'url';

export class EventsModule implements OnModuleDestroy {
  static forRoot(): DynamicModule {
    return {
      module: EventsModule,
      global: true,
      imports: [ConfigModule.forRoot()],
      providers: [
        {
          provide: 'REDIS_PUB_CLIENT',
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
          provide: 'REDIS_SUB_CLIENT',
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
          provide: EventsService,
          useClass: EventsService,
        },
      ],
      exports: [EventsService],
    };
  }

  constructor(
    @Inject('REDIS_PUB_CLIENT') private pub: Redis.Redis,
    @Inject('REDIS_SUB_CLIENT') private sub: Redis.Redis,
  ) {}
  async onModuleDestroy(): Promise<void> {
    await this.pub.disconnect();
    await this.sub.disconnect();
  }
}
