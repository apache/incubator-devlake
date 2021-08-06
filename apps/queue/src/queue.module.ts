import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { ConsumerModule } from './consumer';
import { ProducerModule } from './producer';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    ConsumerModule.forRoot(),
    ProducerModule.forRoot(),
  ],
  providers: [],
})
export class QueueModule {}
