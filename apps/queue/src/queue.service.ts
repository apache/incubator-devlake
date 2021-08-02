import { Injectable } from '@nestjs/common';

@Injectable()
export class QueueService {
  getHello(): string {
    return 'Hello World!';
  }
}
