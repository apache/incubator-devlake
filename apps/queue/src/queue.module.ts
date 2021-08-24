import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { ConsumerModule } from './consumer';
import { ProducerModule } from './producer';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    ProducerModule.forRoot(),
    ConsumerModule.forRoot(),
  ],
  providers: [],
})
export class QueueModule {}
