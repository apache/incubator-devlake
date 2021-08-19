import { Inject, Injectable } from '@nestjs/common';
import redis from 'ioredis';

export type EventHandler = (payload: any) => void;

@Injectable()
export class EventsService {
  constructor(
    @Inject('REDIS_PUB_CLIENT') private pub: redis.Redis,
    @Inject('REDIS_SUB_CLIENT') private sub: redis.Redis,
  ) {}

  async emit(event: string, payload: any): Promise<number> {
    return this.pub.publish(event, JSON.stringify(payload));
  }

  async on(event: string, handler: EventHandler): Promise<void> {
    await this.sub.subscribe(event);
    this.sub.on('message', (channel, message) => {
      if (channel === event) {
        const args = JSON.parse(message);
        handler(args);
      }
    });
  }
}
